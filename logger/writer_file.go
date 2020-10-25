package logger

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	env "github.com/hanjingo/golib/env"
	file "github.com/hanjingo/golib/file"
)

//文件打印器
type FileWriter struct {
	mu         sync.Mutex         //原子锁
	isValid    bool               //是否可用
	writerType uint32             //打印机类型
	name       string             //打印机名字
	checkCount int                //检查计数器
	targets    map[string]string  //地址
	levels     map[int]LevelInfoI //等级
	cache      chan string        //缓存
	cacheCapa  int                //缓存容量
	nfFunc     func() string      //新建文件函数
	checkDur   int                //检查周期(次数)
	expireDur  time.Duration      //保质期
	capacity   int64              //文件容量
	flushDur   time.Duration      //刷新周期
	fm         *FileManager1      //文件管理器
}

//新建一个文件打印器
func NewFileWriter(conf *FileWriterConfig) *FileWriter {
	if !conf.Check() {
		return nil
	}
	var back *FileWriter
	back = &FileWriter{
		isValid:    true,
		writerType: FILE,
		name:       conf.WriterName,
		targets:    make(map[string]string),
		levels:     make(map[int]LevelInfoI),
		cache:      make(chan string, conf.CacheCapa),
		cacheCapa:  conf.CacheCapa,
		checkDur:   conf.CheckDur,
		expireDur:  time.Duration(conf.ExpireDur) * time.Hour,
		capacity:   conf.Capacity,
		fm:         NewFileManager1(conf.MaxFileNum),
	}
	back.Init(conf.OpenLevels)
	return back
}

//初始化打印机
func (w *FileWriter) Init(openLevels []int) { //默认打开 debug 到 fatal 如果没有target new一个
	for id, _ := range w.levels {
		delete(w.levels, id)
	}
	for _, lvl := range openLevels {
		w.SetLevel(NewLevel(lvl, defaultHeadFunc))
	}
}

//返回打印机类型
func (w *FileWriter) Type() uint32 {
	return w.writerType
}

//返回打印机名字
func (w *FileWriter) Name() string {
	return w.name
}

//设置打印地址
func (w *FileWriter) SetTarget(arg string) { //一个打印机只支持一个地址
	target := env.GetCurrPath()
	if !filepath.IsAbs(arg) {
		target = filepath.Join(env.GetCurrPath(), arg)
	} else {
		target = arg
	}
	for key, _ := range w.targets {
		delete(w.targets, key)
	}
	w.targets[target] = target
}

//设置日志等级
func (w *FileWriter) SetLevel(lvl LevelInfoI) {
	w.levels[lvl.Level()] = lvl
}

//关闭某个日志等级
func (w *FileWriter) CloseLevel(lvl int) {
	if _, ok := w.levels[lvl]; ok {
		w.levels[lvl].SetValid(false)
	}
}

//写日志
func (w *FileWriter) Write(lvl int, str string) {
	if !w.isValid {
		return
	}
	if !w.IsLevelOpen(lvl) {
		return
	}
	w.checkCount--
	if w.checkCount <= 0 {
		w.checkAndDeal()
	}
	msg := w.levels[lvl].Head()
	msg += str
	if len(w.cache) >= w.cacheCapa {
		w.Flush()
	}
	w.cache <- msg
	w.Flush()
}

//判断等级是否打开
func (w *FileWriter) IsLevelOpen(lvl int) bool {
	if _, ok := w.levels[lvl]; !ok {
		return false
	}
	if !w.levels[lvl].IsValid() {
		return false
	}
	return true
}

//刷入
func (w *FileWriter) Flush() {
	if len(w.targets) == 0 {
		if w.nfFunc != nil {
			w.SetTarget(w.nfFunc())
		}
	}
	w.doFlush()
}

//关闭打印器
func (w *FileWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isValid = false
	w.levels = nil
	close(w.cache)
}

//设置是否可用
func (w *FileWriter) SetValid(isValid bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isValid = isValid
}

/*******************other func********************/

//设置新建文件函数
func (w *FileWriter) SetNewFileFunc(f func() string) {
	w.nfFunc = f
}

//检查处理
func (w *FileWriter) checkAndDeal() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.checkCount = w.checkDur
	for _, addr := range w.targets {
		if !w.check(addr) {
			if w.nfFunc == nil {
				continue
			}
			w.fm.CloseFile(addr)
			var filePathName string
			if w.nfFunc != nil {
				filePathName = w.nfFunc()
			}
			_, err := w.fm.OpenFile(filePathName)
			if err != nil {
				fmt.Printf("新建日志文件失败,错误:%v", err)
				continue
			}
			w.SetTarget(filePathName)
		}
	}
}

//检查日志文件
func (w *FileWriter) check(addr string) bool {
	info := w.fm.GetFileInfo(addr)
	if info == nil {
		return false
	}
	if info.IsExpired(w.expireDur) {
		return false
	}
	if info.IsFull(file.SIZE(w.capacity)) {
		return false
	}
	return true
}

func (w *FileWriter) doFlush() {
	for len(w.cache) > 0 {
		msg := <-w.cache
		for _, target := range w.targets {
			w.fm.Write(target, []byte(msg))
		}
	}
}

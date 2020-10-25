package logger

import (
	"fmt"
	"sync"
	"time"
)

type ConsoleWriter struct {
	mu         sync.Mutex //原子锁
	isValid    bool
	writerType uint32
	name       string
	levels     map[int]LevelInfoI
	cache      chan string
	cacheCapa  int
	flushDur   time.Duration
}

func NewConsoleWriter(conf *ConsoleWriterConfig) *ConsoleWriter {
	var back *ConsoleWriter
	if !conf.Check() {
		return nil
	} else {
		back = &ConsoleWriter{
			isValid:    true,
			writerType: CONSOLE,
			name:       conf.WriterName,
			levels:     make(map[int]LevelInfoI),
			cache:      make(chan string, conf.CacheCapa),
			cacheCapa:  conf.CacheCapa,
		}
	}
	back.Init(conf.OpenLevels)
	return back
}

//初始化打印机
func (w *ConsoleWriter) Init(openLevels []int) {
	for id, _ := range w.levels {
		delete(w.levels, id)
	}
	for _, lvl := range openLevels {
		w.SetLevel(NewLevel(lvl, defaultHeadFunc))
	}
}

//返回打印机类型
func (w *ConsoleWriter) Type() uint32 {
	return w.writerType
}

//返回打印机名字
func (w *ConsoleWriter) Name() string {
	return w.name
}

//设置打印目标
func (w *ConsoleWriter) SetTarget(target string) {
	//todo
}

//设置日志等级
func (w *ConsoleWriter) SetLevel(lvl LevelInfoI) {
	w.levels[lvl.Level()] = lvl
}

//关闭日志等级
func (w *ConsoleWriter) CloseLevel(lvl int) {
	if _, ok := w.levels[lvl]; ok {
		w.levels[lvl].SetValid(false)
	}
}

//写日志
func (w *ConsoleWriter) Write(lvl int, str string) {
	if !w.isValid {
		return
	}
	if !w.IsLevelOpen(lvl) {
		return
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
func (w *ConsoleWriter) IsLevelOpen(lvl int) bool {
	if _, ok := w.levels[lvl]; !ok {
		return false
	}
	if !w.levels[lvl].IsValid() {
		return false
	}
	return true
}

//输入
func (w *ConsoleWriter) Flush() {
	w.doFlush()
}

//关闭
func (w *ConsoleWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isValid = false
	w.levels = nil
	close(w.cache)
}

//设置是否可用
func (w *ConsoleWriter) SetValid(isValid bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isValid = isValid
}

/***********other func************/

func (w *ConsoleWriter) doWrite(content string) {
	fmt.Printf(content)
}

func (w *ConsoleWriter) check() bool {
	return true
}

func (w *ConsoleWriter) doFlush() {
	for len(w.cache) > 0 {
		str := <-w.cache
		w.doWrite(str)
	}
}

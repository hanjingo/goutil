package logger

import (
	"sync"

	ws "github.com/gorilla/websocket"
)

type NetWriter struct {
	mu         sync.Mutex          //原子锁
	isValid    bool                //是否可用
	name       string              //打印机名字
	checkCount int                 //检查计数器
	writerType uint32              //打印机类型
	targets    map[string]string   //打印类型
	levels     map[int]LevelInfoI  //等级
	conns      map[string]*ws.Conn //key:target value:conn
	cache      chan string         //缓存
	cacheCapa  int                 //缓存容量
	checkDur   int                 //检查周期
}

func NewNetWriter(conf *NetWriterConfig) *NetWriter {
	var back *NetWriter
	if !conf.Check() {
		return nil
	} else {
		back = &NetWriter{
			isValid:    true,
			name:       conf.WriterName,
			writerType: NET,
			targets:    make(map[string]string),
			levels:     make(map[int]LevelInfoI),
			conns:      make(map[string]*ws.Conn),
			cache:      make(chan string, conf.CacheCapa),
			cacheCapa:  conf.CacheCapa,
			checkDur:   conf.CheckDur,
		}
	}
	back.Init(conf.OpenLevels)
	return back
}

//初始化打印机
func (w *NetWriter) Init(openLevels []int) { //默认打开 debug 到 fatal 如果没有target new一个
	for id, _ := range w.levels {
		delete(w.levels, id)
	}
	for _, lvl := range openLevels {
		w.SetLevel(NewLevel(lvl, defaultHeadFunc))
	}
}

//类型
func (w *NetWriter) Type() uint32 {
	return w.writerType
}

//返回打印机名字
func (w *NetWriter) Name() string {
	return w.name
}

//设置打印地址
func (w *NetWriter) SetTarget(target string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for key, _ := range w.targets {
		delete(w.targets, key)
	}
	w.targets[target] = target
}

//设置等级
func (w *NetWriter) SetLevel(lvl LevelInfoI) {
	w.levels[lvl.Level()] = lvl
}

//关闭日志等级
func (w *NetWriter) CloseLevel(lvl int) {
	if _, ok := w.levels[lvl]; ok {
		w.levels[lvl].SetValid(false)
	}
}

//判断等级是否打开
func (w *NetWriter) IsLevelOpen(lvl int) bool {
	if _, ok := w.levels[lvl]; !ok {
		return false
	}
	if !w.levels[lvl].IsValid() {
		return false
	}
	return true
}

//写
func (w *NetWriter) Write(lvl int, str string) {
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
}

//刷
func (w *NetWriter) Flush() {
	w.doFlush()
}

//关闭打印器
func (w *NetWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isValid = false
	w.levels = nil
	w.conns = nil
	close(w.cache)
}

//设置是否可用
func (w *NetWriter) SetValid(isValid bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isValid = isValid
}

/*******************other func********************/

func (w *NetWriter) dial() {
	for addr, _ := range w.targets {
		if _, ok := w.conns[addr]; !ok {
			conn, _, err := ws.DefaultDialer.Dial(addr, nil)
			if err != nil {
				continue
			}
			w.conns[addr] = conn
		}
	}
}

func (w *NetWriter) checkAndDeal() {
	for key, _ := range w.conns {
		if !w.isConnValid(key) {
			w.mu.Lock()
			w.delConn(key)
			w.mu.Unlock()
		}
	}
	w.dial()
	w.checkCount = w.checkDur
}

//判断连接是否可用
func (w *NetWriter) isConnValid(key string) bool {
	if _, ok := w.conns[key]; !ok {
		return false
	}
	if err := w.conns[key].WriteMessage(ws.TextMessage, []byte("ping")); err != nil {
		return false
	}
	return true
}

func (w *NetWriter) doFlush() {
	for len(w.cache) > 0 {
		msg := <-w.cache
		for _, conn := range w.conns {
			conn.WriteMessage(ws.TextMessage, []byte(msg))
		}
	}
}

func (w *NetWriter) delConn(key string) {
	if conn, ok := w.conns[key]; ok {
		conn.Close()
		delete(w.conns, key)
	}
}

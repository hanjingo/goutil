package network

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	ws "github.com/gorilla/websocket"
)

type WsServer struct {
	addr           string       //地址
	upgrader       *ws.Upgrader //websocket更新器
	maxHeaderBytes int          //ws协议头大小限制
	li             net.Listener //监听器
	doNewConn      TFNewConn    //新建连接回调
	doHandleMsg    TFHandleMsg  //新到消息回调
	doConnClose    TFConnClose  //关闭连接回调
}

func NewWsServer(addr string, f1 TFNewConn, f2 TFHandleMsg, f3 TFConnClose) (*WsServer, error) {
	back := &WsServer{
		addr: addr,
		upgrader: &ws.Upgrader{
			HandshakeTimeout: DefultHandShakeTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true }},
		maxHeaderBytes: DefultHeaderBytes,
		doNewConn:      f1,
		doHandleMsg:    f2,
		doConnClose:    f3,
	}
	return back, nil
}

func (s *WsServer) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	go s.listen(wg)
}

func (s *WsServer) listen(wg *sync.WaitGroup) {
	defer wg.Done()
	li, err := net.Listen("tcp", s.addr)
	if li == nil || err != nil {
		return
	}
	s.li = li
	httpServer := &http.Server{
		Addr:           s.addr,
		Handler:        s, //这里指定了wsserver，一有url请求就会跳到wsserver的serverhttp方法
		ReadTimeout:    DefultReadTimeout,
		WriteTimeout:   DefultWriteTimeout,
		MaxHeaderBytes: s.maxHeaderBytes,
	}
	httpServer.Serve(s.li)
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "websocket不支持此方法", 405)
		return
	}
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn, err := NewWsConn(c, s.doConnClose, s.doHandleMsg)
	if err != nil {
		return
	}
	go s.doNewConn(conn)
}

func (s *WsServer) Close() {
	if s.li != nil {
		s.li.Close()
	}
}

func (s *WsServer) Type() TYP_NET {
	return WS
}

func (s *WsServer) Addr() string {
	return s.addr
}

func (s *WsServer) SetParam(key string, value interface{}) error {
	switch strings.ToUpper(key) {
	default:
		return errors.New(fmt.Sprintf("不存在此参数:%s", key))
	}
}

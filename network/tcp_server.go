package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

type TcpServer struct {
	addr        string       //地址
	li          net.Listener //监听器
	doNewConn   TFNewConn    //新建连接回调
	doHandleMsg TFHandleMsg  //新到消息回调
	doConnClose TFConnClose  //关闭连接回调
}

func NewTcpServer(addr string, f1 TFNewConn, f2 TFHandleMsg, f3 TFConnClose) (*TcpServer, error) {
	back := &TcpServer{
		addr:        addr,
		doNewConn:   f1,
		doHandleMsg: f2,
		doConnClose: f3,
	}
	return back, nil
}

func (s *TcpServer) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	go s.listen(wg)
}

func (s *TcpServer) listen(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	if s.li, err = net.Listen("tcp", s.addr); s.li == nil || err != nil {
		return
	}
	for {
		c, err := s.li.Accept()
		if err != nil {
			return
		}
		conn, err := NewTcpConn(c, s.doConnClose, s.doHandleMsg)
		if err != nil {
			continue
		}
		go s.doNewConn(conn)
	}
}

func (s *TcpServer) Close() {
	if s.li != nil {
		s.li.Close()
	}
}

func (s *TcpServer) Type() TYP_NET {
	return TCP
}

func (s *TcpServer) Addr() string {
	return s.addr
}

func (s *TcpServer) SetParam(key string, value interface{}) error {
	switch strings.ToUpper(key) {
	default:
		return errors.New(fmt.Sprintf("不存在此参数:%s", key))
	}
}

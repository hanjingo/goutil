package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type TcpCli struct {
	addr        string
	dialer      *net.Dialer //拨号器
	doNewConn   TFNewConn   //新建连接回调
	doHandleMsg TFHandleMsg //新到消息回调
	doConnClose TFConnClose //关闭连接回调
}

func NewTcpCli(addr string, f1 TFNewConn, f2 TFHandleMsg, f3 TFConnClose) (*TcpCli, error) {
	back := &TcpCli{
		addr:        addr,
		doNewConn:   f1,
		doHandleMsg: f2,
		doConnClose: f3,
	}
	return back, nil
}

func (cli *TcpCli) Dial(addr string) (*TcpConn, error) {
	var conn net.Conn
	var err error
	if cli.dialer != nil {
		conn, err = cli.dialer.Dial("tcp", addr)
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return nil, err
	}
	back, err := NewTcpConn(conn, cli.doConnClose, cli.doHandleMsg)
	if err != nil {
		return nil, err
	}
	go cli.doNewConn(back)
	return back, nil
}

func (cli *TcpCli) SetParam(key string, value interface{}) error {
	switch strings.ToUpper(key) {
	case "DIALER":
		dialer, ok := value.(*net.Dialer)
		if !ok {
			return errors.New(fmt.Sprintf("非法的拨号器类型:%v", value))
		}
		cli.dialer = dialer
		return nil
	default:
		return errors.New(fmt.Sprintf("不存在此参数:%s", key))
	}
}

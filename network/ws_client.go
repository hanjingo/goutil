package network

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	ws "github.com/gorilla/websocket"
)

type WsCli struct {
	addr        string
	dialer      *ws.Dialer  //拨号器
	doNewConn   TFNewConn   //新建连接回调
	doHandleMsg TFHandleMsg //新到消息回调
	doConnClose TFConnClose //关闭连接回调
}

func NewWsCli(addr string, f1 TFNewConn, f2 TFHandleMsg, f3 TFConnClose) (*WsCli, error) {
	back := &WsCli{
		addr:        addr,
		dialer:      ws.DefaultDialer,
		doNewConn:   f1,
		doHandleMsg: f2,
		doConnClose: f3,
	}
	return back, nil
}

func (cli *WsCli) Dial(url string, opts ...http.Header) (*WsConn, error) {
	var header http.Header
	if opts != nil && len(opts) > 0 {
		header = opts[0]
	}
	conn, _, err := cli.dialer.Dial(url, header)
	if err != nil {
		return nil, err
	}
	back, err := NewWsConn(conn, cli.doConnClose, cli.doHandleMsg)
	if err != nil {
		return nil, err
	}
	go cli.doNewConn(back)
	return back, nil
}

func (cli *WsCli) SetParam(key string, value interface{}) error {
	switch strings.ToUpper(key) {
	case "DIALER":
		dialer, ok := value.(*ws.Dialer)
		if !ok {
			return errors.New(fmt.Sprintf("非法的拨号器类型:%v", value))
		}
		cli.dialer = dialer
		return nil
	default:
		return errors.New(fmt.Sprintf("不存在此参数:%s", key))
	}
}

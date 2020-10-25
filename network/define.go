package network

import (
	"sync"
	"time"
)

type TYP_NET string

type TFConnClose func(SessionI)
type TFNewConn func(SessionI)
type TFHandleMsg func(SessionI, int)

type ServerI interface {
	SetParam(string, interface{}) error //设置配置
	Run(wg *sync.WaitGroup)             //服务器跑起来
	Close()                             //关闭
	Type() TYP_NET                      //服务器类型
	Addr() string                       //服务器监听地址
}

type SessionI interface {
	SetParam(string, interface{}) error   //设置配置
	Run() error                           //跑起来
	ReadMsg(arg []byte) (int, error)      //读消息
	WriteMsg(args ...[]byte) (int, error) //写消息
	Close()                               //关闭连接 会等消息发完
	Destroy()                             //摧毁连接 即使有消息也会强制关掉连接
}

type CliI interface {
	SetParam(string, interface{}) error //设置配置
	Dial(string) (SessionI, error)      //拨号
}

var PackSize int = 1400 //数据包大小 默认MTU

var DefultHeaderBytes int = 1024                                                  //默认协议头大小
var DefultHandShakeTimeout time.Duration = time.Duration(3000) * time.Millisecond //默认握手超时(3000ms)
var DefultReadTimeout time.Duration = time.Duration(3000) * time.Millisecond      //默认读超时(3000ms)
var DefultWriteTimeout time.Duration = time.Duration(3000) * time.Millisecond     //默认写超时(3000ms)

const (
	TCP  TYP_NET = "tcp"
	UDP  TYP_NET = "udp"
	WS   TYP_NET = "ws"
	HTTP TYP_NET = "http"
)

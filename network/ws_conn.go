package network

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"

	ws "github.com/gorilla/websocket"
)

type WsConn struct {
	mu           *sync.RWMutex      //互斥锁
	conn         *ws.Conn           //ws框架提供的连接
	bValid       bool               //是否可用
	bRun         bool               //是否已运行
	bCanRead     bool               //是否可读
	bCanWrite    bool               //是否可写
	readBufSize  int                //读缓冲区大小
	readBuf      *bytes.Buffer      //读缓存
	writeBufSize int                //写缓冲区大小
	writeBuf     *bytes.Buffer      //写缓存
	finish       context.CancelFunc //结束函数
	doConnClose  TFConnClose        //连接关闭回调
	doHandleMsg  TFHandleMsg        //消息回调
}

func NewWsConn(conn *ws.Conn, f1 TFConnClose, f2 TFHandleMsg) (*WsConn, error) {
	if conn == nil {
		return nil, errors.New("conn不能为空")
	}
	back := &WsConn{
		mu:           new(sync.RWMutex),
		bValid:       false,
		bRun:         false,
		readBufSize:  0,
		readBuf:      new(bytes.Buffer),
		writeBufSize: 0,
		writeBuf:     new(bytes.Buffer),
		conn:         conn,
		doConnClose:  f1,
		doHandleMsg:  f2,
	}
	return back, nil
}

//跑起来
func (c *WsConn) Run() error {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("tcp连接运行失败,错误:%v\n", string(buf[:n]))
		}
	}()
	if c.bRun {
		return errors.New("已经启动，无需再次启动")
	}
	c.bRun = true
	c.bValid = true
	c.bCanRead = true
	c.bCanWrite = true
	ctx, cancel := context.WithCancel(context.Background())
	c.finish = cancel
	go c.goRead(ctx)
	go c.goWrite(ctx)
	return nil
}

//摧毁连接
func (c *WsConn) Destroy() {
	if !c.IsValid() {
		return
	}
	c.bValid = false
	if c.conn != nil {
		c.conn.Close()
	}
	if c.doConnClose != nil {
		c.doConnClose(c)
	}
	c.writeBuf.Reset()
	c.readBuf.Reset()
}

//读消息(阻塞)
func (c *WsConn) ReadMsg(arg []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.bValid {
		return 0, errors.New("连接已摧毁")
	}
	if !c.bCanRead {
		return 0, errors.New("连接不可读")
	}
	return c.readBuf.Read(arg)
}

//写消息(阻塞)
func (c *WsConn) WriteMsg(args ...[]byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	back := 0
	if !c.bValid {
		return back, errors.New("连接已摧毁")
	}
	if !c.bCanWrite {
		return back, errors.New("连接不可写")
	}
	for _, arg := range args {
		if _, err := c.writeBuf.Write(arg); err != nil {
			return back, err
		}
		back += len(arg)
	}
	return back, nil
}

//安全的关闭连接
func (c *WsConn) Close() {
	if !c.IsValid() {
		return
	}
	if c.finish != nil {
		c.finish()
	}
	if c.doConnClose != nil {
		c.doConnClose(c)
	}
}

//连接是否可用
func (c *WsConn) IsValid() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bValid
}

//连接不可读
func (c *WsConn) disEnbleRead() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bCanRead = false
}

//读
func (c *WsConn) goRead(ctx context.Context) {
	capa := 1
	if c.readBufSize > 0 {
		capa = (c.readBufSize / PackSize) + 1
	}
	readC := make(chan []byte, capa)
	defer func() {
		for len(readC) > 0 {
			c.readBuf.Write(<-readC)
		}
		close(readC)
		c.disEnbleRead()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case arg := <-readC:
			c.readBuf.Write(arg)
			if c.doHandleMsg != nil {
				c.doHandleMsg(c, len(arg))
			}
		default:
			_, buf, err := c.conn.ReadMessage()
			if err != nil {
				return
			}
			readC <- buf
		}
	}
}

//连接不可写
func (c *WsConn) disEnbleWrite() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bCanWrite = false
}

//写
func (c *WsConn) goWrite(ctx context.Context) {
	capa := 1
	if c.writeBufSize > 0 {
		capa = (c.writeBufSize / PackSize) + 1
	}
	writeC := make(chan []byte, capa)
	defer func() {
		close(writeC)
		c.disEnbleWrite()
	}()
	for {
		select {
		case <-ctx.Done():
			for c.writeBuf.Len() > 0 {
				for len(writeC) > 0 {
					if err := c.conn.WriteMessage(ws.TextMessage, <-writeC); err != nil {
						return
					}
				}
				writeC <- c.writeBuf.Next(PackSize)
			}
			return
		case arg := <-writeC:
			if arg == nil {
				continue
			}
			if err := c.conn.WriteMessage(ws.TextMessage, arg); err != nil {
				return
			}
		case writeC <- c.writeBuf.Next(PackSize):
		}
	}
}

//设置参数
func (c *WsConn) SetParam(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch strings.ToUpper(key) {
	case "READBUFSIZE": //读容量
		n, ok := value.(int)
		if !ok {
			return errors.New("错误的参数")
		}
		if n <= 0 {
			return errors.New("参数值非法")
		}
		c.readBufSize = n
	case "WRITEBUFSIZE": //写容量
		n, ok := value.(int)
		if !ok {
			return errors.New("错误的参数")
		}
		if n <= 0 {
			return errors.New("参数值非法")
		}
		c.writeBufSize = n
	case "CanREAD": //是否可读
		b, ok := value.(bool)
		if !ok {
			return errors.New("错误的参数")
		}
		c.bCanRead = b
	case "CanWrite": //是否可写
		b, ok := value.(bool)
		if !ok {
			return errors.New("错误的参数")
		}
		c.bCanWrite = b
	default:
		return errors.New(fmt.Sprintf("不存在此参数:%s", key))
	}
	return nil
}

package network

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"sync"
)

type TcpConn struct {
	mu           *sync.RWMutex
	conn         net.Conn           //连接
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

func NewTcpConn(conn net.Conn, f1 TFConnClose, f2 TFHandleMsg) (*TcpConn, error) {
	if conn == nil {
		return nil, errors.New("conn不能为空")
	}
	back := &TcpConn{
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
	back.conn.(*net.TCPConn).SetNoDelay(true) //默认关闭negal算法
	return back, nil
}

//跑起来
func (c *TcpConn) Run() error {
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
func (c *TcpConn) Destroy() {
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
func (c *TcpConn) ReadMsg(arg []byte) (int, error) {
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
func (c *TcpConn) WriteMsg(args ...[]byte) (int, error) {
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
func (c *TcpConn) Close() {
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
func (c *TcpConn) IsValid() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bValid
}

//连接不可读
func (c *TcpConn) disEnbleRead() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bCanRead = false
}

//读
func (c *TcpConn) goRead(ctx context.Context) {
	capa := 1
	if c.readBufSize > 0 {
		capa = (c.readBufSize / PackSize) + 1
	}
	readC := make(chan []byte, capa)
	freeC := make(chan []byte, capa)
	for i := 0; i < capa; i++ {
		freeC <- make([]byte, PackSize)
	}
	defer func() {
		for len(readC) > 0 {
			c.readBuf.Write(<-readC)
		}
		close(readC)
		close(freeC)
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
			freeC <- arg
		case tmp := <-freeC:
			n, err := c.conn.Read(tmp) //阻塞
			if err != nil {
				return
			}
			readC <- tmp[:n]
		}
	}
}

//连接不可写
func (c *TcpConn) disEnbleWrite() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bCanWrite = false
}

//写
func (c *TcpConn) goWrite(ctx context.Context) {
	capa := 1
	if c.writeBufSize > 0 {
		capa = (c.writeBufSize / PackSize) + 1
	}
	writeC := make(chan []byte, capa)
	freeC := make(chan []byte, capa)
	for i := 0; i < capa; i++ {
		freeC <- make([]byte, capa)
	}
	defer func() {
		close(writeC)
		close(freeC)
		c.disEnbleWrite()
	}()
	for {
		select {
		case <-ctx.Done():
			for c.writeBuf.Len() > 0 {
				for len(writeC) > 0 {
					if _, err := c.conn.Write(<-writeC); err != nil {
						return
					}
				}
				writeC <- c.writeBuf.Next(PackSize)
			}
			return
		case arg := <-writeC:
			if _, err := c.conn.Write(arg); err != nil {
				return
			}
			freeC <- arg
		case arg := <-freeC:
			n, err := c.writeBuf.Read(arg)
			if err != nil {
				return
			}
			writeC <- arg[:n]
		}
	}
}

//设置参数
func (c *TcpConn) SetParam(key string, value interface{}) error {
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

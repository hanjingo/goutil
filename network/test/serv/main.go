package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hanjingo/golib/network"
)

func main() {
	wg := new(sync.WaitGroup)
	s, err := network.NewTcpServer(
		"127.0.0.1:10086",
		func(c network.SessionI) {
			//c.SetParam("readBufSize", 1024)
			//c.SetParam("writeBufSize", 1024)
			c.Run()
		},
		func(c network.SessionI, n int) {
			if n < 5 {
				return
			}
			data := make([]byte, 5)
			for _, err := c.ReadMsg(data); err == nil; {
				if string(data) != "hello" {
					fmt.Println("收到异常数据:", string(data))
					return
				}
				c.WriteMsg([]byte("world"))
			}
		},
		func(c network.SessionI) {
			fmt.Println("连接关闭")
		},
	)
	if err != nil {
		panic(err)
	}
	s.Run(wg)
	dur := time.Duration(1000) * time.Millisecond
	tm := time.NewTimer(dur)
	start := time.Now()
	end := start.Add(time.Duration(30) * time.Second)
	for {
		select {
		case <-tm.C:
			tm.Reset(dur)
		default:
			if time.Now().After(end) {
				return
			}
		}
	}
	wg.Wait()
}

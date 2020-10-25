package main

import (
	"fmt"
	"time"

	"github.com/hanjingo/golib/network"
)

var npack uint64 = 0
var nconn uint64 = 0
var mb uint64 = 1024 * 1024

func main() {
	max_conn := 500
	dur1 := time.Duration(10) * time.Second
	i := 0
	//创建cli
	cli, err := network.NewTcpCli(
		"127.0.0.1:10086",
		func(c network.SessionI) {
			nconn++
		},
		func(c network.SessionI, n int) {
			if n < 5 {
				return
			}
			data := make([]byte, 5)
			for _, err := c.ReadMsg(data); err == nil; {
				if string(data) == "world" {
					npack += 5
				} else {
					fmt.Println("收到异常数据:", string(data))
				}
			}
		},
		func(c network.SessionI) {
			if nconn > 0 {
				nconn--
			}
		},
	)
	if err != nil {
		fmt.Println("创建cli失败:", err)
		return
	}
	//创建连接
	for ; i < max_conn; i++ {
		go func() {
			tm := time.NewTimer(dur1)
			conn, err := cli.Dial("127.0.0.1:10086")
			if err != nil {
				fmt.Println("创建链接失败:", err)
				return
			}
			//conn.SetParam("readBufSize", 640)
			//conn.SetParam("writeBufSize", 640)
			conn.Run()
			for {
				select {
				case <-tm.C:
					conn.Destroy()
				default:
					if _, err := conn.WriteMsg([]byte("hello")); err != nil {
						return
					}
				}
			}
		}()
	}
	dur := time.Duration(500) * time.Millisecond
	tm := time.NewTimer(dur)
	start := time.Now()
	end := start.Add(time.Duration(30) * time.Second)
	for {
		select {
		case <-tm.C:
			fmt.Printf("\x1bc")
			fmt.Println("conn连接数:", nconn)
			fmt.Println("pack大小:", npack/mb, " MB")
			fmt.Println("经过时间:", time.Now().Sub(start).Seconds(), "秒")
			tm.Reset(dur)
		default:
			if time.Now().After(end) {
				return
			}
		}
	}
}

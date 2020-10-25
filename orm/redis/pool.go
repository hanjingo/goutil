package redis

import (
	"errors"
	"sync"

	rds "github.com/gomodule/redigo/redis"
)

type RdsPool struct {
	mu 			*sync.Mutex             //互斥锁
	pool       	chan RdsConn 			//连接池
	isClosed 	bool 					//是否关闭
	Capa       	int          			//连接池容量
}

func NewRdsPool(capa int) *RdsPool {
	if capa <= 0 {
		capa = 1
	}
	back := &RdsPool{
		mu: new(sync.Mutex),
		pool: make(chan RdsConn, capa),
		isClosed: false,
		Capa: capa,
	}
	return back
}

func (p *RdsPool) Clean() {
	for {
		if len(p.pool) > 0 {
			c := <-p.pool
			c.Close()
		} else {
			return
		}
	}
}

/*获得redis的连接*/
func doDial(url string) (RdsConn, error) {
	conn, err := rds.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	if conn.Err() != nil {
		return nil, conn.Err()
	}
	return conn, nil
}

func (p *RdsPool) Dial(url string, num int) error {
	for i := 0; i < num; i++ {
		conn, err := doDial(url)
		if err != nil {
			return err
		}
		p.Put(conn)
	}
	return nil
}

func (p *RdsPool) Put(conn RdsConn) {
	p.pool <- conn
	return
}

func (p *RdsPool) Get() RdsConn {
	return <-p.pool
}

func (p *RdsPool) Close() error {
	if p.isClosed {
		return errors.New("连接池已关闭")
	}
	p.Clean()
	p.isClosed = true
	close(p.pool)
	return nil
}

/*获得redis的连接*/
func Dial(url string) (RdsConn, error) {
	conn, err := rds.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	if conn.Err() != nil {
		return nil, conn.Err()
	}
	return conn, nil
}
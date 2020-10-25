package redis

import (
	redis "github.com/gomodule/redigo/redis"
	def "github.com/hanjingo/golib/orm/define"
)

type RdsConn redis.Conn

type RdsDb struct {
	name string
	pool *RdsPool
}

func NewRdsDb(conf *RdsConfig) (*RdsDb, error) {
	return &RdsDb{
		name: conf.Name(),
		pool: NewRdsPool(conf.PoolCapa),
	}, nil
}

func (rds *RdsDb) Where(...interface{}) def.ModelI {
	//todo
	return NewModel(nil, "")
}

func (rds *RdsDb) Name() string {
	return rds.name
}

func (rds *RdsDb) Close() error {
	rds.pool.Clean()
	return nil
}

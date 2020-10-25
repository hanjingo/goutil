package horm

import (
	"errors"

	mysql "github.com/hanjingo/golib/orm/mysql"
	redis "github.com/hanjingo/golib/orm/redis"
)

type Horm struct {
	dbs map[string]DBI //db集合 key:数据库名字 value:数据库
}

//拨号
func (orm *Horm) Dial(dbType string, conf interface{}) (DBI, error) {
	var db DBI
	var err error
	switch dbType {
	case DB_MYSQL:
		db, err = mysql.NewMsDb(conf.(*mysql.MsDbConfig))
	case DB_REDIS:
		db, err = redis.NewRdsDb(conf.(*redis.RdsConfig))
	default:
		return nil, errors.New("不支持的数据库类型")
	}
	if err != nil {
		return nil, err
	}
	orm.dbs[db.Name()] = db
	return db, err
}

//关闭db
func (orm *Horm) Close(name string) error {
	if _, ok := orm.dbs[name]; !ok {
		return errors.New("db不存在")
	}
	defer delete(orm.dbs, name)
	return orm.dbs[name].Close()
}

//关闭所有的db
func (orm *Horm) CloseAll() {
	for name, db := range orm.dbs {
		db.Close()
		delete(orm.dbs, name)
	}
}

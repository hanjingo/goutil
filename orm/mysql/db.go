package mysql

import (
	"database/sql"
	"fmt"

	def "github.com/hanjingo/golib/orm/define"
)

type MsDb struct {
	name string
	db   *sql.DB
}

func NewMsDb(conf *MsDbConfig) (*MsDb, error) {
	cmd := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=%v&&loc=Local",
		conf.UserName, conf.PassWord, conf.Addr, conf.DataBase, conf.Charset)
	db, err := sql.Open("mysql", cmd)
	if err != nil || db == nil {
		return nil, err
	}
	return &MsDb{name: conf.Name(), db: db}, nil
}

func (msdb *MsDb) Where(args ...interface{}) def.ModelI {
	querryKey := ""
	if args == nil || len(args) < 1 {
		return NewModel(nil, "", querryKey)
	}
	table := fmt.Sprintf("%v", args[0])
	argLen := len(args)
	for i := 1; i < argLen; i++ {
		querryKey += fmt.Sprintf("%v", args[i])
	}
	querryKey = "WHERE " + querryKey
	return NewModel(msdb.db, table, querryKey)
}

func (msdb *MsDb) Name() string {
	return msdb.name
}

func (msdb *MsDb) Close() error {
	return msdb.db.Close()
}

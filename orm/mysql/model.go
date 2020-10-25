package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

const TAG = "horm"

func NewModel(db *sql.DB, table string, querryKey string) *Model {
	return &Model{db: db, table: table, querryKey: querryKey}
}

type Model struct {
	content   interface{} 				//实际内容
	db        *sql.DB 					//数据库连接
	querryKey string					//查询的键
	fieldMap  map[string]interface{}	//查到的值
	table     string
}

func (m *Model) Do(...interface{}) error {
	return nil
}

//在mysql创建数据
func (m *Model) Create(content interface{}) error {
	if content == nil {
		return errors.New("参数不能为空")
	}
	m.content = content
	m.parseModel()
	cmd := fmt.Sprintf("INSERT INTO %s (", m.table)
	i := len(m.fieldMap)
	var values []interface{}
	for key, value := range m.fieldMap {
		i--
		cmd += key
		if i > 0 {
			cmd += ", "
		} else {
			cmd += ")"
		}
		values = append(values, value)
	}
	cmd += " VALUES ("
	for i := len(m.fieldMap); i > 1; i-- {
		cmd += "?, "
	}
	cmd += "?)"
	stmt, err := m.db.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	return err
}

func (m *Model) Update(content interface{}) error {
	if content == nil {
		return errors.New("参数不能为空")
	}
	m.content = content
	m.parseModel()
	cmd := fmt.Sprintf("UPDATE %s SET ", m.table)
	i := len(m.fieldMap)
	var values []interface{}
	for key, value := range m.fieldMap {
		i--
		cmd += key
		cmd += "=?" 
		if i > 0 {
			cmd += ", "
		} else {
			cmd += " "
		}
		values = append(values, value)
	}
	cmd += m.querryKey
	stmt, err := m.db.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	return err
}

func (m *Model) Del() error {
	cmd := fmt.Sprintf("DELETE FROM %s ", m.table)
	cmd += m.querryKey
	stmt, err := m.db.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cmd)
	return err
}

func (m *Model) Get(content interface{}) error {
	if content == nil {
		return errors.New("参数不能为空")
	}
	m.content = content
	m.parseModel()
	cmd := fmt.Sprintf("SELECT 1 FROM %s", m.table)
	i := len(m.fieldMap)
	var values []interface{}
	for key, value := range m.fieldMap {
		i--
		cmd += key
		if i > 0 {
			cmd += ", "
		}
		values = append(values, value)
	}
	cmd += m.querryKey
	row := m.db.QueryRow(cmd)
	return row.Scan(values...)
}

func (m *Model) parseModel() {
	m.fieldMap = make(map[string]interface{})
	v := reflect.ValueOf(m.content)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get(TAG)
		value := v.Field(i).Interface()
		m.fieldMap[key] = value
	}
}

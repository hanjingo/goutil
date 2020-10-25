package redis

import (
	"errors"
)

const (
	STRING = "string"
	HASH   = "hash"
	LIST   = "list"
	SET    = "set"
	ZSET   = "zset"
)

type Model struct {
	pool     *RdsPool      //连接池
	dType    string        //数据结构
	isClosed bool          //是否已关闭
	keys     []interface{} //键
}

func NewModel(p *RdsPool, t string, keys ...interface{}) *Model {
	back := &Model{
		pool:     p,
		dType:    t,
		isClosed: false,
		keys:     keys,
	}
	return back
}

func (m *Model) Do(args ...interface{}) error {
	if m.isClosed {
		return errors.New("模型已关闭")
	}
	if args == nil || len(args) < 1 {
		return errors.New("参数错误")
	}
	cmd, ok := args[0].(string)
	if !ok {
		return errors.New("第一个参数必须为string")
	}
	in := args[1:]
	conn := m.pool.Get()
	defer m.pool.Put(conn)

	_, err := conn.Do(cmd, in...)
	return err
}

func (m *Model) Create(content interface{}) error {
	if m.isClosed {
		return errors.New("模型已关闭")
	}
	if content == nil {
		return errors.New("传入的参数不能为空")
	}
	conn := m.pool.Get()
	defer m.pool.Put(conn)

	var cmd []interface{}
	switch m.dType {
	case STRING:
		cmd = append(cmd, "SET")
	case HASH:
		cmd = append(cmd, "HSET")
		cmd = append(cmd, m.keys...)
	case LIST:
		cmd = append(cmd, "LPUSH")
		cmd = append(cmd, m.keys[0])
	case SET:
		cmd = append(cmd, "SADD")
		cmd = append(cmd, m.keys...)
	case ZSET:
		cmd = append(cmd, "ZADD")
		cmd = append(cmd, m.keys...)
	}
	return SaveItem(conn, content, cmd...)
}

func (m *Model) Update(content interface{}) error {
	if m.isClosed {
		return errors.New("模型已关闭")
	}
	if content == nil {
		return errors.New("对象内容不能为空")
	}
	conn := m.pool.Get()
	defer m.pool.Put(conn)

	var cmd []interface{}
	switch m.dType {
	case STRING:
		cmd = append(cmd, "SET")
	case HASH:
		cmd = append(cmd, "HSET")
		cmd = append(cmd, m.keys...)
	case LIST:
		if err := m.Del(); err != nil {
			return err
		}
		cmd = append(cmd, "LPUSH")
		cmd = append(cmd, m.keys[0])
	case SET:
		if err := m.Del(); err != nil {
			return err
		}
		cmd = append(cmd, "SADD")
		cmd = append(cmd, m.keys...)
	case ZSET:
		cmd = append(cmd, "")
	}
	return nil
}

func (m *Model) Del() error {
	if m.isClosed {
		return errors.New("模型已关闭")
	}
	if !m.isExist() {
		return errors.New("redis中不存在次数据")
	}
	var cmd []interface{}
	switch m.dType {
	case STRING:
		cmd = append(cmd, "DEL")
	case HASH:
		cmd = append(cmd, "HDEL")
	case LIST:
		cmd = append(cmd, "LPUSH")
	case SET:
		cmd = append(cmd, "SADD")
	case ZSET:
		cmd = append(cmd, "ZADD")
	}
	cmd = append(cmd, m.keys...)
	return nil
}

func (m *Model) Get(content interface{}) error {
	if m.isClosed {
		return errors.New("模型已关闭")
	}
	if content == nil {
		return errors.New("参数不能为空")
	}
	if !m.isExist() {
		return errors.New("数据不存在")
	}
	return nil
}

func (m *Model) Close() error {
	if m.isClosed {
		return errors.New("模型已关闭，无需再次关闭")
	}
	m.isClosed = true
	m.pool = nil
	m.keys = nil
	return nil
}

func (m *Model) isExist() bool {
	return false
}

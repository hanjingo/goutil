package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	rds "github.com/gomodule/redigo/redis"
)

//存
func SaveItem(conn RdsConn, target interface{}, args ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println("保存对象失败,错误:" + string(buf[:n]))
		}
	}()

	if len(args) < 1 || reflect.TypeOf(args[0]).Kind() != reflect.String {
		return errors.New(fmt.Sprintf("将对象保存到redis失败,错误:参数错误"))
	}
	data, err := json.Marshal(target)
	if err != nil || data == nil {
		return errors.New(fmt.Sprintf("序列化失败,错误:%v", err))
	}
	return doSave(conn, args[0].(string), append(args[1:], data)...)
}
func doSave(conn RdsConn, cmd string, args ...interface{}) error {
	if conn == nil || conn.Err() != nil {
		return errors.New(fmt.Sprintf("无法连接到redis"))
	}
	if _, err := conn.Do(cmd, args...); err != nil {
		return errors.New(fmt.Sprintf("redis保存信息失败,错误:%v", err))
	}
	return nil
}

//取
func GetItem(conn RdsConn, target interface{}, args ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("信息序列化读取redis失败,panic:" + string(buf[:n]))
		}
	}()

	if len(args) < 1 || reflect.TypeOf(args[0]).Kind() != reflect.String {
		return errors.New(fmt.Sprintf("将从redis读取对象失败,错误:参数错误"))
	}
	data, err := rds.Bytes(doGet(conn, args[0].(string), append(args[1:])...))
	if err != nil {
		return errors.New(fmt.Sprintf("将从redis读取对象失败,错误:", err))
	}
	if data == nil || len(data) == 0 {
		return errors.New(fmt.Sprintf("数据不存在"))
	}
	if err = json.Unmarshal(data, target); err != nil {
		return errors.New(fmt.Sprintf("从redis读取信息时反序列化失败"))
	}
	return nil
}
func doGet(conn RdsConn, cmd string, args ...interface{}) (interface{}, error) {
	if conn == nil || conn.Err() != nil {
		return nil, errors.New("无法连接到redis")
	}
	back, err := conn.Do(cmd, args...)
	return back, err
}

//删
func DelItem(conn RdsConn, args ...interface{}) error {
	if len(args) < 1 || reflect.TypeOf(args[0]).Kind() != reflect.String {
		return errors.New(fmt.Sprintf("将从redis删除对象失败,错误:参数错误"))
	}
	_, err := conn.Do(args[0].(string), args[1:]...)
	return err
}

//是否存在
func IsExist(conn RdsConn, cmd string, args ...interface{}) bool {
	if args == nil || len(args) < 1 {
		return false
	}
	if back, err := rds.Bool(conn.Do(cmd, args...)); err == nil {
		return back
	}
	return false
}

//发布
func PubRds(conn RdsConn, key string, msg []byte) error {
	if err := conn.Send("PUBLISH", key, msg); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

//订阅 阻塞
func SubRds(conn RdsConn, key string, subC chan []byte, ctx context.Context) error {
	psc := rds.PubSubConn{conn}
	if err := psc.Subscribe(key); err != nil {
		return err
	}
	defer psc.Unsubscribe(key)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		switch v := psc.Receive().(type) {
		case rds.Message:
			if v.Channel != key {
				continue
			}
			if v.Data == nil || len(v.Data) == 0 {
				continue
			}
			subC <- v.Data
		}
	}
}

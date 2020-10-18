package goutil

import (
	"reflect"

	glua "github.com/yuin/gopher-lua"
)

var state = glua.NewState()

//Lua2Go 转换 lua -> go: src:lub数据 args:[要转换到的目标， key转换选项]
func Lua2Go(src glua.LValue, typ reflect.Type) interface{} {
	switch v := src.(type) {
	case *glua.LNilType:
		return nil
	case glua.LBool:
		return bool(v)
	case glua.LString:
		return string(v)
	case glua.LNumber:
		switch typ.Kind() {
		case reflect.Int:
			return int(v)
		case reflect.Int8:
			return int8(v)
		case reflect.Int16:
			return int16(v)
		case reflect.Int32:
			return int32(v)
		case reflect.Int64:
			return int64(v)
		case reflect.Float32:
			return float32(v)
		case reflect.Float64:
			return float64(v)
		case reflect.Uint:
			return uint(v)
		case reflect.Uint8:
			return uint8(v)
		case reflect.Uint16:
			return uint16(v)
		case reflect.Uint32:
			return uint32(v)
		case reflect.Uint64:
			return uint64(v)
		}
		return float64(v)
	case *glua.LTable:
		tb := src.(*glua.LTable)
		switch typ.Kind() {
		case reflect.Struct:
			back := reflect.New(typ)
			value := back.Elem()
			for i := 0; i < typ.NumField(); i++ {
				tag := glua.LString(typ.Field(i).Tag)
				tmp := Lua2Go(tb.RawGet(tag), typ.Field(i).Type)
				value.Field(i).Set(reflect.ValueOf(tmp))
			}
			return back.Interface()
		case reflect.Ptr:
			return Lua2Go(src, typ.Elem())
		case reflect.Map:
			back := reflect.MakeMap(typ)
			typK := typ.Key()
			typV := typ.Elem()
			v.ForEach(func(key, value glua.LValue) {
				tmpK := Lua2Go(key, typK)
				tmpV := Lua2Go(value, typV)
				back.SetMapIndex(reflect.ValueOf(tmpK), reflect.ValueOf(tmpV))
			})
			return back.Interface()
		case reflect.Slice:
			maxn := v.MaxN()
			back := reflect.MakeSlice(typ, maxn, 2*maxn)
			typA := typ.Elem()
			for i := 1; i <= maxn; i++ {
				tmp := Lua2Go(v.RawGetInt(i), typA)
				back.Index(i - 1).Set(reflect.ValueOf(tmp))
			}
			return back.Interface()
		case reflect.Interface:
			return Lua2Go(src, typ.Elem())
		}
	}
	return src
}

//Go2Lua 转换 go -> lua
func Go2Lua(arg interface{}) glua.LValue {
	return go2Lua(reflect.ValueOf(arg))
}
func go2Lua(src reflect.Value) glua.LValue {
	typ := src.Type()
	switch src.Kind() {
	case reflect.Bool:
		return glua.LBool(src.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return glua.LNumber(src.Int())
	case reflect.Float32, reflect.Float64:
		return glua.LNumber(src.Float())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return glua.LNumber(src.Uint())
	case reflect.String:
		return glua.LString(src.String())
	case reflect.Struct:
		tb := state.NewTable()
		for i := 0; i < typ.NumField(); i++ {
			tmp := src.Field(i)
			tag := string(typ.Field(i).Tag)
			tb.RawSet(glua.LString(tag), go2Lua(tmp))
		}
		return tb
	case reflect.Ptr:
		value := src.Elem()
		if !value.IsValid() {
			return nil
		}
		return go2Lua(value)
	case reflect.Map:
		tb := state.NewTable()
		keys := src.MapKeys()
		for _, key := range keys {
			value := src.MapIndex(key)
			lkey := go2Lua(key)
			tb.RawSet(lkey, go2Lua(value))
		}
		return tb
	case reflect.Slice:
		tb := state.NewTable()
		for i := 0; i < src.Len(); i++ {
			value := go2Lua(src.Index(i))
			tb.RawSet(glua.LNumber(i+1), value)
		}
		return tb
	case reflect.Interface:
		value := src.Elem()
		if !value.IsValid() {
			return nil
		}
		return go2Lua(value)
	}
	return nil
}

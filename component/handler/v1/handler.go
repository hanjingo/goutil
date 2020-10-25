package v1_1

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	net "github.com/hanjingo/golib/network"
	p1 "github.com/hanjingo/golib/protocol/v1"
	types "github.com/hanjingo/golib/types"
)

type Handler struct {
	codec *p1.Codec
	fmap  map[uint32]func(net.SessionI, []byte) error
	buf   *bytes.Buffer
}

func NewHandler() *Handler {
	return &Handler{
		codec: p1.NewCodec(),
		fmap:  make(map[uint32]func(net.SessionI, []byte) error),
		buf:   new(bytes.Buffer),
	}
}

//注册回调函数，函数类型：func(session, req, rsp)
func (h *Handler) Reg(f interface{}, reqId uint32, req reflect.Type, rspId uint32, rsp reflect.Type) error {
	if _, ok := h.fmap[reqId]; ok {
		return errors.New(fmt.Sprintf("id:%v已存在", reqId))
	}
	h.fmap[reqId] = func(s net.SessionI, data []byte) error {
		in := []reflect.Value{reflect.ValueOf(s)}

		if req != nil {
			reqMsg := &p1.Messager{Id: reqId, Content: types.New(req)}
			if err := h.codec.UnFormat(data, reqMsg); err != nil {
				return err
			}
			in = append(in, reflect.ValueOf(reqMsg.Content))
		}

		var rspContent interface{}
		if rsp != nil {
			rspContent = types.New(rsp)
			in = append(in, reflect.ValueOf(rspContent))
		}

		fv := reflect.ValueOf(f)
		fv.Call(in)

		if rsp != nil {
			rspMsg := &p1.Messager{Id: rspId, Content: rspContent}
			back, err := h.codec.Format(rspMsg)
			if err != nil {
				return err
			}
			if _, err := s.WriteMsg(back); err != nil {
				return err
			}
		}

		return nil
	}
	return nil
}

func (h *Handler) Handle(s net.SessionI, n int) {
	tmp := make([]byte, n)
	_, err := s.ReadMsg(tmp)
	if err != nil {
		return
	}
	h.buf.Write(tmp)
	if h.buf.Len() < p1.MsgHeadLen {
		return
	}
	length, err := h.codec.ParseLen(h.buf.Bytes())
	if err != nil {
		return
	}
	id, err := h.codec.ParseId(h.buf.Bytes())
	if err != nil {
		return
	}
	f, ok := h.fmap[id]
	if !ok {
		return
	}
	if h.buf.Len() < length {
		return
	}
	if err := f(s, h.buf.Next(length)); err != nil {
		return
	}
}

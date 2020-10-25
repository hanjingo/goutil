package v2

import (
	"bytes"
	"errors"

	"github.com/hanjingo/golib/network"
)

type Handler struct {
	buf      *bytes.Buffer
	doHandle func(s network.SessionI, handler *Handler, n int)
}

func NewHandler(h func(s network.SessionI, handler *Handler, n int)) *Handler {
	return &Handler{
		buf:      new(bytes.Buffer),
		doHandle: h,
	}
}

func (h *Handler) Handle(s network.SessionI, n int) {
	if h.doHandle == nil {
		return
	}
	h.doHandle(s, h, n)
}

func (h *Handler) Read(n int) ([]byte, error) {
	if h.buf.Len() < n {
		return nil, errors.New("可读长度不够")
	}
	back := make([]byte, n)
	length, err := h.buf.Read(back)
	if err != nil {
		return nil, err
	}
	if length != n {
		return back[:length], errors.New("读到的长度不够")
	}
	return back, nil
}

package v1_1

import (
	"reflect"
	"testing"
)

type Req struct {
	Content string
}

type Rsp struct {
	Content string
}

func doHello(req *Req, rsp *Rsp) {
	rsp.Content = "world"
}

func TestNewHandler(t *testing.T) {
	h := NewHandler()
	h.Reg(1, doHello, reflect.TypeOf(&Req), reflect.TypeOf(&Rsp))
}

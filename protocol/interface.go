package protocol

import (
	p1 "github.com/hanjingo/golib/protocol/v1"
	p2 "github.com/hanjingo/golib/protocol/v2"
	p3 "github.com/hanjingo/golib/protocol/v3"
	p4 "github.com/hanjingo/golib/protocol/v4"
)

type CodecI interface {
	Version() string                    //版本
	Format(interface{}) ([]byte, error) //格式化
	UnFormat([]byte, interface{}) error //反格式化
}

var VERSION1 = p1.VERSION
var VERSION2 = p2.VERSION
var VERSION3 = p3.VERSION
var VERSION4 = p4.VERSION
var VERSION_TRSP = rtsp.VERSION

var NewCodec1 = p1.NewCodec
var NewCodec2 = p2.NewCodec
var NewCodec3 = p3.NewCodec
var NewCodec4 = p4.NewCodec
var NewRtspCodec = rtsp.NewCodec

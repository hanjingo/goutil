package v2

import (
	"encoding/json"
	"errors"

	types "github.com/hanjingo/golib/types"
)

const VERSION = "2.0.0"
const MsgLenLen int = 2

/*
	版本: 		2.0.0
	序列化:		json
	格式:		len(2) + content(0~65535)
*/

//消息体
type Messager struct {
	Content interface{} //消息内容
}

func (m *Messager) Version() string {
	return VERSION
}

//编解码
type Codec struct{}

func NewCodec() *Codec {
	return &Codec{}
}

//版本
func (c *Codec) Version() string {
	return VERSION
}

//格式化
func (c *Codec) Format(arg interface{}) ([]byte, error) {
	if _, ok := arg.(*Messager); !ok {
		return nil, errors.New("参数不对")
	}
	msg := arg.(*Messager)
	if msg.Version() != c.Version() {
		return nil, errors.New("消息和编码器版本不一致")
	}
	var data []byte
	data2, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	length := MsgLenLen + len(data2)
	data1, err := types.UintToBytes(uint16(length))
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)
	data = append(data, data2...)

	return data, nil
}

//反格式化
func (c *Codec) UnFormat(data []byte, arg interface{}) error {
	if _, ok := arg.(*Messager); !ok {
		return errors.New("参数不对")
	}
	msg := arg.(*Messager)
	if msg.Version() != c.Version() {
		return errors.New("消息和解码器版本不一致")
	}

	if data == nil {
		return errors.New("data不能为空")
	}
	length, err := c.ParseLen(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data[MsgLenLen:length], msg.Content); err != nil {
		return err
	}
	return nil
}

//解析长度
func (c *Codec) ParseLen(data []byte) (int, error) {
	if len(data) < MsgLenLen {
		return 0, errors.New("消息过短")
	}
	length, err := types.BytesToUint(data[:MsgLenLen])
	if err != nil {
		return 0, err
	}
	return int(length.(uint16)), nil
}

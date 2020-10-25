package v1

import (
	"encoding/json"
	"errors"

	types "github.com/hanjingo/golib/types"
)

const VERSION = "1.0.0"
const MsgIdLen int = 4
const MsgLenLen int = 2
const MsgHeadLen int = 6

/*
	版本: 			1.0.0
	id长度:			4位
	内容长度:		2位,Content(0~65525位)
	序列化:			json
	格式:			id(4)+len(2)+content(0~65535)
*/

//消息体
type Messager struct {
	Content interface{} //内容
	Id      uint32      //Id
}

func (m *Messager) Version() string {
	return VERSION
}

//编解码器
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
	data3, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, err
	}

	data1, err := types.UintToBytes(msg.Id)
	if err != nil {
		return nil, err
	}
	data = append(data, data1...)

	length := MsgHeadLen + len(data3)
	data2, err := types.UintToBytes(uint16(length))
	if err != nil {
		return nil, err
	}
	data = append(data, data2...)

	data = append(data, data3...)
	return data, nil
}

//反格式化
func (c *Codec) UnFormat(data []byte, arg interface{}) error {
	if _, ok := arg.(*Messager); !ok {
		return errors.New("参数不对")
	}
	msg := arg.(*Messager)
	if msg.Version() != c.Version() {
		return errors.New("消息和编码器版本不一致")
	}

	if len(data) < MsgHeadLen {
		return errors.New("消息过短")
	}

	var err error
	//Id
	if msg.Id, err = c.ParseId(data); err != nil {
		return err
	}

	//len
	length, err := c.ParseLen(data)
	if err != nil {
		return err
	}
	if length < MsgHeadLen || length > len(data) {
		return errors.New("消息不完整")
	}

	//Content
	if err := json.Unmarshal(data[MsgHeadLen:length], msg.Content); err != nil {
		return err
	}
	return nil
}

//解析id
func (c *Codec) ParseId(data []byte) (uint32, error) {
	if len(data) < MsgIdLen {
		return 0, errors.New("消息过短")
	}
	Id, err := types.BytesToUint(data[:MsgIdLen])
	if err != nil {
		return 0, err
	}
	return Id.(uint32), nil
}

//解析长度
func (c *Codec) ParseLen(data []byte) (int, error) {
	if len(data) < MsgHeadLen {
		return 0, errors.New("消息过短")
	}
	length, err := types.BytesToUint(data[MsgIdLen:MsgHeadLen])
	if err != nil {
		return 0, err
	}
	return int(length.(uint16)), nil
}

//解析content
func (c *Codec) ParseContent(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("消息为空")
	}
	length, err := c.ParseLen(data)
	if err != nil {
		return nil, err
	}
	if len(data) < length {
		return nil, errors.New("消息不完整")
	}
	return data[MsgHeadLen:length], nil
}

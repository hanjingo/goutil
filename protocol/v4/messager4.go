package v4

import (
	"encoding/json"
	"errors"

	types "github.com/hanjingo/golib/types"
)

const VERSION string = "4.0.0"
const MsgLenLen int = 2
const MsgOpcodeLen int = 4
const MsgSenderLen int = 8
const MsgRecvLenLen int = 2

/*
	版本:			"4.0.0"
	序列化:			json
	格式:			消息总长(2)+操作码(4)+收信人名单长度(2)+发送者(8)+收信人名单(0~65535)+内容(0~65535)
*/

//消息体
type Messager struct {
	OpCode   uint32
	Receiver []uint64
	Sender   uint64
	Content  interface{}
}

//版本
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

	//操作码
	data2, err := types.UintToBytes(msg.OpCode)
	if err != nil {
		return nil, err
	}

	//收信人
	var data5 []byte
	for _, recv := range msg.Receiver {
		temp, err := types.UintToBytes(recv)
		if err != nil {
			return nil, err
		}
		data5 = append(data5, temp...)
	}

	//收信人名单长度
	data3, err := types.UintToBytes(uint16(len(data5)))
	if err != nil {
		return nil, err
	}

	//送信人
	data4, err := types.UintToBytes(msg.Sender)
	if err != nil {
		return nil, err
	}

	//内容
	data6, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, err
	}

	//消息总长
	length := uint16(MsgLenLen + MsgOpcodeLen + MsgSenderLen + MsgRecvLenLen + len(data5) + len(data6))
	data1, err := types.UintToBytes(length)
	if err != nil {
		return nil, err
	}

	var data []byte
	data = append(data, data1...)
	data = append(data, data2...)
	data = append(data, data3...)
	data = append(data, data4...)
	data = append(data, data5...)
	data = append(data, data6...)
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
	var err error
	if msg.OpCode, err = c.ParseOpCode(data); err != nil {
		return err
	}
	if msg.Sender, err = c.ParseSender(data); err != nil {
		return err
	}
	if msg.Receiver, err = c.ParseRecv(data); err != nil {
		return err
	}
	if err = c.ParseContent(data, msg.Content); err != nil {
		return err
	}
	return nil
}

//解析总长
func (c *Codec) ParseTotalLen(data []byte) (uint16, error) {
	if data == nil || len(data) < MsgLenLen {
		return 0, errors.New("消息过短")
	}
	length, err := types.BytesToUint(data[:MsgLenLen])
	if err != nil {
		return 0, err
	}
	return length.(uint16), nil
}

//解析opcode
func (c *Codec) ParseOpCode(data []byte) (uint32, error) {
	start := MsgLenLen
	end := start + MsgOpcodeLen
	if data == nil || len(data) < end {
		return 0, errors.New("消息过短")
	}
	opcode, err := types.BytesToUint(data[start:end])
	if err != nil {
		return 0, err
	}
	return opcode.(uint32), nil
}

//解析收信人名单长度
func (c *Codec) ParseRecvLen(data []byte) (uint16, error) {
	start := MsgLenLen + MsgOpcodeLen
	end := start + MsgRecvLenLen
	if data == nil || len(data) < end {
		return 0, errors.New("消息过短")
	}
	recvLen, err := types.BytesToUint(data[start:end])
	if err != nil {
		return 0, err
	}
	return recvLen.(uint16), nil
}

//解析发送者
func (c *Codec) ParseSender(data []byte) (uint64, error) {
	start := MsgLenLen + MsgOpcodeLen + MsgRecvLenLen
	end := start + MsgSenderLen
	if data == nil || len(data) < end {
		return 0, errors.New("消息过短")
	}
	sender, err := types.BytesToUint(data[start:end])
	if err != nil {
		return 0, err
	}
	return sender.(uint64), nil
}

//解析收信人名单
func (c *Codec) ParseRecv(data []byte) ([]uint64, error) {
	start := MsgLenLen + MsgOpcodeLen + MsgRecvLenLen + MsgSenderLen
	recvLen, err := c.ParseRecvLen(data)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data) < int(recvLen) {
		return nil, errors.New("消息过短")
	}
	var back []uint64
	i := start
	j := start + 8
	end := int(start) + int(recvLen)
	for j <= end {
		recv, err := types.BytesToUint(data[i:j])
		if err != nil {
			return nil, err
		}
		back = append(back, recv.(uint64))
		i = j
		j += 8
	}
	return back, nil
}

//解析内容
func (c *Codec) ParseContent(data []byte, content interface{}) error {
	temp := c.GetContentData(data)
	if temp == nil {
		return errors.New("截取内容部分失败")
	}
	if err := json.Unmarshal(temp, content); err != nil {
		return err
	}
	return nil
}

//获得内容二进制
func (c *Codec) GetContentData(data []byte) []byte {
	recvLen, err := c.ParseRecvLen(data)
	if err != nil {
		return nil
	}
	start := MsgLenLen + MsgOpcodeLen + MsgRecvLenLen + MsgSenderLen + int(recvLen)
	totalLen, err := c.ParseTotalLen(data)
	if err != nil {
		return nil
	}
	if data == nil || len(data) < int(totalLen) {
		return nil
	}
	return data[start:totalLen]
}

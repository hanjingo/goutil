package v4

import (
	"fmt"
	"testing"
)

type hello struct {
	CMD string
}

func TestFormat(t *testing.T) {
	c := NewCodec()
	msg := &Messager{
		OpCode:   1,
		Receiver: []uint64{1, 2},
		Sender:   1,
		Content:  &hello{CMD: "007"},
	}
	data, err := c.Format(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg1 := &Messager{
		Content: &hello{},
	}
	n1, err := c.ParseTotalLen(data)
	fmt.Println("totalLen>>", n1, " err>>", err)

	n2, err := c.ParseOpCode(data)
	fmt.Println("opcode>>", n2, " err>>", err)

	n3, err := c.ParseRecv(data)
	fmt.Println("recvs>>", n3, " err>>", err)

	n5, err := c.ParseRecvLen(data)
	fmt.Println("recvLen>>", n5, " err>>", err)

	n6, err := c.ParseSender(data)
	fmt.Println("sender>>", n6, " err>>", err)

	err = c.UnFormat(data, msg1)
	if err != nil {
		fmt.Println(err)
	}
}

func TestUnFormat(t *testing.T) {

}

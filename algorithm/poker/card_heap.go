package poker

import (
	"github.com/hanjingo/golib/container"
)

type CardHeap struct {
	mi *container.MultiIndexTable //多索引表
}

func (ch *CardHeap) AddCard(cards ...PokerCard) {
}

package poker

//一张牌
type PokerCard interface {
	Point() byte //点数
	Color() byte //花色
}

//牌堆检查函数
type CheckF func(cards ...PokerCard) bool

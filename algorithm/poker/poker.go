package poker

//是否花色一样
func IsColorEqual(cards ...PokerCard) bool {
	if cards == nil || len(cards) < 2 {
		return false
	}
	for i := 0; i < len(cards); i++ {
		if !IsColorEqual(cards[0], cards[i]) {
			return false
		}
	}
	return true
}

//是否点数一样
func IsPointEqual(cards ...PokerCard) bool {
	if cards == nil || len(cards) < 2 {
		return false
	}
	for i := 0; i < len(cards); i++ {
		if !IsPointEqual(cards[0], cards[i]) {
			return false
		}
	}
	return true
}

//获得点数最大的牌
func MaxPointCard(cards ...PokerCard) PokerCard {
	if cards == nil || len(cards) == 0 {
		return nil
	}
	var back PokerCard = cards[0]
	for _, v := range cards {
		if back.Point() > v.Point() {
			back = v
		}
	}
	return back
}

//获得点数最小的牌
func MinPointCard(cards ...PokerCard) PokerCard {
	if cards == nil || len(cards) == 0 {
		return nil
	}
	var back PokerCard = cards[0]
	for _, v := range cards {
		if back.Point() < v.Point() {
			back = v
		}
	}
	return back
}

//按点数从小到大排序
func AscendingByPoint(cards ...PokerCard) []PokerCard {
	if cards == nil || len(cards) < 2 {
		return cards
	}
	for i := 0; i < len(cards); i++ {
		for j := i; j < len(cards); j++ {
			if cards[i].Point() > cards[j].Point() {
				//swap
				temp := cards[j]
				cards[i] = cards[j]
				cards[j] = temp
			}
		}
	}
	return cards
}

//按点数从大到小排序
func DescendingByPoint(cards ...PokerCard) []PokerCard {
	if cards == nil || len(cards) < 2 {
		return cards
	}
	for i := 0; i < len(cards); i++ {
		for j := i; j < len(cards); j++ {
			if cards[i].Point() < cards[j].Point() {
				//swap
				temp := cards[j]
				cards[i] = cards[j]
				cards[j] = temp
			}
		}
	}
	return cards
}

//获得点数相同牌集合(注意剔除牌数为1的item)  返回 key:point value:牌集合
func GetSamePointCardsIndex(cards ...PokerCard) map[byte][]PokerCard {
	if cards == nil || len(cards) < 2 {
		return nil
	}
	back := make(map[byte][]PokerCard)
	for i := 0; i < len(cards); i++ {
		back[cards[i].Point()] = []PokerCard{}
	}
	for point, _ := range back {
		for j := 0; j < len(cards); j++ {
			if cards[j].Point() == point {
				back[point] = append(back[point], cards[j])
			}
		}
	}
	return back
}

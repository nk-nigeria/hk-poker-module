package hand

import (
	"fmt"
	"sort"

	"github.com/nakamaFramework/cgp-chinese-poker-module/entity"
)

type autoHand struct {
	cards         entity.ListCard
	cardsCount    map[entity.Rank]int
	trackCardTake map[entity.Card]int
}

func NewAutoHand(cards entity.ListCard) *autoHand {
	h := &autoHand{
		cards:         cards.Clone(),
		cardsCount:    make(map[entity.Rank]int),
		trackCardTake: make(map[entity.Card]int),
	}
	sort.Slice(h.cards, func(i, j int) bool {
		return h.cards[i].GetRank() > h.cards[j].GetRank()
	})
	// tần suất lá bài xuất hiện
	for _, c := range h.cards {
		h.cardsCount[c.GetRank()]++
	}
	// xóa lá bài chỉ xuất hiện 1 lần
	for k, v := range h.cardsCount {
		if v <= 1 {
			delete(h.cardsCount, k)
		}
	}
	return h
}

func (h *autoHand) FindStraighFlush() []entity.ListCard {
	cardsSameColor := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameColor[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameColor[c.GetRank()] = list
	}
	// remove card same color <5 card
	listStraighFlush := make([]entity.ListCard, 0)
	for _, v := range cardsSameColor {
		if len(v) < 5 {
			continue
		}
		for i := 0; i <= len(v)-5; i++ {
			ml := v[i : 5+i]
			binCards := entity.NewBinListCards(ml)
			_, ok := CheckStraightFlush(binCards)
			if !ok {
				continue
			}
			listStraighFlush = append(listStraighFlush, ml.Clone())
		}
	}
	return listStraighFlush
}

func (h *autoHand) FindStraigh() []entity.ListCard {
	cardsSameRank := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameRank[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameRank[c.GetRank()] = list
	}
	// remove card same color <5 card
	listStraigh := make([]entity.ListCard, 0)
	for _, v := range cardsSameRank {
		if len(v) < 5 {
			continue
		}
		for i := 0; i <= len(v)-5; i++ {
			ml := v[i : 5+i]
			listStraigh = append(listStraigh, ml.Clone())
		}
	}
	return listStraigh
}
func (h *autoHand) FindFullHouse() []entity.ListCard {
	threeKinds := h.FindThreeKind()
	// tempTrack := make(map[entity.Card]struct{})
	// for _, cards := range threeKinds {
	// 	for _, card := range cards {
	// 		tempTrack[card] = struct{}{}
	// 	}
	// }
	for _, threeKind := range threeKinds {
		doubles := h.FindPair()
		for i := len(doubles) - 1; i >= 0; i-- {
			double := doubles[i]
			if double[0].GetRank() == threeKind[0].GetRank() {
				continue
			}
			list := append(threeKind, double[:2]...)
			return []entity.ListCard{list}
		}
	}
	return nil
}
func (h *autoHand) FindFlush() []entity.ListCard {
	cardsSameSuit := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameSuit[c.GetSuit()]
		v := c
		list = append(list, v)
		cardsSameSuit[c.GetSuit()] = list
	}
	// remove card same suit <5 card
	listFlush := make([]entity.ListCard, 0)
	for _, v := range cardsSameSuit {
		if len(v) < 5 {
			continue
		}
		for i := 0; i < len(v)-5; i++ {
			ml := v[:5+i]
			binCards := entity.NewBinListCards(ml)
			_, ok := CheckFlush(binCards)
			if !ok {
				continue
			}
			listFlush = append(listFlush, ml.Clone())
		}
	}
	return listFlush
}

func (h *autoHand) FindFourKind() []entity.ListCard {
	cardsSameRank := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameRank[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameRank[c.GetRank()] = list
	}
	list := make([]entity.ListCard, 0)
	for _, v := range cardsSameRank {
		if len(v) < 4 {
			continue
		}
		list = append(list, v.Clone())
	}
	return list
}

func (h *autoHand) FindThreeKind() []entity.ListCard {
	cardsSameRank := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameRank[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameRank[c.GetRank()] = list
	}
	list := make([]entity.ListCard, 0)
	for _, v := range cardsSameRank {
		if len(v) < 3 {
			continue
		}
		list = append(list, v.Clone())
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i][0].GetRank() > list[j][0].GetRank()
	})
	return list
}

func (h *autoHand) FindTwoPair() []entity.ListCard {
	cardsSameRank := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameRank[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameRank[c.GetRank()] = list
	}
	list := make([]entity.ListCard, 0)
	for _, v := range cardsSameRank {
		if len(v) < 2 {
			continue
		}
		list = append(list, v.Clone())
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i][0].GetRank() > list[j][0].GetRank()
	})
	if len(list) >= 2 {
		return list
	}
	return nil
}
func (h *autoHand) FindPair() []entity.ListCard {
	cardsSameRank := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}
		list := cardsSameRank[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameRank[c.GetRank()] = list
	}
	list := make([]entity.ListCard, 0)
	for _, v := range cardsSameRank {
		if len(v) < 2 {
			continue
		}
		list = append(list, v.Clone())
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i][0].GetRank() > list[j][0].GetRank()
	})
	return list
}

func (h *autoHand) FindHighCard() entity.ListCard {
	cardsSameRank := make(map[entity.Rank]entity.ListCard)
	for _, c := range h.cards {
		if h.trackCardTake[c] > 0 {
			continue
		}

		v := c
		list := cardsSameRank[c.GetRank()]
		list = append(list, v)
		cardsSameRank[c.GetRank()] = list
	}
	list := make(entity.ListCard, 0)
	for _, v := range cardsSameRank {
		if len(v) >= 2 {
			continue
		}
		list = append(list, v...)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetRank() > list[j].GetRank()
	})
	return list
}
func (h *autoHand) TakeCard(cards ...entity.Card) {
	for _, card := range cards {
		c := card
		h.trackCardTake[c]++
		if h.trackCardTake[c] > 1 {
			fmt.Println("hoooo")
		}
	}
}
func (h *autoHand) PreferCardsNotTakeDouble(ml ...entity.ListCard) entity.ListCard {
	if len(ml) == 0 {
		return nil
	}
	if len(ml) == 1 {
		return ml[0]
	}
	// ưu tiên straigh flush ko chứa card tạo đôi
	listMin := ml[0]
	minCount := 0
	for _, list := range ml {
		count := 0
		for _, c := range list {
			if h.cardsCount[c.GetRank()] >= 2 {
				count += h.cardsCount[c.GetRank()]
			}
		}
		if count == 0 {
			return list
		}
		if count < minCount {
			listMin = list
			minCount = count
		}
	}
	return listMin
}

func (h *autoHand) GetHighCardWithCond(num int, excludeCards ...entity.Card) entity.ListCard {
	if num == 0 {
		return nil
	}
	excludeCardsMap := make(map[entity.Rank]struct{})
	for _, c := range excludeCards {
		excludeCardsMap[c.GetRank()] = struct{}{}
	}
	highcard := make(entity.ListCard, 0)
	for _, card := range h.FindHighCard() {
		if _, ok := excludeCardsMap[card.GetRank()]; ok {
			continue
		}
		highcard = append(highcard, card)
	}
	if len(highcard) < num {
		for _, pair := range h.FindPair() {
			for _, card := range pair {
				if _, ok := excludeCardsMap[card.GetRank()]; ok {
					continue
				}
				highcard = append(highcard, card)
			}
		}
	}
	sort.Slice(highcard, func(i, j int) bool {
		return highcard[i].GetRank() < highcard[j].GetRank()
	})
	if len(highcard) >= num {
		return highcard[:num]
	}
	return highcard
}

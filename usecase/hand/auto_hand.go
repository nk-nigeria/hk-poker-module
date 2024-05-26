package hand

import (
	"sort"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
)

type AutoHand struct {
	cards         entity.ListCard
	cardsCount    map[uint8]int
	trackCardTake map[entity.Card]struct{}
}

func NewAutoHand(cards entity.ListCard) *AutoHand {

	h := &AutoHand{
		cards:         cards.Clone(),
		cardsCount:    make(map[uint8]int),
		trackCardTake: make(map[entity.Card]struct{}),
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

func (h *AutoHand) FindStraighFlush() []entity.ListCard {
	cardsSameColor := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
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
		for i := 0; i < len(v)-5; i++ {
			ml := v[:5+i]
			binCards := entity.NewBinListCards(ml)
			_, ok := CheckStraightFlush(binCards)
			if !ok {
				continue
			}
			listStraighFlush = append(listStraighFlush, ml)
		}
	}
	return listStraighFlush
}

func (h *AutoHand) FindStraigh() []entity.ListCard {
	cardsSameColor := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
			continue
		}
		list := cardsSameColor[c.GetRank()]
		v := c
		list = append(list, v)
		cardsSameColor[c.GetRank()] = list
	}
	// remove card same color <5 card
	listStraigh := make([]entity.ListCard, 0)
	for _, v := range cardsSameColor {
		if len(v) < 5 {
			continue
		}
		listStraigh = append(listStraigh, v)
	}
	return listStraigh
}
func (h *AutoHand) FindFullHouse() []entity.ListCard {
	threeKinds := h.FindThreeKind()
	tempTrack := make(map[entity.Card]struct{})
	for _, cards := range threeKinds {
		for _, card := range cards {
			tempTrack[card] = struct{}{}
		}
	}
	for _, threeKind := range threeKinds {
		doubles := h.FindPair()
		for _, double := range doubles {
			if double[0].GetRank() == threeKind[0].GetRank() {
				continue
			}
			list := append(threeKind, double...)
			return []entity.ListCard{list}
		}
	}
	return nil
}
func (h *AutoHand) FindFlush() []entity.ListCard {
	cardsSameSuit := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
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
			listFlush = append(listFlush, ml)
		}
	}
	return listFlush
}

func (h *AutoHand) FindFourKind() []entity.ListCard {
	cardsSameRank := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
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
		list = append(list, v)
	}
	return list
}

func (h *AutoHand) FindThreeKind() []entity.ListCard {
	cardsSameRank := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
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
		list = append(list, v)
	}
	return list
}

func (h *AutoHand) FindTwoPair() []entity.ListCard {
	cardsSameRank := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
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
		list = append(list, v)
		if len(list) == 2 {
			return list
		}
	}
	return nil
}
func (h *AutoHand) FindPair() []entity.ListCard {
	cardsSameRank := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
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
		list = append(list, v)
	}
	return list
}

func (h *AutoHand) FindHighCard() entity.ListCard {
	cardsSameRank := make(map[uint8]entity.ListCard)
	for _, c := range h.cards {
		if _, exist := h.trackCardTake[c]; exist {
			continue
		}
		list := cardsSameRank[c.GetRank()]
		v := c
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
	return list
}
func (h *AutoHand) TakeCard(cards ...entity.Card) {
	for _, card := range cards {
		h.trackCardTake[card] = struct{}{}
	}
}
func (h *AutoHand) PreferCardsNotTakeDouble(ml ...entity.ListCard) entity.ListCard {
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

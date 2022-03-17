package api

import (
	"sort"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

type HandPoint struct {
	rankingType pb.HandRanking
	point       int
}

type HandCards struct {
	ListCard    entity.ListCard
	MapCardType map[pb.HandRanking]entity.ListCard
}

func NewHandCards() *HandCards {
	return &HandCards{
		MapCardType: make(map[pb.HandRanking]entity.ListCard),
	}
}

func (hc *HandCards) CopyMapCardType(mapCardType map[pb.HandRanking]entity.ListCard) {
	for k, v := range mapCardType {
		l := hc.MapCardType[k]
		l = append(l, v...)
		hc.MapCardType[k] = l
	}
}

func (h *HandPoint) IsStraight() bool {
	return h.rankingType == pb.HandRanking_Straight
}

func (h *HandPoint) IsFlush() bool {
	return h.rankingType == pb.HandRanking_Flush
}

type CheckFunc func(entity.ListCard) (*HandCards, bool)

var HandChecker *linkedhashmap.Map

var HandCheckerFront *linkedhashmap.Map

func init() {
	HandChecker = linkedhashmap.New()
	HandChecker.Put(pb.HandRanking_StraightFlush, CheckStraightFlush)
	HandChecker.Put(pb.HandRanking_FourOfAKind, CheckFourOfAKind)
	HandChecker.Put(pb.HandRanking_FullHouse, CheckFullHouse)
	HandChecker.Put(pb.HandRanking_Flush, CheckFlush)
	HandChecker.Put(pb.HandRanking_Straight, CheckStraight)
	HandChecker.Put(pb.HandRanking_ThreeOfAKind, CheckThreeOfAKind)
	HandChecker.Put(pb.HandRanking_TwoPairs, CheckTwoPairs)
	HandChecker.Put(pb.HandRanking_Pair, CheckPair)

	HandCheckerFront = linkedhashmap.New()
	HandCheckerFront.Put(pb.HandRanking_ThreeOfAKind, CheckThreeOfAKind)
	HandCheckerFront.Put(pb.HandRanking_Pair, CheckPair)
}

// var HandChecker = map[pb.HandRanking]CheckFunc{
// 	pb.HandRanking_StraightFlush: CheckStraightFlush,
// 	pb.HandRanking_FourOfAKind:   CheckFourOfAKind,
// 	pb.HandRanking_FullHouse:     CheckFullHouse,
// 	pb.HandRanking_Flush:         CheckFlush,
// 	pb.HandRanking_Straight:      CheckStraight,
// 	pb.HandRanking_ThreeOfAKind:  CheckThreeOfAKind,
// 	pb.HandRanking_TwoPairs:      CheckTwoPairs,
// 	pb.HandRanking_Pair:          CheckPair,
// }

// var HandCheckerFront = map[pb.HandRanking]CheckFunc{
// 	pb.HandRanking_ThreeOfAKind: CheckThreeOfAKind,
// 	pb.HandRanking_Pair:         CheckPair,
// }

func CaculatorPoint(listCard entity.ListCard) (*HandPoint, *HandCards) {
	if len(listCard) == 3 {
		// For check front
		for _, key := range HandCheckerFront.Keys() {
			rank := key.(pb.HandRanking)
			val, _ := (HandCheckerFront.Get(rank))
			fn := val.(func(listCard entity.ListCard) (*HandCards, bool))
			if handCard, valid := fn(listCard); valid {
				return &HandPoint{
					rankingType: rank,
				}, handCard
			}
		}
	} else {
		for _, key := range HandChecker.Keys() {
			rank := key.(pb.HandRanking)
			val, exist := HandChecker.Get(rank)
			if !exist {
				return nil, nil
			}
			fn := val.(func(listCard entity.ListCard) (*HandCards, bool))
			if handCard, valid := fn(listCard); valid {
				return &HandPoint{
					rankingType: rank,
				}, handCard
			}
		}
	}

	// Sort for high card
	listCard = SortCard(listCard)
	handCard := HandCards{
		ListCard:    listCard,
		MapCardType: make(map[pb.HandRanking]entity.ListCard),
	}
	return &HandPoint{
		rankingType: pb.HandRanking_HighCard,
	}, &handCard
}

// CheckStraightFlush
// Thùng phá sảnh (en: Straight Flush)
// Năm lá bài cùng màu, đồng chất, cùng một chuỗi số
// Là Flush, có cùng chuỗi
func CheckStraightFlush(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	handCard, valid := CheckFlush(listCard)
	if !valid {
		return nil, false
	}
	handCard2, valid := CheckStraight(listCard)
	if !valid {
		return nil, false
	}
	if handCard.MapCardType == nil {
		handCard.MapCardType = make(map[pb.HandRanking]entity.ListCard)
	}
	for k, v := range handCard2.MapCardType {
		handCard.MapCardType[k] = v
	}
	handCard.MapCardType[pb.HandRanking_StraightFlush] = listCard
	handCard.ListCard = handCard2.ListCard
	return handCard, true
}

// CheckFourOfAKind
// Tứ quý (en: Four of a Kind)
// Bốn lá đồng số
func CheckFourOfAKind(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)

	newListCard := make([]entity.Card, 0, len(listCard))
	handCard := NewHandCards()
	isFourOfAKind := false
	var list entity.ListCard
	var remain entity.ListCard
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 4 {
			isFourOfAKind = true

			handCard.MapCardType[pb.HandRanking_FourOfAKind] = list
			newListCard = append(newListCard, list...)
		} else {
			remain = append(remain, list...)
		}
	}
	if isFourOfAKind {
		remain = SortCard(remain)
		handCard.ListCard = append(handCard.ListCard, newListCard...)
		handCard.ListCard = append(handCard.ListCard, remain...)
		return handCard, true
	}
	return nil, false
}

// CheckFullHouse
// Cù lũ (en: Full House)
// Một bộ ba và một bộ đôi
// Bốn lá đồng số
func CheckFullHouse(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() != 2 {
		return nil, false
	}

	newListCard := make(entity.ListCard, 0, len(listCard))
	hasTriangle := false
	hasDouble := false

	var list entity.ListCard
	handCard := NewHandCards()
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			handCard.MapCardType[pb.HandRanking_ThreeOfAKind] = list
			continue
		}
		if len(list) == 2 {
			hasDouble = true
			newListCard = append(list, newListCard...)
			handCard.MapCardType[pb.HandRanking_Pair] = list

			continue
		}
		newListCard = append(newListCard, list...)
	}
	if hasTriangle && hasDouble {
		newListCard = SortCard(newListCard)
		handCard.ListCard = newListCard
		return handCard, true
	}
	return nil, false
}

// CheckFlush
// Thùng (en: Flush)
// Năm lá bài cùng màu, đồng chất (nhưng không cùng một chuỗi số)
func CheckFlush(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard = SortCard(listCard)

	prevSuitPoint := listCard[0].GetSuit()
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		suitPoint := card.GetSuit()
		if suitPoint != prevSuitPoint {
			return nil, false
		}
	}
	handCard := HandCards{
		ListCard:    listCard,
		MapCardType: make(map[pb.HandRanking]entity.ListCard),
	}
	handCard.MapCardType[pb.HandRanking_Flush] = listCard
	return &handCard, true
}

// CheckStraight
// Sảnh (en: Straight)
// Năm lá bài trong một chuỗi số (nhưng không đồng chất)
func CheckStraight(listCard entity.ListCard) (*HandCards, bool) {
	listCard = SortCard(listCard)
	prevRankPoint := listCard[0].GetRank()
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		rankPoint := card.GetRank()
		prevRankPoint--
		if rankPoint != prevRankPoint {
			return nil, false
		}
	}
	handCards := HandCards{
		ListCard:    listCard,
		MapCardType: make(map[pb.HandRanking]entity.ListCard),
	}
	handCards.MapCardType[pb.HandRanking_Straight] = listCard
	return &handCards, true
}

// CheckThreeOfAKind
// Xám chi/Xám cô (en: Three of a Kind)
// Ba lá bài đồng số
func CheckThreeOfAKind(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 2 {
		return nil, false
	}

	newListCard := make(entity.ListCard, 0, len(listCard))
	hasTriangle := false
	handCard := NewHandCards()

	var list entity.ListCard
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			handCard.MapCardType[pb.HandRanking_ThreeOfAKind] = listCard
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if hasTriangle {
		newListCard = SortCard(newListCard)
		handCard.ListCard = newListCard
		return handCard, true
	}
	return nil, false
}

// CheckTwoPairs
// Thú (en: Two Pairs)
// Hai đôi
func CheckTwoPairs(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 3 {
		return nil, false
	}

	newListCard := make(entity.ListCard, 0, len(listCard))
	numPair := 0

	var list entity.ListCard
	handCard := NewHandCards()
	listTwoPair := entity.ListCard{}
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 2 {
			numPair++
			newListCard = append(newListCard, list...)
			listTwoPair = append(listTwoPair, list...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if numPair == 2 {
		newListCard = SortCard(newListCard)
		handCard.MapCardType[pb.HandRanking_TwoPairs] = SortCard(listTwoPair)
		handCard.ListCard = newListCard
		return handCard, true
	}
	return nil, false
}

// CheckPair
// Đôi (en: Pair)
// Hai lá bài đồng số
func CheckPair(listCard entity.ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 3 {
		return nil, false
	}

	newListCard := make(entity.ListCard, 0, len(listCard))
	numPair := 0

	var list entity.ListCard
	cards := NewHandCards()
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 2 {
			numPair++
			newListCard = append(list, newListCard...)
			cards.MapCardType[pb.HandRanking_Pair] = list
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if numPair > 0 {
		newListCard = SortCard(newListCard)
		cards.ListCard = newListCard
		return cards, true
	}
	return nil, false
}

func ToMapRank(listCard entity.ListCard) *linkedhashmap.Map {
	mapCardRank := linkedhashmap.New()
	for i := range listCard {
		var list entity.ListCard
		card := listCard[i]
		rankPoint := card.GetRank()
		if val, exist := mapCardRank.Get(rankPoint); !exist {
			list = entity.ListCard{}
		} else {
			list = *(val.(*entity.ListCard))
		}
		list = append(list, card)
		mapCardRank.Put(rankPoint, &list)
	}

	return mapCardRank
}

func ToMapSuit(listCard entity.ListCard) *linkedhashmap.Map {
	m := linkedhashmap.New()
	for _, card := range listCard {
		var list entity.ListCard
		suitPoint := card.GetSuit()
		if val, exist := m.Get(suitPoint); !exist {
			list = entity.ListCard{}
		} else {
			list = *(val.(*entity.ListCard))
		}

		list = append(list, card)
		m.Put(suitPoint, &list)
	}

	return m
}

// SortCard
// sort card increase by rank, equal rank will check suit
func SortCard(listCard entity.ListCard) entity.ListCard {
	sort.Slice(listCard, func(a, b int) bool {
		cardA := listCard[a]
		cardB := listCard[b]
		rankPointA := cardA.GetRank()
		rankPointB := cardB.GetRank()
		if rankPointA > rankPointB {
			return true
		}
		if rankPointA < rankPointB {
			return false
		}
		suitPointA := cardA.GetSuit()
		suitPointB := cardB.GetSuit()
		return suitPointA < suitPointB
	})
	return listCard
}

//Mậu binh tới trắng: (Người chơi chiến thắng trực tiếp mà không cần so từng chi)
// Sảnh rồng: 13 lá từ 2 -> A ko đồng chất.
func IsDragon(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	_, isStraight := CheckStraight(listCard)
	// check is straight
	if !isStraight {
		return nil, false
	}
	h := NewHandCards()
	h.ListCard = SortCard(listCard)
	return h, true
}

// Sảnh rồng: 13 lá từ 2 -> A đồng chất.
func IsCleanDragon(listCard entity.ListCard) (*HandCards, bool) {
	_, isDragon := CheckDragon(listCard)
	// check is straight
	if !isDragon {
		return nil, false
	}
	_, isFullColor := IsFullColored(listCard)
	if !isFullColor {
		return nil, false
	}
	h := NewHandCards()
	h.ListCard = SortCard(listCard)
	return h, true
}

// 3 thùng phá sảnh: 3 thùng phá sảnh ở cả ba chi
func IsThreeStraightFlushes(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	listCard = SortCard(listCard)
	front := listCard[:3]
	mid := listCard[3:8]
	back := listCard[8:]
	h := NewHandCards()
	var handCard *HandCards
	isStraightFlush := false
	if handCard, isStraightFlush = CheckStraightFlush(front); !isStraightFlush {
		return nil, false
	}

	h.CopyMapCardType(handCard.MapCardType)
	if _, isStraightFlush = CheckStraight(mid); !isStraightFlush {
		return nil, false
	}
	h.CopyMapCardType(handCard.MapCardType)
	if _, isStraightFlush = CheckStraight(back); !isStraightFlush {
		return nil, false
	}
	h.CopyMapCardType(handCard.MapCardType)
	h.ListCard = SortCard(listCard)
	return h, true
}

// Đồng màu 1: 13 lá đồng màu đen/đỏ.
func IsFullColored(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	mapSuit := ToMapSuit(listCard)
	// check all card same suit
	if len(mapSuit.Keys()) == 1 {
		h := NewHandCards()
		h.ListCard = SortCard(listCard)
		return h, true
	}
	return nil, false
}

// Đồng màu 2: bài có 12 lá đồng màu đen/đỏ hoặc đỏ/đen.
func IsFullColored2(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	mapSuit := ToMapSuit(listCard)
	// check all card same suit
	if len(mapSuit.Keys()) != 2 {
		return nil, false
	}
	for _, v := range mapSuit.Values() {
		list := *(v.(*entity.ListCard))
		if len(list) == 12 {
			h := NewHandCards()
			h.ListCard = SortCard(listCard)
			return h, true
		}
	}
	return nil, false
}

// 5 đôi 1 xám: bài có 5 đôi và 1 xám cô. Giống nhau so sánh đến lá lớn nhất trong xám.
func IsFivePairThreeOfAKind(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	listCard = SortCard(listCard)
	mapRank := ToMapRank(listCard)
	numPair := 0
	hasThreeOfAKind := false
	h := NewHandCards()
	for _, v := range mapRank.Values() {
		list := *(v.(*entity.ListCard))
		size := len(list)
		if size > 3 {
			return nil, false
		}
		if size%2 == 0 {
			numPair += size / 2
		}
		if size == 3 {
			h.MapCardType[pb.HandRanking_ThreeOfAKind] = list
			hasThreeOfAKind = true
		}
	}
	if numPair == 5 && hasThreeOfAKind {
		h.ListCard = listCard
		return h, true
	}
	return nil, false
}

// Lục phé bôn: bài có 6 đôi và 1 mậu thầu. Giống nhau so đến đôi cao nhất.
func IsSixAndAHalfPairs(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	listCard = SortCard(listCard)
	mapRank := ToMapRank(listCard)
	numPair := 0
	h := NewHandCards()
	for _, v := range mapRank.Values() {
		list := *(v.(*entity.ListCard))
		size := len(list)
		if size%2 == 0 {
			numPair += size / 2
			l := h.MapCardType[pb.HandRanking_TwoPairs]
			l = append(l, list...)
			h.MapCardType[pb.HandRanking_TwoPairs] = l
			continue
		}

		return nil, false
	}
	valid := numPair == 6
	if !valid {
		return nil, false
	}

	h.ListCard = listCard
	return h, true
}

// 3 tứ quý: bài có 3 tứ quý. Giống nhau so đến tứ quý cao nhất.
func ThreeQuads(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	listCard = SortCard(listCard)
	mapRank := ToMapRank(listCard)
	h := NewHandCards()
	if len(mapRank.Keys()) != 3 {
		return nil, false
	}
	for _, v := range mapRank.Values() {
		list := *(v.(*entity.ListCard))
		if len(list) != 4 {
			return nil, false
		}
		l := h.MapCardType[pb.HandRanking_FourOfAKind]
		l = append(l, list...)
		h.MapCardType[pb.HandRanking_FourOfAKind] = l
	}
	h.ListCard = listCard
	return h, true
}

// 3 cái thùng: 3 chi mỗi chi là một thùng. Giống nhau so đến các thùng ở các chi. Có thể hoà.
func IsThreeFlushes(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	listCard = SortCard(listCard)
	front := listCard[:3]
	mid := listCard[3:8]
	back := listCard[8:]
	h := NewHandCards()
	var handCard *HandCards
	isFlush := false
	if handCard, isFlush = CheckFlush(front); !isFlush {
		return nil, false
	}

	h.CopyMapCardType(handCard.MapCardType)
	if handCard, isFlush = CheckFlush(mid); !isFlush {
		return nil, false
	}
	h.CopyMapCardType(handCard.MapCardType)
	if handCard, isFlush = CheckFlush(back); !isFlush {
		return nil, false
	}
	h.CopyMapCardType(handCard.MapCardType)
	h.ListCard = listCard
	return h, true
}

// 3 cái sảnh: 3 chi mỗi chi là một sảnh. Giống nhau so đến các sảnh ở các chi. Có thể hoà.
func IsThreeStraight(listCard entity.ListCard) (*HandCards, bool) {
	if len(listCard) != 13 {
		return nil, false
	}
	listCard = SortCard(listCard)
	front := listCard[:3]
	mid := listCard[3:8]
	back := listCard[8:]
	h := NewHandCards()
	var handCard *HandCards
	isStraight := false
	if handCard, isStraight = CheckStraight(front); !isStraight {
		return nil, false
	}

	h.CopyMapCardType(handCard.MapCardType)
	if handCard, isStraight = CheckStraight(mid); !isStraight {
		return nil, false
	}

	h.CopyMapCardType(handCard.MapCardType)
	if handCard, isStraight = CheckStraight(back); !isStraight {
		return nil, false
	}

	h.CopyMapCardType(handCard.MapCardType)
	h.ListCard = listCard
	return h, true
}

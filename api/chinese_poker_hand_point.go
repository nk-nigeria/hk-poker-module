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

func (h *HandPoint) IsStraight() bool {
	return h.rankingType == pb.HandRanking_Straight
}

func (h *HandPoint) IsFlush() bool {
	return h.rankingType == pb.HandRanking_Flush
}

type CheckFunc func(entity.ListCard) (*HandCards, bool)

var HandChecker = map[pb.HandRanking]CheckFunc{
	pb.HandRanking_StraightFlush: CheckStraightFlush,
	pb.HandRanking_FourOfAKind:   CheckFourOfAKind,
	pb.HandRanking_FullHouse:     CheckFullHouse,
	pb.HandRanking_Flush:         CheckFlush,
	pb.HandRanking_Straight:      CheckStraight,
	pb.HandRanking_ThreeOfAKind:  CheckThreeOfAKind,
	pb.HandRanking_TwoPairs:      CheckTwoPairs,
	pb.HandRanking_Pair:          CheckPair,
}

var HandCheckerFront = map[pb.HandRanking]CheckFunc{
	pb.HandRanking_ThreeOfAKind: CheckThreeOfAKind,
	pb.HandRanking_Pair:         CheckPair,
}

func GetHandPoint(listCard entity.ListCard) (*HandPoint, *HandCards) {
	if len(listCard) == 3 {
		// For check front
		for rank, check := range HandCheckerFront {
			if handCard, valid := check(listCard); valid {
				return &HandPoint{
					rankingType: rank,
				}, handCard
			}
		}
	} else {
		for rank, check := range HandChecker {
			if handCard, valid := check(listCard); valid {
				return &HandPoint{
					rankingType: rank,
				}, handCard
			}
		}
	}

	// Sort for high card
	listCard = SortCard(listCard)
	handCard := HandCards{
		ListCard: listCard,
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
	for k, v := range handCard2.MapCardType {
		handCard.MapCardType[k] = v
	}
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
	if mapCardRank.Size() != 2 {
		return nil, false
	}

	newListCard := make(entity.ListCard, 0, len(listCard))
	handCard := NewHandCards()
	isFourOfAKind := false
	var list entity.ListCard
	for _, value := range mapCardRank.Values() {
		list = value.(entity.ListCard)
		if len(list) == 4 {
			isFourOfAKind = true
		}
		handCard.MapCardType[pb.HandRanking_FourOfAKind] = list
		newListCard = append(newListCard, list...)
	}
	if isFourOfAKind {
		newListCard = SortCard(newListCard)
		handCard.ListCard = newListCard
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
		list = value.(entity.ListCard)
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
		ListCard: listCard,
	}
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
		ListCard: listCard,
	}
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
	handCard := HandCards{}

	var list entity.ListCard
	for _, value := range mapCardRank.Values() {
		list = value.(entity.ListCard)
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
		return &handCard, true
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
	handCard := HandCards{}
	listTwoPair := entity.ListCard{}
	for _, value := range mapCardRank.Values() {
		list = value.(entity.ListCard)
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
		return &handCard, true
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
	cards := &HandCards{}
	for _, value := range mapCardRank.Values() {
		list = value.(entity.ListCard)
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
	var list entity.ListCard
	for i := range listCard {
		card := listCard[i]
		rankPoint := card.GetRank()
		if val, exist := mapCardRank.Get(rankPoint); !exist {
			list = entity.ListCard{}
			mapCardRank.Put(rankPoint, list)
		} else {
			list = val.(entity.ListCard)
		}
		list = append(list, card)
	}

	return mapCardRank
}

func ToMapSuit(listCard entity.ListCard) *linkedhashmap.Map {
	m := linkedhashmap.New()
	var list entity.ListCard
	for _, card := range listCard {
		suitPoint := card.GetSuit()
		if val, exist := m.Get(suitPoint); !exist {
			list = entity.ListCard{}
			m.Put(suitPoint, list)
		} else {
			list = val.(entity.ListCard)
		}

		list = append(list, card)
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

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

type ListCard []*pb.Card

func (ls ListCard) GetMaxPointCard() int {
	isStraight := true
	isContainCardRanK := false
	prevRankPoint := entity.GetCardRankPoint(ls[0].GetRank())
	for i := 1; i < len(ls); i++ {
		card := ls[i]
		if isStraight {
			rankPoint := entity.GetCardRankPoint(card.GetRank())
			prevRankPoint--
			if rankPoint != prevRankPoint {
				isStraight = false
			}
		}
		isContainCardRanK = card.GetRank() == pb.CardRank_RANK_K
	}
	maxCard := ls[len(ls)-1]
	// Chú ý rằng trong Mậu Binh có thể xếp sảnh (hoặc thùng phá sảnh) con A ghép với 2,3,4,5
	// (tuy nhiên đây là bài sảnh hay thùng phá sảnh nhỏ nhất),
	// còn con A ghép với 10,J,Q,K là lá bài lớn nhất.
	if isStraight && !isContainCardRanK {
		return 1
	}
	return entity.GetCardRankPoint(maxCard.GetRank())
}

// -1 lower
// 0 equal
// 1 higher
func (ls ListCard) CompareStraightCards(other ListCard) int {
	if len(ls) != len(other) {
		return 0
	}
	for i := len(ls); i >= 0; i-- {
		point1 := entity.GetCardRankPoint(ls[i].GetRank())
		point2 := entity.GetCardRankPoint(other[i].GetRank())
		if point1 > point2 {
			return 1
		}
		if point2 < point1 {
			return -1
		}
	}
	return 0
}

type HandCards struct {
	ListCard    ListCard
	MapCardType map[pb.HandRanking]ListCard
}

func (h *HandPoint) IsStraight() bool {
	return h.rankingType == pb.HandRanking_Straight
}

func (h *HandPoint) IsFlush() bool {
	return h.rankingType == pb.HandRanking_Flush
}

type CheckFunc func(ListCard) (*HandCards, bool)

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

func GetHandPoint(listCard ListCard) (*HandPoint, *HandCards) {
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
func CheckStraightFlush(listCard ListCard) (*HandCards, bool) {
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
func CheckFourOfAKind(listCard ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() != 2 {
		return nil, false
	}

	newListCard := make(ListCard, 0, len(listCard))
	handCard := HandCards{}
	isFourOfAKind := false
	var list *pb.ListCard
	for _, value := range mapCardRank.Values() {
		list = value.(*pb.ListCard)
		if len(list.Cards) == 4 {
			isFourOfAKind = true
		}
		handCard.MapCardType[pb.HandRanking_FourOfAKind] = list.GetCards()
		newListCard = append(newListCard, list.Cards...)
	}
	if isFourOfAKind {
		newListCard = SortCard(newListCard)
		handCard.ListCard = newListCard
		return &handCard, true
	}
	return nil, false
}

// CheckFullHouse
// Cù lũ (en: Full House)
// Một bộ ba và một bộ đôi
// Bốn lá đồng số
func CheckFullHouse(listCard ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() != 2 {
		return nil, false
	}

	newListCard := make(ListCard, 0, len(listCard))
	hasTriangle := false
	hasDouble := false

	var list *pb.ListCard
	handCard := HandCards{}
	for _, value := range mapCardRank.Values() {
		list = value.(*pb.ListCard)
		if len(list.Cards) == 3 {
			hasTriangle = true
			newListCard = append(list.Cards, newListCard...)
			handCard.MapCardType[pb.HandRanking_ThreeOfAKind] = list.GetCards()
			continue
		}
		if len(list.Cards) == 2 {
			hasDouble = true
			newListCard = append(list.Cards, newListCard...)
			handCard.MapCardType[pb.HandRanking_Pair] = list.GetCards()

			continue
		}
		newListCard = append(newListCard, list.Cards...)
	}
	if hasTriangle && hasDouble {
		newListCard = SortCard(newListCard)
		handCard.ListCard = newListCard
		return &handCard, true
	}
	return nil, false
}

// CheckFlush
// Thùng (en: Flush)
// Năm lá bài cùng màu, đồng chất (nhưng không cùng một chuỗi số)
func CheckFlush(listCard ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard = SortCard(listCard)

	prevSuitPoint := entity.GetCardSuitPoint(listCard[0].GetSuit())
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		suitPoint := entity.GetCardSuitPoint(card.GetSuit())
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
func CheckStraight(listCard ListCard) (*HandCards, bool) {
	listCard = SortCard(listCard)
	prevRankPoint := entity.GetCardRankPoint(listCard[0].GetRank())
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.GetRank())
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
func CheckThreeOfAKind(listCard ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 2 {
		return nil, false
	}

	newListCard := make(ListCard, 0, len(listCard))
	hasTriangle := false
	handCard := HandCards{}

	var list *pb.ListCard
	for _, value := range mapCardRank.Values() {
		list = value.(*pb.ListCard)
		if len(list.Cards) == 3 {
			hasTriangle = true
			newListCard = append(list.Cards, newListCard...)
			handCard.MapCardType[pb.HandRanking_ThreeOfAKind] = listCard
			continue
		}
		newListCard = append(newListCard, list.Cards...)
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
func CheckTwoPairs(listCard ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 3 {
		return nil, false
	}

	newListCard := make(ListCard, 0, len(listCard))
	numPair := 0

	var list *pb.ListCard
	handCard := HandCards{}
	for _, value := range mapCardRank.Values() {
		list = value.(*pb.ListCard)
		if len(list.Cards) == 2 {
			numPair++
			newListCard = append(list.Cards, newListCard...)
			continue
		}
		newListCard = append(newListCard, list.Cards...)
	}
	if numPair == 2 {
		newListCard = SortCard(newListCard)
		handCard.MapCardType[pb.HandRanking_TwoPairs] = newListCard
		handCard.ListCard = newListCard
		return &handCard, true
	}
	return nil, false
}

// CheckPair
// Đôi (en: Pair)
// Hai lá bài đồng số
func CheckPair(listCard ListCard) (*HandCards, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 3 {
		return nil, false
	}

	newListCard := make(ListCard, 0, len(listCard))
	numPair := 0

	var list *pb.ListCard
	cards := &HandCards{}
	for _, value := range mapCardRank.Values() {
		list = value.(*pb.ListCard)
		if len(list.Cards) == 2 {
			numPair++
			newListCard = append(list.Cards, newListCard...)
			cards.MapCardType[pb.HandRanking_Pair] = list.GetCards()
			continue
		}
		newListCard = append(newListCard, list.Cards...)
	}
	if numPair > 0 {
		newListCard = SortCard(newListCard)
		cards.ListCard = newListCard
		return cards, true
	}
	return nil, false
}

func ToMapRank(listCard ListCard) *linkedhashmap.Map {
	mapCardRank := linkedhashmap.New()
	var list *pb.ListCard
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		if val, exist := mapCardRank.Get(rankPoint); !exist {
			list = &pb.ListCard{}
			mapCardRank.Put(rankPoint, list)
		} else {
			list = val.(*pb.ListCard)
		}
		list.Cards = append(list.Cards, card)
	}

	return mapCardRank
}

func ToMapSuit(listCard ListCard) *linkedhashmap.Map {
	m := linkedhashmap.New()
	var list *pb.ListCard
	for _, card := range listCard {
		suitPoint := entity.GetCardSuitPoint(card.Suit)
		if val, exist := m.Get(suitPoint); !exist {
			list = &pb.ListCard{}
			m.Put(suitPoint, list)
		} else {
			list = val.(*pb.ListCard)
		}

		list.Cards = append(list.Cards, card)
	}

	return m
}

// SortCard
// sort card increase by rank, equal rank will check suit
func SortCard(listCard ListCard) ListCard {
	sort.Slice(listCard, func(a, b int) bool {
		cardA := listCard[a]
		cardB := listCard[b]
		rankPointA := entity.GetCardRankPoint(cardA.GetRank())
		rankPointB := entity.GetCardRankPoint(cardB.GetRank())
		if rankPointA > rankPointB {
			return true
		}
		if rankPointA < rankPointB {
			return false
		}
		suitPointA := entity.GetCardSuitPoint(cardA.GetSuit())
		suitPointB := entity.GetCardSuitPoint(cardB.GetSuit())
		return suitPointA < suitPointB
	})
	return listCard
}

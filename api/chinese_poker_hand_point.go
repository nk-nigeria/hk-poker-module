package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"sort"
)

type HandPoint struct {
	rankingType pb.HandRanking
	point       int
}

type CheckFunc func([]*pb.Card) ([]*pb.Card, bool)

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

func GetHandPoint(listCard []*pb.Card) (*HandPoint, []*pb.Card) {
	if len(listCard) == 3 {
		// For check front
		for rank, check := range HandCheckerFront {
			if sortedListCard, valid := check(listCard); valid {
				return &HandPoint{
					rankingType: rank,
				}, sortedListCard
			}
		}
	} else {
		for rank, check := range HandChecker {
			if sortedListCard, valid := check(listCard); valid {
				return &HandPoint{
					rankingType: rank,
				}, sortedListCard
			}
		}
	}

	// TODO: sort high card
	return &HandPoint{
		rankingType: pb.HandRanking_HighCard,
	}, listCard
}

// CheckStraightFlush
// Thùng phá sảnh (en: Straight Flush)
// Năm lá bài cùng màu, đồng chất, cùng một chuỗi số
// Là Flush, có cùng chuỗi
func CheckStraightFlush(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard, valid := CheckFlush(listCard)
	if !valid {
		return nil, false
	}
	listCard, valid = CheckStraight(listCard)
	if !valid {
		return nil, false
	}
	return listCard, true
}

// CheckFourOfAKind
// Tứ quý (en: Four of a Kind)
// Bốn lá đồng số
func CheckFourOfAKind(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; exist {
			list = append(list, card)
		} else {
			list = make([]*pb.Card, 0)
			list = append(list, card)
		}
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) != 2 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	isFourOfAKind := false
	for _, list := range mapCardRank {
		if len(list) == 4 {
			isFourOfAKind = true
		}
		newListCard = append(newListCard, list...)
	}
	if isFourOfAKind {
		newListCard = SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// CheckFullHouse
// Cù lũ (en: Full House)
// Một bộ ba và một bộ đôi
// Bốn lá đồng số
func CheckFullHouse(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; exist {
			list = append(list, card)
		} else {
			list = make([]*pb.Card, 0)
			list = append(list, card)
		}
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) != 2 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	hasTriangle := false
	hasDouble := false

	for _, list := range mapCardRank {
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			continue
		}
		if len(list) == 2 {
			hasDouble = true
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if hasTriangle && hasDouble {
		newListCard = SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// CheckFlush
// Thùng (en: Flush)
// Năm lá bài cùng màu, đồng chất (nhưng không cùng một chuỗi số)
func CheckFlush(listCard []*pb.Card) ([]*pb.Card, bool) {
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
	return listCard, true
}

// CheckStraight
// Sảnh (en: Straight)
// Năm lá bài trong một chuỗi số (nhưng không đồng chất)
func CheckStraight(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard = SortCard(listCard)
	prevRankPoint := entity.GetCardRankPoint(listCard[0].GetRank())
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.GetRank())
		prevRankPoint++
		if rankPoint != prevRankPoint {
			return nil, false
		}
	}
	return listCard, true
}

// CheckThreeOfAKind
// Xám chi/Xám cô (en: Three of a Kind)
// Ba lá bài đồng số
func CheckThreeOfAKind(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; !exist {
			list = make([]*pb.Card, 0)
		}
		list = append(list, card)
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) < 2 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	hasTriangle := false

	for _, list := range mapCardRank {
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if hasTriangle {
		newListCard = SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// CheckTwoPairs
// Thú (en: Two Pairs)
// Hai đôi
func CheckTwoPairs(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; !exist {
			list = make([]*pb.Card, 0)
		}
		list = append(list, card)
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) < 3 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	numPair := 0

	for _, list := range mapCardRank {
		if len(list) == 2 {
			numPair++
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if numPair == 2 {
		newListCard = SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// CheckPair
// Đôi (en: Pair)
// Hai lá bài đồng số
func CheckPair(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; !exist {
			list = make([]*pb.Card, 0)
		}
		list = append(list, card)
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) < 3 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	numPair := 0

	for _, list := range mapCardRank {
		if len(list) == 2 {
			numPair++
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if numPair > 0 {
		newListCard = SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// SortCard
// sort card increase by rank, equal rank will check suit
func SortCard(listCard []*pb.Card) []*pb.Card {
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

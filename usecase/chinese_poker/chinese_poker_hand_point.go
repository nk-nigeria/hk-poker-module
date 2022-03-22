package chinese_poker

import (
	"fmt"
	"sort"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

//  				t1		s1		s2		s3		s4		s5
//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF

const (
	//	3				t1		m1		m2		m3
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	5				t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointHighCard = uint8(0x01)
	//	3				t1		d1		m1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	5				t1		d1		m1		m2		m3
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointPair = uint8(0x02)
	//					t1		d1		d2		m1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointTwoPairs = uint8(0x03)
	//	3				t1		s1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	//	5				t1		s1		m1		m2
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointThreeOfAKind = uint8(0x04)
	//					t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointStraight = uint8(0x05)
	//					t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointFlush = uint8(0x06)
	//					t1		s1		d1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointFullHouse = uint8(0x07)
	//					t1		q1		m1
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointFourOfAKind = uint8(0x08)
	//					t1		m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointStraightFlush = uint8(0x09)

	//			n1				m1		m2		m3		m4		m5
	//	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF 	0xFF
	ScorePointNaturalThreeStraights = uint8(0x01)
	ScorePointNaturalThreeOfFlushes = uint8(0x02)
	ScorePointNaturalSixPairs       = uint8(0x03)
	ScorePointNaturalFullColors     = uint8(0x04)
	ScorePointNaturalDragon         = uint8(0x05)
	ScorePointNaturalCleanDragon    = uint8(0x06)
)

func createPoint(t, p1, p2, p3, p4, p5 uint8) uint64 {
	var point uint64 = 0
	point |= uint64(t) << (5 * 8)
	point |= uint64(p1) << (4 * 8)
	point |= uint64(p2) << (3 * 8)
	point |= uint64(p3) << (2 * 8)
	point |= uint64(p4) << (1 * 8)
	point |= uint64(p5)

	return point
}

func createPointNaturalCard(t uint8, cards entity.ListCard) (uint64, uint64) {
	var hpoint uint64 = 0
	var lpoint uint64 = 0

	hpoint |= uint64(t) << (6 * 8)
	hpoint |= uint64(cards[0].GetRank()) << (4 * 8)
	hpoint |= uint64(cards[1].GetRank()) << (3 * 8)
	hpoint |= uint64(cards[2].GetRank()) << (2 * 8)
	hpoint |= uint64(cards[3].GetRank()) << (1 * 8)
	hpoint |= uint64(cards[4].GetRank())

	lpoint |= uint64(cards[5].GetRank()) << (7 * 8)
	lpoint |= uint64(cards[6].GetRank()) << (6 * 8)
	lpoint |= uint64(cards[7].GetRank()) << (5 * 8)
	lpoint |= uint64(cards[8].GetRank()) << (4 * 8)
	lpoint |= uint64(cards[9].GetRank()) << (3 * 8)
	lpoint |= uint64(cards[10].GetRank()) << (2 * 8)
	lpoint |= uint64(cards[11].GetRank()) << (1 * 8)
	lpoint |= uint64(cards[12].GetRank())

	return hpoint, lpoint
}

type HandPoint struct {
	rankingType pb.HandRanking
	point       uint64
	lpoint      uint64
}

func (h HandPoint) String() string {
	return fmt.Sprintf("Rank %v, Point: %v", h.rankingType, h.point)
}

func (h *HandPoint) IsStraight() bool {
	return h.rankingType == pb.HandRanking_Straight
}

func (h *HandPoint) IsFlush() bool {
	return h.rankingType == pb.HandRanking_Flush
}

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

func CalculatePoint(listCard entity.ListCard) *HandPoint {
	if len(listCard) == 3 {
		// For check front
		for _, key := range HandCheckerFront.Keys() {
			rank := key.(pb.HandRanking)
			val, exists := HandCheckerFront.Get(rank)
			if !exists {
				return nil
			}
			fn := val.(func(listCard entity.ListCard) (*HandPoint, bool))
			if handPoint, valid := fn(listCard); valid {
				return handPoint
			}
		}
	} else {
		for _, key := range HandChecker.Keys() {
			rank := key.(pb.HandRanking)
			val, exist := HandChecker.Get(rank)
			if !exist {
				return nil
			}
			fn := val.(func(listCard entity.ListCard) (*HandPoint, bool))
			if handPoint, valid := fn(listCard); valid {
				return handPoint
			}
		}
	}

	listCard = SortCard(listCard)

	if len(listCard) == 3 {
		return &HandPoint{
			rankingType: pb.HandRanking_HighCard,
			point: createPoint(ScorePointHighCard,
				listCard[0].GetRank(),
				listCard[1].GetRank(),
				listCard[2].GetRank(),
				0,
				0,
			),
		}
	} else {
		return &HandPoint{
			rankingType: pb.HandRanking_HighCard,
			point: createPoint(ScorePointHighCard,
				listCard[0].GetRank(),
				listCard[1].GetRank(),
				listCard[2].GetRank(),
				listCard[3].GetRank(),
				listCard[4].GetRank(),
			),
		}
	}
}

// CheckStraightFlush
// Thùng phá sảnh (en: Straight Flush)
// Năm lá bài cùng màu, đồng chất, cùng một chuỗi số
// Là Flush, có cùng chuỗi
func CheckStraightFlush(listCard entity.ListCard) (*HandPoint, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	_, valid := CheckFlush(listCard)
	if !valid {
		return nil, false
	}
	_, valid = CheckStraight(listCard)
	if !valid {
		return nil, false
	}

	handPoint := &HandPoint{
		rankingType: pb.HandRanking_StraightFlush,
		point: createPoint(ScorePointStraightFlush,
			listCard[0].GetRank(),
			listCard[1].GetRank(),
			listCard[2].GetRank(),
			listCard[3].GetRank(),
			listCard[4].GetRank()),
	}

	return handPoint, true
}

// CheckFourOfAKind
// Tứ quý (en: Four of a Kind)
// Bốn lá đồng số
func CheckFourOfAKind(listCard entity.ListCard) (*HandPoint, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)

	isFourOfAKind := false
	var list entity.ListCard
	var remain entity.ListCard
	var fourOfKindCard entity.Card
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 4 {
			isFourOfAKind = true
			fourOfKindCard = list[0]
		} else {
			remain = append(remain, list...)
		}
	}
	if isFourOfAKind {
		remain = SortCard(remain)

		handPoint := &HandPoint{
			rankingType: pb.HandRanking_FourOfAKind,
			point: createPoint(ScorePointFourOfAKind,
				fourOfKindCard.GetRank(),
				remain[0].GetRank(),
				0,
				0,
				0),
		}
		return handPoint, true
	}
	return nil, false
}

// CheckFullHouse
// Cù lũ (en: Full House)
// Một bộ ba và một bộ đôi
// Bốn lá đồng số
func CheckFullHouse(listCard entity.ListCard) (*HandPoint, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() != 2 {
		return nil, false
	}

	hasTriangle := false
	hasDouble := false

	var list entity.ListCard
	var p1, p2 entity.Card
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 3 {
			hasTriangle = true
			p1 = list[0]
			continue
		}
		if len(list) == 2 {
			hasDouble = true
			p2 = list[0]
			continue
		}
	}
	if hasTriangle && hasDouble {
		handPoint := &HandPoint{
			rankingType: pb.HandRanking_FullHouse,
			point: createPoint(ScorePointFullHouse,
				p1.GetRank(),
				p2.GetRank(),
				0,
				0,
				0),
		}

		return handPoint, true
	}
	return nil, false
}

// CheckFlush
// Thùng (en: Flush)
// Năm lá bài cùng màu, đồng chất (nhưng không cùng một chuỗi số)
func CheckFlush(listCard entity.ListCard) (*HandPoint, bool) {
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

	handPoint := &HandPoint{
		rankingType: pb.HandRanking_Flush,
		point: createPoint(ScorePointFlush,
			listCard[0].GetRank(),
			listCard[1].GetRank(),
			listCard[2].GetRank(),
			listCard[3].GetRank(),
			listCard[4].GetRank()),
	}

	return handPoint, true
}

// CheckStraight
// Sảnh (en: Straight)
// Năm lá bài trong một chuỗi số (nhưng không đồng chất)
func CheckStraight(listCard entity.ListCard) (*HandPoint, bool) {
	listCard = SortCard(listCard)

	cardAIsOnePoint := false
	if listCard[0].GetRank() == entity.RankA {
		if listCard[1].GetRank() != entity.RankK {
			cardAIsOnePoint = true
			newList := listCard.Clone()
			newList = append(newList[1:], newList[0])
			listCard = newList
		}
	}
	prevRankPoint := listCard[0].GetRank()

	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		rankPoint := card.GetRank()
		if rankPoint == entity.RankA && cardAIsOnePoint {
			rankPoint = 1
		}
		prevRankPoint--
		if rankPoint != prevRankPoint {
			return nil, false
		}
	}

	handPoint := &HandPoint{
		rankingType: pb.HandRanking_Straight,
		point: createPoint(ScorePointStraight,
			listCard[0].GetRank(),
			listCard[1].GetRank(),
			listCard[2].GetRank(),
			listCard[3].GetRank(),
			listCard[4].GetRank()),
	}

	return handPoint, true
}

// CheckThreeOfAKind
// Xám chi/Xám cô (en: Three of a Kind)
// Ba lá bài đồng số
func CheckThreeOfAKind(listCard entity.ListCard) (*HandPoint, bool) {
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

	var list entity.ListCard
	var remain entity.ListCard
	var threeOfCard entity.Card
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			threeOfCard = list[0]
			continue
		} else {
			remain = append(remain, list...)
		}
	}
	if hasTriangle {
		remain = SortCard(remain)

		var handPoint *HandPoint
		if l == 3 {
			handPoint = &HandPoint{
				rankingType: pb.HandRanking_ThreeOfAKind,
				point: createPoint(ScorePointThreeOfAKind,
					threeOfCard.GetRank(),
					0,
					0,
					0,
					0),
			}
		} else {
			handPoint = &HandPoint{
				rankingType: pb.HandRanking_ThreeOfAKind,
				point: createPoint(ScorePointThreeOfAKind,
					threeOfCard.GetRank(),
					remain[0].GetRank(),
					remain[1].GetRank(),
					0,
					0),
			}
		}

		return handPoint, true
	}
	return nil, false
}

// CheckTwoPairs
// Thú (en: Two Pairs)
// Hai đôi
func CheckTwoPairs(listCard entity.ListCard) (*HandPoint, bool) {
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
	var remain entity.ListCard
	var pairs entity.ListCard
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 2 {
			numPair++
			newListCard = append(newListCard, list...)
			pairs = append(pairs, list[0])
			continue
		} else {
			remain = append(remain, list...)
		}
	}
	if numPair == 2 {
		pairs = SortCard(pairs)

		handPoint := &HandPoint{
			rankingType: pb.HandRanking_TwoPairs,
			point: createPoint(ScorePointTwoPairs,
				pairs[0].GetRank(),
				pairs[1].GetRank(),
				remain[0].GetRank(),
				0,
				0),
		}

		return handPoint, true
	}
	return nil, false
}

// CheckPair
// Đôi (en: Pair)
// Hai lá bài đồng số
func CheckPair(listCard entity.ListCard) (*HandPoint, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := ToMapRank(listCard)
	if mapCardRank.Size() < 3 {
		return nil, false
	}

	numPair := 0

	var list entity.ListCard
	var pair entity.Card
	var remain entity.ListCard
	for _, value := range mapCardRank.Values() {
		list = *(value.(*entity.ListCard))
		if len(list) == 2 {
			numPair++
			pair = list[0]
			continue
		} else {
			remain = append(remain, list...)
		}
	}
	if numPair > 0 {
		remain = SortCard(remain)
		var handPoint *HandPoint
		if l == 3 {
			handPoint = &HandPoint{
				rankingType: pb.HandRanking_Pair,
				point: createPoint(ScorePointPair,
					pair.GetRank(),
					remain[0].GetRank(),
					0,
					0,
					0),
			}
		} else {
			handPoint = &HandPoint{
				rankingType: pb.HandRanking_ThreeOfAKind,
				point: createPoint(ScorePointThreeOfAKind,
					pair.GetRank(),
					remain[0].GetRank(),
					remain[1].GetRank(),
					remain[2].GetRank(),
					0),
			}
		}

		return handPoint, true
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

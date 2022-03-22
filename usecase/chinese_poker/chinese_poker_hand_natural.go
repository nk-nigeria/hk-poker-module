package chinese_poker

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"log"
)

var naturalCardChecker map[uint8]func(entity.ListCard) (*HandPoint, bool)
var naturalHandChecker map[uint8]func(*Hand) (*HandPoint, bool)

func init() {
	naturalCardChecker = make(map[uint8]func(entity.ListCard) (*HandPoint, bool))
	naturalCardChecker[ScorePointNaturalCleanDragon] = CheckCleanDragon
	naturalCardChecker[ScorePointNaturalDragon] = CheckDragon
	naturalCardChecker[ScorePointNaturalSixPairs] = CheckSixPairs
	naturalCardChecker[ScorePointNaturalFullColors] = CheckFullColor

	naturalHandChecker = make(map[uint8]func(*Hand) (*HandPoint, bool))
	naturalHandChecker[ScorePointNaturalThreeOfFlushes] = CheckThreeFlushes
	naturalHandChecker[ScorePointNaturalThreeStraights] = CheckThreeStraights
}

func CheckNaturalCards(h *Hand) (*HandPoint, bool) {
	// check natural win
	for _, checkerFn := range naturalCardChecker {
		handPoint, valid := checkerFn(h.GetCards())
		if valid {
			return handPoint, valid
		}
	}

	return nil, false
}

func CheckNaturalHands(hand *Hand) (*HandPoint, bool) {
	// check natural win
	for _, checkerFn := range naturalHandChecker {
		handPoint, valid := checkerFn(hand)
		if valid {
			return handPoint, valid
		}
	}

	return nil, false
}

// CheckCleanDragon
// Sảnh rồng đồng màu
func CheckCleanDragon(listCard entity.ListCard) (*HandPoint, bool) {
	mapCardSuit := ToMapSuit(listCard)
	log.Printf("key %v", len(mapCardSuit.Keys()))
	for _, k := range mapCardSuit.Keys() {
		log.Printf("key %v", k)
	}
	if len(mapCardSuit.Keys()) > 1 {
		return nil, false
	}

	listCard = SortCard(listCard)
	hpoint, lpoint := createPointNaturalCard(ScorePointNaturalCleanDragon, listCard)

	return &HandPoint{
		rankingType: pb.HandRanking_NaturalCleanDragon,
		point:       hpoint,
		lpoint:      lpoint,
	}, true
}

// CheckFullColor
// Đồng màu 12 lá
func CheckFullColor(listCard entity.ListCard) (*HandPoint, bool) {
	mapCardSuit := ToMapSuit(listCard)
	redCount := 0
	blackCount := 0
	var list entity.ListCard
	if val, exist := mapCardSuit.Get(entity.SuitHearts); exist {
		list = *(val.(*entity.ListCard))
		redCount += len(list)
	}
	if val, exist := mapCardSuit.Get(entity.SuitDiamonds); exist {
		list = *(val.(*entity.ListCard))
		redCount += len(list)
	}

	if val, exist := mapCardSuit.Get(entity.SuitSpides); exist {
		list = *(val.(*entity.ListCard))
		blackCount += len(list)
	}
	if val, exist := mapCardSuit.Get(entity.SuitClubs); exist {
		list = *(val.(*entity.ListCard))
		blackCount += len(list)
	}

	if redCount >= 12 || blackCount >= 12 {
		listCard = SortCard(listCard)
		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalFullColors, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalFullColors,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}

	return nil, false
}

// CheckDragon
// Sảnh rồng
func CheckDragon(listCard entity.ListCard) (*HandPoint, bool) {
	_, ok := CheckStraight(listCard)
	if ok {
		listCard = SortCard(listCard)
		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalDragon, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalDragon,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}
	return nil, false
}

// CheckSixPairs
// 6 đôi
func CheckSixPairs(listCard entity.ListCard) (*HandPoint, bool) {
	mapRank := ToMapRank(listCard)
	if len(mapRank.Keys()) < 6 {
		return nil, false
	}

	var list entity.ListCard
	var numPairs = 0
	for _, val := range mapRank.Values() {
		list = *(val.(*entity.ListCard))
		if len(list)%2 == 0 {
			numPairs++
		}
	}

	if numPairs == 6 {
		listCard = SortCard(listCard)
		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalSixPairs, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalSixPairs,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}

	return nil, false
}

// CheckThreeStraight
// 3 sảnh
func CheckThreeStraights(hand *Hand) (*HandPoint, bool) {
	threeStraight := hand.frontHand.Point.IsStraight() && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	if threeStraight {
		var listCard entity.ListCard
		listCard = append(listCard, SortCard(hand.backHand.Cards)...)
		listCard = append(listCard, SortCard(hand.middleHand.Cards)...)
		listCard = append(listCard, SortCard(hand.frontHand.Cards)...)

		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalThreeStraights, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalThreeStraights,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}

	return nil, false
}

// CheckThreeFlushes
// 3 cái thùng
func CheckThreeFlushes(hand *Hand) (*HandPoint, bool) {
	threeFlush := hand.frontHand.Point.IsStraight() && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	if threeFlush {
		var listCard entity.ListCard
		listCard = append(listCard, SortCard(hand.backHand.Cards)...)
		listCard = append(listCard, SortCard(hand.middleHand.Cards)...)
		listCard = append(listCard, SortCard(hand.frontHand.Cards)...)

		hpoint, lpoint := createPointNaturalCard(ScorePointNaturalThreeOfFlushes, listCard)

		return &HandPoint{
			rankingType: pb.HandRanking_NaturalThreeOfFlushes,
			point:       hpoint,
			lpoint:      lpoint,
		}, true
	}

	return nil, false
}

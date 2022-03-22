package chinese_poker

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"log"
)

var naturalCardChecker map[pb.NaturalRanking]func(entity.ListCard) bool
var naturalHandChecker map[pb.NaturalRanking]func(*Hand) bool

func init() {
	naturalCardChecker = make(map[pb.NaturalRanking]func(entity.ListCard) bool)
	naturalCardChecker[pb.NaturalRanking_CleanDragon] = CheckCleanDragon
	naturalCardChecker[pb.NaturalRanking_Dragon] = CheckDragon
	naturalCardChecker[pb.NaturalRanking_SixPairs] = CheckSixPairs
	naturalCardChecker[pb.NaturalRanking_FullColors] = CheckFullColor

	naturalHandChecker = make(map[pb.NaturalRanking]func(*Hand) bool)
	naturalHandChecker[pb.NaturalRanking_ThreeOfFlushes] = CheckThreeFlushes
	naturalHandChecker[pb.NaturalRanking_ThreeStraights] = CheckThreeStraights
}

func CheckNaturalCards(h *Hand) (bool, pb.NaturalRanking) {
	// check natural win
	for k, checkerFn := range naturalCardChecker {
		valid := checkerFn(h.GetCards())
		if valid {
			return valid, k
		}
	}

	return false, pb.NaturalRanking_None
}

func CheckNaturalHands(hand *Hand) (bool, pb.NaturalRanking) {
	// check natural win
	for k, checkerFn := range naturalHandChecker {
		valid := checkerFn(hand)
		if valid {
			return valid, k
		}
	}

	return false, pb.NaturalRanking_None
}

// CheckCleanDragon
// Sảnh rồng đồng màu
func CheckCleanDragon(listCard entity.ListCard) bool {
	mapCardSuit := ToMapSuit(listCard)
	log.Printf("key %v", len(mapCardSuit.Keys()))
	for _, k := range mapCardSuit.Keys() {
		log.Printf("key %v", k)
	}
	if len(mapCardSuit.Keys()) > 1 {
		return false
	}

	return true
}

// CheckFullColor
// Đồng màu 12 lá
func CheckFullColor(listCard entity.ListCard) bool {
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
		return true
	}

	return false
}

// CheckDragon
// Sảnh rồng
func CheckDragon(listCard entity.ListCard) bool {
	_, ok := CheckStraight(listCard)
	if ok {
		return true
	}
	return false
}

// CheckSixPairs
// 6 đôi
func CheckSixPairs(listCard entity.ListCard) bool {
	mapRank := ToMapRank(listCard)
	if len(mapRank.Keys()) < 6 {
		return false
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
		return true
	}

	return false
}

// CheckThreeStraight
// 3 sảnh
func CheckThreeStraights(hand *Hand) bool {
	threeStraight := hand.frontHand.Point.IsStraight() && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	if threeStraight {
		return true
	}

	return false
}

// CheckThreeFlushes
// 3 cái thùng
func CheckThreeFlushes(hand *Hand) bool {
	threeFlush := hand.frontHand.Point.IsStraight() && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	if threeFlush {
		return true
	}

	return false
}

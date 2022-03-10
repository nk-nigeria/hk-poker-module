package api

import (
	"log"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

// CheckSameSuit
// Đồng nước 13 lá
func CheckSameSuit(listCard ListCard) (ListCard, bool) {
	mapCardSuit := ToMapSuit(listCard)
	log.Printf("key %v", len(mapCardSuit.Keys()))
	for _, k := range mapCardSuit.Keys() {
		log.Printf("key %v", k)
	}
	if len(mapCardSuit.Keys()) > 1 {
		return nil, false
	}

	return nil, true
}

// CheckSameColor
// Đồng nước 13 lá
func CheckSameColor(listCard ListCard) (ListCard, bool) {
	mapCardSuit := ToMapSuit(listCard)
	redCount := 0
	blackCount := 0
	var list *pb.ListCard
	if val, exist := mapCardSuit.Get(entity.GetCardSuitPoint(pb.CardSuit_SUIT_HEARTS)); exist {
		list = val.(*pb.ListCard)
		redCount += len(list.Cards)
	}
	if val, exist := mapCardSuit.Get(entity.GetCardSuitPoint(pb.CardSuit_SUIT_DIAMONDS)); exist {
		list = val.(*pb.ListCard)
		redCount += len(list.Cards)
	}

	if val, exist := mapCardSuit.Get(entity.GetCardSuitPoint(pb.CardSuit_SUIT_SPADES)); exist {
		list = val.(*pb.ListCard)
		blackCount += len(list.Cards)
	}
	if val, exist := mapCardSuit.Get(entity.GetCardSuitPoint(pb.CardSuit_SUIT_CLUBS)); exist {
		list = val.(*pb.ListCard)
		blackCount += len(list.Cards)
	}

	if redCount >= 12 || blackCount >= 12 {
		return nil, true
	}

	return nil, false
}

// CheckDragon
// Sảnh rồng
func CheckDragon(listCard ListCard) (ListCard, bool) {
	_, ok := CheckStraight(listCard)
	return nil, ok
}

// CheckSixPairs
// 6 đôi
func CheckSixPairs(listCard ListCard) (ListCard, bool) {
	mapRank := ToMapRank(listCard)
	if len(mapRank.Keys()) < 7 {
		return nil, false
	}

	var list *pb.ListCard
	var numPairs = 0
	for _, val := range mapRank.Values() {
		list = val.(*pb.ListCard)
		if len(list.Cards) == 2 {
			numPairs++
		}
	}

	if numPairs == 6 {
		return nil, true
	}

	return nil, false
}

// CheckThreeStraight
// 3 sảnh
func CheckThreeStraight(hand *Hand) (ListCard, bool) {
	threeStraight := hand.frontHand.Point.IsStraight() && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	return nil, threeStraight
}

// CheckThreeFlush
// 3 cái thùng
func CheckThreeFlush(hand *Hand) (ListCard, bool) {
	threeFlush := hand.frontHand.Point.IsStraight() && hand.middleHand.Point.IsStraight() && hand.backHand.Point.IsStraight()
	return nil, threeFlush
}

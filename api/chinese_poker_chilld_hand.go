package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

type ChildHand struct {
	Cards *HandCards
	Point *HandPoint
}

func (ch *ChildHand) calculatePoint() {
	if ch.Point != nil {
		return
	}
	ch.Point, ch.Cards = GetHandPoint(ch.Cards.ListCard)
}

func NewChildHand(cards entity.ListCard) *ChildHand {
	child := ChildHand{
		Cards: &HandCards{
			ListCard: cards,
		},
	}
	return &child
}

func compareChildHand(h1, h2 *ChildHand) int {
	h1.calculatePoint()
	h2.calculatePoint()

	resultPoint := 0

	rank1 := entity.GetHandRankingPoint(h1.Point.rankingType)
	rank2 := entity.GetHandRankingPoint(h2.Point.rankingType)
	if rank1 > rank2 {
		resultPoint++
		return resultPoint
	}
	if rank1 < rank2 {
		resultPoint--
		return resultPoint
	}

	// compare same rank
	point1 := uint8(0)
	point2 := uint8(0)
	extraScore := int8(0)
	switch h1.Point.rankingType {
	case pb.HandRanking_StraightFlush:
		x1 := h1.Cards.MapCardType[pb.HandRanking_StraightFlush]
		x2 := h2.Cards.MapCardType[pb.HandRanking_StraightFlush]
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1[0].GetRank()
		point2 = x2[0].GetRank()
	case pb.HandRanking_FourOfAKind:
		x1 := h1.Cards.MapCardType[pb.HandRanking_FourOfAKind]
		x2 := h2.Cards.MapCardType[pb.HandRanking_FourOfAKind]
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1[0].GetRank()
		point2 = x2[0].GetRank()
	case pb.HandRanking_FullHouse:
		x1 := h1.Cards.MapCardType[pb.HandRanking_ThreeOfAKind]
		x2 := h2.Cards.MapCardType[pb.HandRanking_ThreeOfAKind]
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1[0].GetRank()
		point2 = x2[0].GetRank()
	case pb.HandRanking_Flush:
		x1 := h1.Cards.ListCard
		x2 := h2.Cards.ListCard
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1.GetMaxPointCard()
		point2 = x2.GetMaxPointCard()
	case pb.HandRanking_Straight:
		x1 := h1.Cards.ListCard
		x2 := h2.Cards.ListCard
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1.GetMaxPointCard()
		point2 = x2.GetMaxPointCard()
	case pb.HandRanking_ThreeOfAKind:
		x1 := h1.Cards.MapCardType[pb.HandRanking_ThreeOfAKind]
		x2 := h2.Cards.MapCardType[pb.HandRanking_ThreeOfAKind]
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1[0].GetRank()
		point2 = x2[0].GetRank()
	case pb.HandRanking_TwoPairs:
		x1 := h1.Cards.MapCardType[pb.HandRanking_TwoPairs]
		x2 := h2.Cards.MapCardType[pb.HandRanking_TwoPairs]
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}

		//
		point1 = x1[2].GetRank()
		point2 = x2[2].GetRank()
		if point1 == point2 {
			point1 = x1[0].GetRank()
			point2 = x2[0].GetRank()
		}
		if point1 == point2 {
			extraScore = h1.Cards.ListCard.CompareHighCard(h2.Cards.ListCard)
		}
	case pb.HandRanking_Pair:
		x1 := h1.Cards.MapCardType[pb.HandRanking_Pair]
		x2 := h2.Cards.MapCardType[pb.HandRanking_Pair]
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		point1 = x1[0].GetRank()
		point2 = x2[0].GetRank()
		if point1 == point2 {
			extraScore = h1.Cards.ListCard.CompareHighCard(h2.Cards.ListCard)
		}
	case pb.HandRanking_HighCard:
		x1 := h1.Cards.ListCard
		x2 := h2.Cards.ListCard
		if len(x1) == 0 || len(x1) != len(x2) {
			return resultPoint
		}
		extraScore = x1.CompareHighCard(x2)
	}

	if point1 > point2 {
		resultPoint++
		return resultPoint
	}
	if point1 < point2 {
		resultPoint--
		return resultPoint
	}
	if extraScore > 0 {
		resultPoint++
		return resultPoint
	}
	if extraScore < 0 {
		resultPoint--
		return resultPoint
	}

	return resultPoint
}

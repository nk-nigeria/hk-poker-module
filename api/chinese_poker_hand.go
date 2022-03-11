package api

import (
	"errors"
	"log"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

// Hand
// Contain all presence card
type Hand struct {
	cards   entity.ListCard
	ranking pb.HandRanking

	// frontHand  [3]*pb.Card
	// middleHand [5]*pb.Card
	// backHand   [5]*pb.Card

	// frontHandPoint  *HandPoint
	// middleHandPoint *HandPoint
	// backHandPoint   *HandPoint

	frontHand  *ChildHand
	middleHand *ChildHand
	backHand   *ChildHand

	burned    bool
	royalties bool
}
type ChildHand struct {
	Cards *HandCards
	Point *HandPoint
}

func (ch *ChildHand) calculatePoint() {
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

func NewHand(cards *pb.ListCard) (*Hand, error) {
	if cards == nil {
		h := &Hand{}
		return h, nil
	}
	listCard := make(entity.ListCard, 0, len(cards.Cards))
	// deep copy card
	for _, c := range cards.GetCards() {
		listCard = append(listCard, entity.NewCard(c.GetRank(), c.GetSuit()))
	}
	hand := &Hand{
		cards: listCard,
	}

	if hand.parse() != nil {
		return nil, errors.New("hand.new.error.invalid")
	}

	return hand, nil
}

func (h *Hand) GetCards() entity.ListCard {
	return h.cards
}

func (h *Hand) parse() error {
	cards := h.cards
	if len(cards) != MaxPresenceCard {
		return errors.New("hand.parse.error.invalid-len")
	}

	// h.frontHand[0] = cards[0]
	// h.frontHand[1] = cards[1]
	// h.frontHand[2] = cards[2]
	h.frontHand = NewChildHand(cards[:3])

	// h.middleHand[0] = cards[3]
	// h.middleHand[1] = cards[4]
	// h.middleHand[2] = cards[5]
	// h.middleHand[3] = cards[6]
	// h.middleHand[4] = cards[7]
	h.middleHand = NewChildHand(cards[3:8])

	// h.backHand[0] = cards[8]
	// h.backHand[1] = cards[9]
	// h.backHand[2] = cards[10]
	// h.backHand[3] = cards[11]
	// h.backHand[4] = cards[12]
	h.backHand = NewChildHand(cards[8:])

	return nil
}

func (h *Hand) calculatePoint() int {
	// Check royalties

	// h.calculatePointFrontHand()
	// h.calculatePointMiddleHand()
	// h.calculatePointBackHand()
	h.frontHand.calculatePoint()
	h.middleHand.calculatePoint()
	h.backHand.calculatePoint()

	// Check 3 flush
	// Check 3 straight

	return 0
}

func (h *Hand) calculatePointFrontHand() {
	var handCard *HandCards
	h.frontHand.Point, handCard = GetHandPoint(h.frontHand.Cards.ListCard)
	// copy(h.backHand.Child[:], sortedCard[:3])
	h.frontHand.Cards = handCard
}

func (h *Hand) calculatePointMiddleHand() {
	var handCard *HandCards
	h.middleHand.Point, handCard = GetHandPoint(h.middleHand.Cards.ListCard)
	h.middleHand.Cards = handCard
}

func (h *Hand) calculatePointBackHand() {
	var handCard *HandCards
	h.backHand.Point, handCard = GetHandPoint(h.middleHand.Cards.ListCard)
	h.backHand.Cards = handCard

}

func CompareHand(h1, h2 *Hand) *pb.CompareResult {
	result := &pb.CompareResult{}
	result.FrontFactor = int64(compareChildHand(h1.frontHand, h2.frontHand))
	log.Printf("result1 %d\n", result.FrontFactor)
	result.MiddleFactor = int64(compareChildHand(h1.middleHand, h2.middleHand))
	log.Printf("result2 %d\n", result.MiddleFactor)
	result.BackFactor = int64(compareChildHand(h1.backHand, h2.backHand))
	log.Printf("result3 %d\n", result.BackFactor)
	return result
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

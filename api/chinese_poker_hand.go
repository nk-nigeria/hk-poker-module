package api

import (
	"errors"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

// Hand
// Contain all presence card
type Hand struct {
	cards   *pb.ListCard
	ranking pb.HandRanking

	frontHand  [3]*pb.Card
	middleHand [5]*pb.Card
	backHand   [5]*pb.Card

	frontHandPoint  *HandPoint
	middleHandPoint *HandPoint
	backHandPoint   *HandPoint

	burned    bool
	royalties bool
}

func NewHand(cards *pb.ListCard) (*Hand, error) {
	if cards == nil {
		h := &Hand{}
		return h, nil
	}
	listCard := make([]*pb.Card, 0, len(cards.Cards))
	// deep copy card
	for _, c := range cards.GetCards() {
		x := pb.Card{
			Rank:   c.GetRank(),
			Suit:   c.GetSuit(),
			Status: c.GetStatus(),
		}
		listCard = append(listCard, &x)
	}
	hand := &Hand{
		cards: &pb.ListCard{
			Cards: cards.Cards,
		},
	}

	if hand.parse() != nil {
		return nil, errors.New("hand.new.error.invalid")
	}

	return hand, nil
}

func (h *Hand) GetCards() *pb.ListCard {
	return h.cards
}

func (h *Hand) parse() error {
	var cards = h.cards.GetCards()

	if len(cards) != MaxPresenceCard {
		return errors.New("hand.parse.error.invalid-len")
	}

	// h.frontHand[0] = cards[0]
	// h.frontHand[1] = cards[1]
	// h.frontHand[2] = cards[2]
	copy(h.frontHand[:], cards[:3])

	// h.middleHand[0] = cards[3]
	// h.middleHand[1] = cards[4]
	// h.middleHand[2] = cards[5]
	// h.middleHand[3] = cards[6]
	// h.middleHand[4] = cards[7]

	copy(h.frontHand[:], cards[3:8])

	// h.backHand[0] = cards[8]
	// h.backHand[1] = cards[9]
	// h.backHand[2] = cards[10]
	// h.backHand[3] = cards[11]
	// h.backHand[4] = cards[12]
	copy(h.frontHand[:], cards[8:])

	return nil
}

func (h *Hand) calculatePoint() int {
	// Check royalties

	h.calculatePointFrontHand()
	h.calculatePointMiddleHand()
	h.calculatePointBackHand()

	// Check 3 flush
	// Check 3 straight

	return 0
}

func (h *Hand) calculatePointFrontHand() {
	var sortedCard []*pb.Card
	h.frontHandPoint, sortedCard = GetHandPoint(h.frontHand[:])
	h.frontHand[0] = sortedCard[0]
	h.frontHand[1] = sortedCard[1]
	h.frontHand[2] = sortedCard[2]
}

func (h *Hand) calculatePointMiddleHand() {
	var sortedCard []*pb.Card
	h.middleHandPoint, sortedCard = GetHandPoint(h.middleHand[:])
	h.middleHand[0] = sortedCard[0]
	h.middleHand[1] = sortedCard[1]
	h.middleHand[2] = sortedCard[2]
	h.middleHand[3] = sortedCard[3]
	h.middleHand[4] = sortedCard[4]
}

func (h *Hand) calculatePointBackHand() {
	var sortedCard []*pb.Card
	h.backHandPoint, sortedCard = GetHandPoint(h.backHand[:])
	h.backHand[0] = sortedCard[0]
	h.backHand[1] = sortedCard[1]
	h.backHand[2] = sortedCard[2]
	h.backHand[3] = sortedCard[3]
	h.backHand[4] = sortedCard[4]
}

func CompareHand(h1, h2 *Hand) *pb.CompareResult {
	h1.calculatePoint()
	h2.calculatePoint()
	result := &pb.CompareResult{}
	if entity.GetHandRankingPoint(h1.frontHandPoint.rankingType) > entity.GetHandRankingPoint(h2.frontHandPoint.rankingType) {
		result.FrontFactor++
	} else if entity.GetHandRankingPoint(h1.frontHandPoint.rankingType) < entity.GetHandRankingPoint(h2.frontHandPoint.rankingType) {
		result.FrontFactor--
	} else {
		for _, card1 := range h1.frontHand {
			for _, card2 := range h2.frontHand {
				if entity.GetCardRankPoint(card1.Rank) > entity.GetCardRankPoint(card2.Rank) {

				}
			}
		}
	}

	return result
}

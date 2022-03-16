package api

import (
	"errors"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

// Hand
// Contain all presence card
type Hand struct {
	cards   entity.ListCard
	ranking pb.HandRanking

	frontHand  *ChildHand
	middleHand *ChildHand
	backHand   *ChildHand

	burned    bool
	royalties bool
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

	h.frontHand = NewChildHand(cards[:3])
	h.middleHand = NewChildHand(cards[3:8])
	h.backHand = NewChildHand(cards[8:])

	return nil
}

func (h *Hand) calculatePoint() int {
	// Check royalties

	h.frontHand.calculatePoint()
	h.middleHand.calculatePoint()
	h.backHand.calculatePoint()

	// Check 3 flush
	// Check 3 straight

	return 0
}

func (h *Hand) calculatePointFrontHand() {
	var handCard *HandCards
	h.frontHand.Point, handCard = CaculatorPoint(h.frontHand.Cards.ListCard)
	// copy(h.backHand.Child[:], sortedCard[:3])
	h.frontHand.Cards = handCard
}

func (h *Hand) calculatePointMiddleHand() {
	var handCard *HandCards
	h.middleHand.Point, handCard = CaculatorPoint(h.middleHand.Cards.ListCard)
	h.middleHand.Cards = handCard
}

func (h *Hand) calculatePointBackHand() {
	var handCard *HandCards
	h.backHand.Point, handCard = CaculatorPoint(h.middleHand.Cards.ListCard)
	h.backHand.Cards = handCard
}

func CompareHand(h1, h2 *Hand) *pb.ComparisonResult {
	result := &pb.ComparisonResult{}
	result.FrontFactor = int64(h1.frontHand.CompareHand(h2.frontHand))
	result.MiddleFactor = int64(h1.middleHand.CompareHand(h2.middleHand))
	result.BackFactor = int64(h1.backHand.CompareHand(h2.backHand))
	return result
}

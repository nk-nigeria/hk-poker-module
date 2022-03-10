package api

import (
	"errors"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

// Hand
// Contain all presence card
type Hand struct {
	cards   pb.ListCard
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

func NewHand(cards pb.ListCard) (*Hand, error) {
	hand := &Hand{
		cards: cards,
	}

	if hand.parse() != nil {
		return nil, errors.New("hand.new.error.invalid")
	}

	return hand, nil
}

func (h *Hand) GetCards() pb.ListCard {
	return h.cards
}

func (h *Hand) parse() error {
	var cards = h.cards.GetCards()

	if len(cards) != MaxPresenceCard {
		return errors.New("hand.parse.error.invalid-len")
	}

	h.frontHand[0] = cards[0]
	h.frontHand[1] = cards[1]
	h.frontHand[2] = cards[2]

	h.middleHand[0] = cards[3]
	h.middleHand[1] = cards[4]
	h.middleHand[2] = cards[5]
	h.middleHand[3] = cards[6]
	h.middleHand[4] = cards[7]

	h.backHand[0] = cards[8]
	h.backHand[1] = cards[9]
	h.backHand[2] = cards[10]
	h.backHand[3] = cards[11]
	h.backHand[4] = cards[12]

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
	h.frontHandPoint, _ = GetHandPoint(h.frontHand[:])
}

func (h *Hand) calculatePointMiddleHand() {
	h.middleHandPoint, _ = GetHandPoint(h.middleHand[:])
}

func (h *Hand) calculatePointBackHand() {
	h.backHandPoint, _ = GetHandPoint(h.backHand[:])
}

func CompareHand(h1, h2 *Hand) int {

	return 0
}

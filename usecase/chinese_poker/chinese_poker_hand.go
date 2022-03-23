package chinese_poker

import (
	"errors"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

type Result struct {
	FrontFactor       int `json:"front_factor"`
	MiddleFactor      int `json:"middle_factor"`
	BackFactor        int `json:"back_factor"`
	FrontBonusFactor  int `json:"front_bonus_factor"`
	MiddleBonusFactor int `json:"middle_bonus_factor"`
	BackBonusFactor   int `json:"back_bonus_factor"`
	NaturalFactor     int `json:"natural_factor"`
	ScoopFactor       int `json:"scoop_factor"`
}

// ComparisonResult
type ComparisonResult struct {
	r1 Result `json:"r1"`
	r2 Result `json:"r1"`
}

// Hand
// Contain all presence card
type Hand struct {
	cards   entity.ListCard
	ranking pb.HandRanking

	frontHand  *ChildHand
	middleHand *ChildHand
	backHand   *ChildHand

	naturalPoint *HandPoint
	pointType    pb.PointType
	calculated   bool
}

func NewHand(cards *pb.ListCard) (*Hand, error) {
	if cards == nil {
		h := &Hand{
			calculated: false,
		}
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

func (h Hand) GetCards() entity.ListCard {
	return h.cards
}

func (h Hand) IsNatural() bool {
	return h.pointType == pb.PointType_Point_Natural
}

func (h Hand) IsMisSet() bool {
	return h.pointType == pb.PointType_Point_Mis_Set
}

func (h Hand) IsNormal() bool {
	return h.pointType == pb.PointType_Point_Normal
}

func (h *Hand) parse() error {
	cards := h.cards
	if len(cards) != MaxPresenceCard {
		return errors.New("hand.parse.error.invalid-len")
	}

	h.frontHand = NewChildHand(cards[:3], kFronHand)
	h.middleHand = NewChildHand(cards[3:8], kMidHand)
	h.backHand = NewChildHand(cards[8:], kBackHand)

	return nil
}

func (h *Hand) calculatePoint() error {
	if h.calculated {
		return errors.New("hand.calculate.already")
	}
	defer func() {
		// mark as already calculated
		h.calculated = true
	}()

	// check cards naturals
	handPoint, natural := CheckNaturalCards(h)
	if natural {
		h.pointType = pb.PointType_Point_Natural
		h.naturalPoint = handPoint
		return nil
	}

	// calculate hand by hand
	h.frontHand.calculatePoint()
	h.middleHand.calculatePoint()
	h.backHand.calculatePoint()

	// check mis set
	if IsMisSets(h) {
		h.pointType = pb.PointType_Point_Mis_Set
		return nil
	}

	// check hand naturals
	handPoint, natural = CheckNaturalHands(h)
	if natural {
		h.pointType = pb.PointType_Point_Natural
		h.naturalPoint = handPoint
		return nil
	}

	return nil
}

func (h Hand) GetPointResult() *pb.PointResult {
	result := &pb.PointResult{
		Type: h.pointType,
	}

	switch h.pointType {
	case pb.PointType_Point_Normal:
		result = &pb.PointResult{
			Front:  h.frontHand.Point.ToHandResultPB(),
			Middle: h.middleHand.Point.ToHandResultPB(),
			Back:   h.backHand.Point.ToHandResultPB(),
		}
	case pb.PointType_Point_Natural:
		result = &pb.PointResult{
			Natural: h.naturalPoint.ToHandResultPB(),
		}
	case pb.PointType_Point_Mis_Set:

	}

	return result
}

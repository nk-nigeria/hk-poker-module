package chinese_poker

import (
	"errors"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

// ComparisonResult
type ComparisonResult struct {
	UserId            string
	FrontFactor       int64
	MiddleFactor      int64
	BackFactor        int64
	FrontBonusFactor  int64
	MiddleBonusFactor int64
	BackBonusFactor   int64
	ScoopFactor       int64
}

// Hand
// Contain all presence card
type Hand struct {
	cards   entity.ListCard
	ranking pb.HandRanking

	frontHand  *ChildHand
	middleHand *ChildHand
	backHand   *ChildHand

	pointType pb.PointType
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

func (h Hand) GetCards() entity.ListCard {
	return h.cards
}

func (h Hand) IsNatural() bool {
	return h.pointType > pb.PointType_Point_NaturalMin && h.pointType < pb.PointType_Point_NaturalMax
}

func (h Hand) IsMisSet() bool {
	return h.pointType == pb.PointType_Point_Mis_Set
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

func naturalTypeToPointType(ntype pb.NaturalRanking) pb.PointType {
	switch ntype {
	case pb.NaturalRanking_Dragon:
		return pb.PointType_Point_NaturalDragon
	case pb.NaturalRanking_CleanDragon:
		return pb.PointType_Point_NaturalCleanDragon
	case pb.NaturalRanking_SixPairs:
		return pb.PointType_Point_NaturalSixPairs
	case pb.NaturalRanking_FullColors:
		return pb.PointType_Point_NaturalFullColors
	case pb.NaturalRanking_ThreeStraights:
		return pb.PointType_Point_NaturalThreeStraights
	case pb.NaturalRanking_ThreeOfFlushes:
		return pb.PointType_Point_NaturalThreeOfFlushes
	}

	return pb.PointType_Point_Normal
}

func (h *Hand) calculatePoint() error {
	// check cards naturals
	natural, naturalType := CheckNaturalCards(h)
	if natural {
		h.pointType = naturalTypeToPointType(naturalType)
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
	natural, naturalType = CheckNaturalHands(h)
	if natural {
		h.pointType = naturalTypeToPointType(naturalType)
		return nil
	}

	return nil
}

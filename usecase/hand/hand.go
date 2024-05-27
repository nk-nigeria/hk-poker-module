package hand

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var (
	kWinScoop  = 1
	kLoseScoop = -1

	kWinMisset  = 2
	kLoseMisset = -2
)

var (
	kFronHand = 0
	kMidHand  = 1
	kBackHand = 2
)

type Result struct {
	FrontFactor       int `json:"front_factor"`
	MiddleFactor      int `json:"middle_factor"`
	BackFactor        int `json:"back_factor"`
	FrontBonusFactor  int `json:"front_bonus_factor"`
	MiddleBonusFactor int `json:"middle_bonus_factor"`
	BackBonusFactor   int `json:"back_bonus_factor"`
	NaturalFactor     int `json:"natural_factor"`
	BonusFactor       int `json:"bonus_factor"`
	Scoop             int `json:"scoop"`
}

type ComparisonResult struct {
	r1      Result `json:"r1"`
	r2      Result `json:"r2"`
	bonuses []*pb.HandBonus
}

func (r *ComparisonResult) swap() {
	tmp := r.r1
	r.r1 = r.r2
	r.r2 = tmp
}

func (r *ComparisonResult) addHandBonus(win, lose string, bonusType pb.HandBonusType, factor int64) {
	r.bonuses = append(r.bonuses, &pb.HandBonus{
		Win:    win,
		Lose:   lose,
		Type:   bonusType,
		Factor: factor,
	})
}

func (r ComparisonResult) GetR1() Result {
	return r.r1
}

func (r ComparisonResult) GetR2() Result {
	return r.r2
}

func (r ComparisonResult) GetBonuses() []*pb.HandBonus {
	return r.bonuses
}

// Hand
// Contain all presence card
type Hand struct {
	cards   entity.ListCard
	ranking pb.HandRanking

	frontHand  *ChildHand // :3
	middleHand *ChildHand // 3:8
	backHand   *ChildHand // 8:13

	naturalPoint *HandPoint
	pointType    pb.PointType
	calculated   bool

	owner   string
	jackpot bool
}

func (h Hand) String() string {
	var str string
	str += fmt.Sprintf("front: %s\n", h.frontHand)
	str += fmt.Sprintf("middle: %s\n", h.middleHand)
	str += fmt.Sprintf("back: %s\n", h.backHand)

	return str
}

func NewHandFromPb(cards *pb.ListCard) (*Hand, error) {
	if cards == nil {
		h := &Hand{
			calculated: false,
		}
		return h, nil
	}
	listCard := make(entity.ListCard, 0, len(cards.Cards))
	// deep copy card
	for _, c := range cards.GetCards() {
		listCard = append(listCard, entity.NewCardFromPb(c.GetRank(), c.GetSuit()))
	}
	hand := &Hand{
		cards: listCard,
	}

	if hand.parse() != nil {
		return nil, errors.New("hand.new.error.invalid")
	}

	return hand, nil
}

func NewHand(cards entity.ListCard) (*Hand, error) {
	if cards == nil {
		h := &Hand{
			calculated: false,
		}
		return h, nil
	}

	hand := &Hand{
		cards: cards,
	}

	if hand.parse() != nil {
		return nil, errors.New("hand.new.error.invalid")
	}

	return hand, nil
}

func (h *Hand) SetOwner(owner string) {
	h.owner = owner
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
	if len(cards) != entity.MaxPresenceCard {
		return errors.New("hand.parse.error.invalid-len")
	}

	h.frontHand = NewChildHand(cards[:3], kFronHand)
	h.middleHand = NewChildHand(cards[3:8], kMidHand)
	h.backHand = NewChildHand(cards[8:], kBackHand)

	return nil
}

func (h *Hand) CalculatePoint() error {
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

func (h *Hand) GetPointResult() *pb.PointResult {
	h.CalculatePoint()

	// result := &pb.PointResult{
	// 	Type: h.pointType,
	// }

	var result *pb.PointResult
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
	// case pb.PointType_Point_Mis_Set:

	// }
	default:
		{
			result = &pb.PointResult{}
		}
	}
	result.Type = h.pointType
	h.jackpot = CheckJackpot(h.middleHand) || CheckJackpot(h.backHand)
	return result
}

func (h *Hand) IsJackpot() bool {
	return h.jackpot
}

func (h *Hand) AutoOrgCards() *autoHand {
	if h.IsNatural() {
		return nil
	}
	// priority
	// Straight Flush -> Four of a Kind -> Full House -> Flush -> Straight -> Three of a Kind -> Two Pair -> One Pair -> High Card
	// currentHand := h.backHand
	autoHand := NewAutoHand(h.cards)
	handIndex := kBackHand
	// handArr := []*ChildHand{h.backHand, h.middleHand, h.frontHand}
	type FnF func() []entity.ListCard
	listFn := make([]FnF, 0)
	listFn = append(listFn, autoHand.FindStraighFlush, autoHand.FindFourKind, autoHand.FindFullHouse, autoHand.FindFlush, autoHand.FindStraigh)
	for i := 0; i < 2; i++ {
		for _, fn := range listFn {
			if handIndex < 0 {
				break
			}
			list := autoHand.PreferCardsNotTakeDouble(fn()...)
			if list != nil {
				// handArr[handIndex] = NewChildHand(list, handIndex)
				if len(list) == 4 {
					// highcard := autoHand.FindHighCard()
					// sort.Slice(h.cards, func(i, j int) bool {
					// 	return h.cards[i].GetRank() < h.cards[j].GetRank()
					// })
					highcard := autoHand.GetHighCardWithCond(1, list...)
					list = append(list, highcard[:1]...)
				}
				h.setHand(list, handIndex)
				autoHand.TakeCard(list...)
				handIndex--
				continue
				// if handIndex <= 0 {
				// 	currentHand = handArr[handIndex]
				// }
			}
		}
	}

	for i := 0; i < 3; i++ {
		if handIndex < 0 {
			break
		} // three of a kind
		if handIndex > 0 {
			threeOfKind := autoHand.FindThreeKind()
			if len(threeOfKind) > 0 {
				list := make(entity.ListCard, 0)
				list = append(list, threeOfKind[0]...)
				highcard := autoHand.GetHighCardWithCond(2, list...)
				sort.Slice(h.cards, func(i, j int) bool {
					return h.cards[i].GetRank() < h.cards[j].GetRank()
				})
				if len(highcard) < 2 {
					// panic(h.cards)
					return autoHand
				}
				list = append(list, highcard[:2]...)
				// handArr[handIndex] = NewChildHand(list, handIndex)
				h.setHand(list, handIndex)
				autoHand.TakeCard(list...)
				handIndex--
				continue
			}
		}
		// two pair
		if handIndex > 0 {
			twoPair := autoHand.FindTwoPair()
			if len(twoPair) >= 2 {
				list := make([]entity.Card, 0)
				// for _, pair := range twoPair[:2] {
				// 	list = append(list, pair...)
				// }
				list = append(list, twoPair[0]...)
				list = append(list, twoPair[len(twoPair)-1]...)
				highcard := autoHand.GetHighCardWithCond(1, list...)
				if len(highcard) < 1 {
					// panic(h.cards)
					return autoHand
				}
				list = append(list, highcard[:1]...)
				// handArr[handIndex] = NewChildHand(list, handIndex)
				h.setHand(list, handIndex)
				autoHand.TakeCard(list...)
				handIndex--
				continue
			}
		}
		// pair
		if handIndex >= 0 {
			pair := autoHand.PreferCardsNotTakeDouble(autoHand.FindPair()...)
			if len(pair) > 0 {
				numMoreCard := 3
				if handIndex == kFronHand {
					numMoreCard = 1
				}
				list := make(entity.ListCard, 0)
				list = append(list, pair...)
				highcard := autoHand.GetHighCardWithCond(numMoreCard, list...)
				if len(highcard) >= numMoreCard {
					list = append(list, highcard[:numMoreCard]...)
					// handArr[handIndex] = NewChildHand(list, handIndex)
					h.setHand(list, handIndex)
					autoHand.TakeCard(list...)
					handIndex--
					continue
				}
			}
		}
	}

	if handIndex == kBackHand {
		highCards := autoHand.GetHighCardWithCond(5)
		h.backHand = NewChildHand(highCards[:5], 2)
		autoHand.TakeCard(highCards[:5]...)
		handIndex--
	}
	if handIndex == kMidHand {
		highCards := autoHand.GetHighCardWithCond(5)
		h.middleHand = NewChildHand(highCards[:5], 1)
		autoHand.TakeCard(highCards[:5]...)
		handIndex--
	}
	if handIndex == kFronHand {
		highCards := autoHand.GetHighCardWithCond(3)
		h.frontHand = NewChildHand(highCards[:3], 0)
		autoHand.TakeCard(highCards[:3]...)
		// handIndex--
	}

	h.cards = make(entity.ListCard, 0)
	h.cards = append(h.cards, h.frontHand.Cards...)
	h.cards = append(h.cards, h.middleHand.Cards...)
	h.cards = append(h.cards, h.backHand.Cards...)
	return autoHand
}

func (h *Hand) setHand(ml entity.ListCard, handType int) {
	switch handType {
	case kFronHand:
		h.frontHand = NewChildHand(ml, handType)
	case kMidHand:
		h.middleHand = NewChildHand(ml, handType)
	case kBackHand:
		h.backHand = NewChildHand(ml, handType)
	}
}

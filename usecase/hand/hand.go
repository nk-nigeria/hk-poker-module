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

func (h *Hand) AutoOrgCards() {
	if h.IsNatural() {
		return
	}
	// priority
	// Straight Flush -> Four of a Kind -> Full House -> Flush -> Straight -> Three of a Kind -> Two Pair -> One Pair -> High Card
	// currentHand := h.backHand
	autoHand := NewAutoHand(h.cards)
	handIndex := kBackHand
	handArr := []*ChildHand{h.backHand, h.middleHand, h.frontHand}
	type FnF func() []entity.ListCard
	fn := make([]FnF, 0)
	listFn := append(fn, autoHand.FindStraighFlush, autoHand.FindFourKind, autoHand.FindFullHouse, autoHand.FindFlush, autoHand.FindStraigh)
	for _, fn := range listFn {
		if handIndex < 0 {
			break
		}
		list := autoHand.PreferCardsNotTakeDouble(fn()...)
		if list != nil {
			handArr[handIndex] = NewChildHand(list, handIndex)

			autoHand.TakeCard(list...)
			handIndex--
			// if handIndex <= 0 {
			// 	currentHand = handArr[handIndex]
			// }
		}
	}
	// three of a kind
	if handIndex > 0 {
		threeOfKind := autoHand.PreferCardsNotTakeDouble(autoHand.FindThreeKind()...)
		if len(threeOfKind) > 0 {
			highcard := autoHand.FindHighCard()
			sort.Slice(h.cards, func(i, j int) bool {
				return h.cards[i].GetRank() < h.cards[j].GetRank()
			})
			list := append(threeOfKind, highcard[:2]...)
			handArr[handIndex] = NewChildHand(list, handIndex)
			autoHand.TakeCard(list...)
			handIndex--
		}
	}
	// two pair
	if handIndex > 0 {
		twoPair := autoHand.FindTwoPair()
		if len(twoPair) > 0 {
			highcard := autoHand.FindHighCard()
			sort.Slice(h.cards, func(i, j int) bool {
				return h.cards[i].GetRank() < h.cards[j].GetRank()
			})
			list := make([]entity.Card, 0)
			for _, pair := range twoPair {
				list = append(list, pair...)
			}
			list = append(list, highcard[:1]...)
			handArr[handIndex] = NewChildHand(list, handIndex)
			autoHand.TakeCard(list...)
			handIndex--
		}
	}
	// pair
	if handIndex >= 0 {
		pair := autoHand.PreferCardsNotTakeDouble(autoHand.FindPair()...)
		if len(pair) > 0 {
			highcard := autoHand.FindHighCard()
			sort.Slice(h.cards, func(i, j int) bool {
				return h.cards[i].GetRank() < h.cards[j].GetRank()
			})
			numMoreCard := 3
			if handIndex == kFronHand {
				numMoreCard = 1
			}
			if len(highcard) >= numMoreCard {
				list := append(pair, highcard[:numMoreCard]...)
				handArr[handIndex] = NewChildHand(list, handIndex)
				autoHand.TakeCard(list...)
				handIndex--
			}
		}
	}
	highCards := autoHand.FindHighCard()
	if len(highCards) > 0 {
		if handIndex == kBackHand && len(highCards) >= 5 {
			h.backHand = NewChildHand(highCards[:5], 2)
			handIndex--
		}

		if handIndex == kMidHand && len(highCards) >= 5 {
			h.middleHand = NewChildHand(highCards[:5], 1)
			highCards = highCards[5:]
			handIndex--
		}
		if handIndex == kFronHand && len(highCards) >= 3 {
			h.frontHand = NewChildHand(highCards[:3], 0)
			// highCards = highCards[3:]
			// handIndex--
		}
	}
	h.cards = make(entity.ListCard, 0)
	h.cards = append(h.cards, h.frontHand.Cards...)
	h.cards = append(h.cards, h.middleHand.Cards...)
	h.cards = append(h.cards, h.backHand.Cards...)
}

func (h *Hand) AutoSortForBest() {
	cards := h.cards.Clone()
	// srt from A->2
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetRank() > cards[j].GetRank()
	})
	trackCardTake := make(map[entity.Card]struct{})
	// tần suất lá bài xuất hiện
	cardsCount := make(map[entity.Card]int)
	for _, c := range cards {
		cardsCount[c]++
	}
	// xóa lá bài chỉ xuất hiện 1 lần
	for k, v := range cardsCount {
		if v <= 1 {
			delete(cardsCount, k)
		}
	}
	// handCard := make([]entity.ListCard, 0)
	// tìm lá bài cùng màu, đồng chất
	fnFindStraighFlush := func() []entity.ListCard {
		cardsSameColor := make(map[uint8]entity.ListCard)
		for _, c := range cards {
			if _, exist := trackCardTake[c]; exist {
				continue
			}
			list := cardsSameColor[c.GetRank()]
			list = append(list, c)
			cardsSameColor[c.GetRank()] = list
		}
		// remove card same color <5 card
		listStraighFlush := make([]entity.ListCard, 0)
		for _, v := range cardsSameColor {
			if len(v) < 5 {
				continue
			}
			for i := 0; i < len(v)-5; i++ {
				ml := v[:5+i]
				binCards := entity.NewBinListCards(ml)
				_, ok := CheckStraightFlush(binCards)
				if !ok {
					continue
				}
				listStraighFlush = append(listStraighFlush, ml)
			}
		}
		return listStraighFlush
	}

	// Năm lá bài cùng màu, đồng chất nhưng không cùng một chuỗi số
	fnFindFlush := func() []entity.ListCard {
		cardsSameColor := make(map[uint8]entity.ListCard)
		for _, c := range cards {
			if _, exist := trackCardTake[c]; exist {
				continue
			}
			list := cardsSameColor[c.GetRank()]
			list = append(list, c)
			cardsSameColor[c.GetRank()] = list
		}
		// remove card same color <5 card
		listFlush := make([]entity.ListCard, 0)
		for _, v := range cardsSameColor {
			if len(v) < 5 {
				continue
			}
			for i := 0; i < len(v)-5; i++ {
				ml := v[:5+i]
				binCards := entity.NewBinListCards(ml)
				_, ok := CheckFlush(binCards)
				if !ok {
					continue
				}
				listFlush = append(listFlush, ml)
			}
		}
		return listFlush
	}

	// ưu tiên lấy list card  ko chưa card có thể tạo đôi
	fnSelect := func(ml []entity.ListCard) entity.ListCard {
		if len(ml) == 0 {
			return nil
		}
		if len(ml) == 1 {
			return ml[0]
		}
		// ưu tiên straigh flush ko chứa card tạo đôi
		listMin := ml[0]
		minCount := 0
		for _, ml := range ml {
			count := 0
			for _, c := range ml {
				count += cardsCount[c]
			}
			if count == 0 {
				return ml
			}
			if count < minCount {
				listMin = ml
				minCount = count
			}
		}
		return listMin
	}
	_ = fnFindStraighFlush()
	_ = fnSelect(nil)
	_ = fnFindFlush()
}

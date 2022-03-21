package chinese_poker

import (
	"errors"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

var cleanWinChecker map[int64]func(entity.ListCard) (*HandCards, bool)

func init() {
	cleanWinChecker = make(map[int64]func(entity.ListCard) (*HandCards, bool))
	cleanWinChecker[entity.WIN_TYPE_WIN_CLEAN_DRAGON] = IsCleanDragon
	cleanWinChecker[entity.WIN_TYPE_WIN_DRAGON] = IsDragon
	cleanWinChecker[entity.WIN_TYPE_WIN_FIVE_PAIR_THREE_OF_A_KIND] = IsFivePairThreeOfAKind
	cleanWinChecker[entity.WIN_TYPE_WIN_SIX_AND_A_HALF_PAIRS] = IsSixAndAHalfPairs
	cleanWinChecker[entity.WIN_TYPE_WIN_THREE_STRAIGHT_FLUSH] = IsThreeStraightFlush
	cleanWinChecker[entity.WIN_TYPE_WIN_THREE_STRAIGHT] = IsThreeStraight
	cleanWinChecker[entity.WIN_TYPE_WIN_THREE_FLUSH] = IsThreeFlushes
	cleanWinChecker[entity.WIN_TYPE_WIN_FULL_COLORED] = IsFullColored
}

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

	WinType       int64
	CleanWinBonus int64
	ScoopBonus    int64
}

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

func calculateBonus(result *ComparisonResult) *ComparisonResult {
	t := result.WinType
	if t&entity.WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND != 0 {
		result.FrontBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND))
	}
	if t&entity.WIN_TYPE_WIN_MID_FULL_HOUSE != 0 {
		result.MiddleBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_MID_FULL_HOUSE))
	}
	if t&entity.WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND != 0 {
		result.BackBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND))
	}
	if t&entity.WIN_TYPE_WIN_MID_FOUR_OF_A_KIND != 0 {
		result.MiddleBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_MID_FOUR_OF_A_KIND))
	}
	if t&entity.WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH != 0 {
		result.BackBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH))
	}
	if t&entity.WIN_TYPE_WIN_MID_STRAIGHT_FLUSH != 0 {
		result.MiddleBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_MID_STRAIGHT_FLUSH))
	}

	// Scoop check
	if result.FrontFactor > 0 && result.MiddleFactor > 0 && result.BackFactor > 0 {
		result.ScoopFactor = 1

		result.ScoopBonus = int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_SCOOP))
	}

	return result
}

func (h *Hand) CompareHand(h2 *Hand) *ComparisonResult {
	result := ComparisonResult{
		WinType: entity.WinType_WIN_TYPE_UNSPECIFIED,
	}
	// check clean win
	for k, checkerFn := range cleanWinChecker {
		l1, isHand1CleanWin := checkerFn(h.GetCards())
		if isHand1CleanWin {
			if k == entity.WIN_TYPE_WIN_CLEAN_DRAGON ||
				k == entity.WIN_TYPE_WIN_DRAGON ||
				k == entity.WIN_TYPE_WIN_THREE_STRAIGHT_FLUSH ||
				k == entity.WIN_TYPE_WIN_FULL_COLORED {
				result.WinType = k
				result.CleanWinBonus = int64(entity.GetWinFactorBonus(result.WinType))
				return &result
			}
		}
		l2, isHand2CleanWin := checkerFn(h2.GetCards())
		if !isHand1CleanWin && !isHand2CleanWin {
			continue
		}
		result.WinType = k
		result.CleanWinBonus = int64(entity.GetWinFactorBonus(result.WinType))
		if isHand1CleanWin != isHand2CleanWin {
			if isHand1CleanWin {
				return &result
			}
			result.CleanWinBonus = -result.CleanWinBonus
			return &result
		}
		switch k {
		// 5 đôi 1 xám: bài có 5 đôi và 1 xám cô. Giống nhau so sánh đến lá lớn nhất trong xám.
		case entity.WIN_TYPE_WIN_FIVE_PAIR_THREE_OF_A_KIND:
			x1 := l1.MapCardType[pb.HandRanking_ThreeOfAKind]
			x2 := l2.MapCardType[pb.HandRanking_ThreeOfAKind]
			if x1[len(x1)-1].GetRank() < x2[len(x2)-1].GetRank() {
				result.CleanWinBonus = -result.CleanWinBonus
			}
			return &result
		// Lục phé bôn: bài có 6 đôi và 1 mậu thầu. Giống nhau so đến đôi cao nhất.
		case entity.WIN_TYPE_WIN_SIX_AND_A_HALF_PAIRS:
			x1 := l1.MapCardType[pb.HandRanking_TwoPairs]
			x2 := l2.MapCardType[pb.HandRanking_TwoPairs]
			for i := len(x1) - 1; i >= 0; i -= 2 {
				r1 := x1[i].GetRank()
				r2 := x2[i].GetRank()
				if r1 == r2 {
					continue
				}
				if r1 < r2 {
					result.CleanWinBonus = -result.CleanWinBonus
				}
				return &result
			}
		case entity.WIN_TYPE_WIN_THREE_STRAIGHT:
			x1 := l1.MapCardType[pb.HandRanking_Straight]
			x2 := l2.MapCardType[pb.HandRanking_Straight]
			arr1 := make([]entity.ListCard, 3)
			arr1 = append(arr1, x1[8:])
			arr1 = append(arr1, x1[3:8])
			arr1 = append(arr1, x1[:3])

			arr2 := make([]entity.ListCard, 3)
			arr2 = append(arr2, x2[8:])
			arr2 = append(arr2, x2[3:8])
			arr2 = append(arr2, x2[:3])

			for i := 0; i < 3; i++ {
				compare := arr1[i].CompareHighCard(arr2[i])
				if compare == 0 {
					continue
				}
				if compare < 0 {
					result.CleanWinBonus = -result.CleanWinBonus
				}
				return &result
			}
			// draw
			result.CleanWinBonus = 0
			return &result

		case entity.WIN_TYPE_WIN_THREE_FLUSH:
			x1 := l1.MapCardType[pb.HandRanking_Flush]
			x2 := l2.MapCardType[pb.HandRanking_Flush]
			arr1 := make([]entity.ListCard, 3)
			arr1 = append(arr1, x1[8:])
			arr1 = append(arr1, x1[3:8])
			arr1 = append(arr1, x1[:3])

			arr2 := make([]entity.ListCard, 3)
			arr2 = append(arr2, x2[8:])
			arr2 = append(arr2, x2[3:8])
			arr2 = append(arr2, x2[:3])

			for i := 0; i < 3; i++ {
				compare := arr1[i].CompareHighCard(arr2[i])
				if compare == 0 {
					continue
				}
				if compare < 0 {
					result.CleanWinBonus = -result.CleanWinBonus
				}
				return &result
			}
			// draw
			result.CleanWinBonus = 0
			return &result
		}
		// draw, comapre child
	}
	return &result
}

func CompareHand(h1, h2 *Hand) *ComparisonResult {
	result := h1.CompareHand(h2)
	if result.WinType != 0 {
		return result
	}

	//  chi dau
	result.BackFactor = int64(h1.backHand.CompareHand(h2.backHand))
	if result.BackFactor > 0 {
		r := h1.backHand.Point.rankingType
		switch r {
		case pb.HandRanking_FourOfAKind:
			result.WinType |= entity.WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND
		case pb.HandRanking_StraightFlush:
			result.WinType |= entity.WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH
		}

	}
	// chi giua
	result.MiddleFactor = int64(h1.middleHand.CompareHand(h2.middleHand))
	if result.MiddleFactor > 0 {
		r := h1.middleHand.Point.rankingType
		switch r {
		case pb.HandRanking_FullHouse:
			result.WinType |= entity.WIN_TYPE_WIN_MID_FULL_HOUSE
		case pb.HandRanking_FourOfAKind:
			if h1.backHand.Point.rankingType == pb.HandRanking_FourOfAKind {
				result.WinType |= entity.WIN_TYPE_WIN_MID_FOUR_OF_A_KIND
			}
		case pb.HandRanking_StraightFlush:
			if h1.backHand.Point.rankingType == pb.HandRanking_StraightFlush {
				result.WinType |= entity.WIN_TYPE_WIN_MID_STRAIGHT_FLUSH
			}
		}
	}
	// chi cuoi
	result.FrontFactor = int64(h1.frontHand.CompareHand(h2.frontHand))
	if result.FrontBonusFactor > 0 && h1.frontHand.Point.rankingType == pb.HandRanking_ThreeOfAKind {
		result.WinType |= entity.WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND
	}
	return calculateBonus(result)
}

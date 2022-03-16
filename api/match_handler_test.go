package api

import (
	"testing"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/stretchr/testify/assert"
)

// Thùng phá sảnh (en: Straight Flush) vs Thùng phá sảnh (en: Straight Flush)
// Same level card
func TestCompareBasicStraightFlushVsStraightFlushDraw(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng phá sảnh (en: Straight Flush)
// list card 1 higher
func TestCompareBasicStraightFlushHigherStraightFlush(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng phá sảnh (en: Straight Flush)
// list card 1 lower
func TestCompareBasicStraightFlushLowerStraightFlushLower(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(-1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Tứ quý (en: Four of a Kind)
func TestCompareBasicStraightFlushVsFourOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Cù lũ (en: Full House
func TestCompareBasicStraightFlushVsFullhouse(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thùng (en: Flush)
func TestCompareBasicStraightFlushVsFlush(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Sảnh (en: Straight)
func TestCompareBasicStraightFlushVsStraight(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Xám chi/Xám cô (en: Three of a Kind)
func TestCompareBasicStraightFlushVsThreeOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Thú (en: Two Pairs)
func TestCompareBasicStraightFlushVsTwoPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Đôi (en: Pair)
func TestCompareBasicStraightFlushVsPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng phá sảnh (en: Straight Flush) vs Mậu Thầu (en: High Card)
func TestCompareBasicStraightFlushVsHighCard(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Tứ quý (en: Four of a Kind)
// Same level card
func TestCompareFourOfAKindVsFourOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Tứ quý (en: Four of a Kind) vs Cù lũ (en: Full House)
func TestCompareFourOfAKindVsFullHouse(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thùng (en: Flush) vs Thùng (en: Flush)
func TestCompareFlushVsFlushHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_6,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// No contain A card
func TestCompareStraightVsStraightNoACardEqual(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

func TestCompareStraightVsStraightNoACardHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// Contain A card
func TestCompareStraightVsStraightContainACardEqual(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// Contain A card, No card K
func TestCompareStraightVsStraightContainACardNotCardKLower(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_6,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(-1), point)
}

// Sảnh (en: Straight) vs Sảnh (en: Straight)
// Contain A card, contain K card
func TestCompareStraightVsStraightContainACardKCard(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Xám chi/Xám cô (en: Three of a Kind) vs Xám chi/Xám cô (en: Three of a Kind)
func TestCompareThreeOfAKindVsThreeOfAKind(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairDraw(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairHigher1(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairHigher2(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Thú (en: Two Pairs) vs Thú (en: Two Pairs) Draw
func TestCompareTwoPairVsTwoPairHigher3(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Đôi (en: Pair) vs Đôi (en: Pair)
func TestComparePairVsPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

// Đôi (en: Pair) vs Đôi (en: Pair)
func TestComparePairVsPairHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

// Mậu Thầu (en: High Card) vs Mậu Thầu (en: High Card)
func TestCompareHighCardVsHighCard(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(0), point)
}

func TestCompareHighCardVsHighCardHigher(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
	}
	strainghtFlush1 := NewChildHand(entity.NewListCard(cards))
	cards = []*pb.Card{
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_Q,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	}
	strainghtFlush2 := NewChildHand(entity.NewListCard(cards))
	point := strainghtFlush1.CompareHand(strainghtFlush2)
	assert.Equal(t, int(1), point)
}

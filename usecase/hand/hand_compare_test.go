package hand

import (
	"testing"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

func mockHandNatural1() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_4,
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
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
		},
	})
}

func mockHandNatural2() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func mockHandNatural3Flush() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},

			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	})
}

func mockHandNormal() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},

			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_A,
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
		},
	})
}

func TestCompareHand(t *testing.T) {

	h1, _ := mockHandNatural1()
	h2, _ := mockHandNatural2()

	result := CompareHand(nil, h1, h2)
	t.Logf("result %v", result)
}

func TestCompareHand3FlushWithNormal(t *testing.T) {

	h1, _ := mockHandNatural3Flush()
	h2, _ := mockHandNormal()

	result := CompareHand(nil, h1, h2)
	t.Logf("result %v", result)
}

package api

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"testing"
)

func mockHand1() (*Hand, error) {
	return NewHand(pb.ListCard{
		Cards: []*pb.Card{
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
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
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
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func mockHand2() (*Hand, error) {
	return NewHand(pb.ListCard{
		Cards: []*pb.Card{
			// Front
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
			// Middle
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			// Back
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
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
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func TestHand(t *testing.T) {
	t.Logf("Test Hand")

	h1, err := mockHand1()
	if err != nil {
		t.Errorf("invalid hand %v", err)
	}

	for _, card := range h1.GetCards().Cards {
		t.Logf("hand %v", card)
	}

	// test calculate
	h1.calculatePoint()
	t.Logf("caculate front %v", h1.frontHandPoint)
	t.Logf("caculate middle %v", h1.middleHandPoint)
	t.Logf("caculate back %v", h1.backHandPoint)

	// test compare
	h2, err := mockHand2()
	if err != nil {
		t.Errorf("invalid hand %v", err)
	}

	for _, card := range h2.GetCards().Cards {
		t.Logf("hand2 %v", card)
	}

	comp := CompareHand(h1, h2)

	t.Logf("compare result: %v", comp)
}

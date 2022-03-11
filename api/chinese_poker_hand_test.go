package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"testing"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

func mockHand1() (*Hand, error) {
	return NewHand(&pb.ListCard{
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
	return NewHand(&pb.ListCard{
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
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_3,
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

	for _, card := range h1.GetCards() {
		t.Logf("hand %v", card)
	}

	// test calculate
	h1.calculatePoint()
	t.Logf("caculate front %v", h1.frontHand.Point)
	t.Logf("caculate middle %v", h1.middleHand.Point)
	t.Logf("caculate back %v", h1.backHand.Point)

	// test compare
	h2, err := mockHand2()
	if err != nil {
		t.Errorf("invalid hand %v", err)
	}

	for _, card := range h2.GetCards() {
		t.Logf("hand2 %v", card)
	}

	// test calculate
	h2.calculatePoint()
	t.Logf("caculate front %v", h2.frontHand.Point)
	t.Logf("caculate middle %v", h2.middleHand.Point)
	t.Logf("caculate back %v", h2.backHand.Point)

	comp := CompareHand(h1, h2)

	t.Logf("compare result: %v", comp)
}

func TestCheck(t *testing.T) {
	t.Logf("check begin")

	unsortCard := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_8,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	sortedCard := SortCard(entity.NewListCard(unsortCard))
	t.Logf("sorted %v", sortedCard)

	cards := []*pb.Card{
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
	}

	if _, ok := CheckFlush(entity.NewListCard(cards)); !ok {
		t.Errorf("wrong check flush")
	} else {
		t.Logf("check flush ok")
	}

	if _, ok := CheckStraight(entity.NewListCard(cards)); !ok {
		t.Errorf("wrong check straight")
	} else {
		t.Logf("check straight ok")
	}

	if _, ok := CheckStraightFlush(entity.NewListCard(cards)); !ok {
		t.Errorf("wrong check straight flush")
	} else {
		t.Logf("check straight flush ok")
	}

	fourOfAKindCards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_A,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckFourOfAKind(entity.NewListCard(fourOfAKindCards)); !ok {
		t.Errorf("wrong check four of a kind")
	} else {
		t.Logf("check four of a kind ok")
	}

	fullHouseCards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
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
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckFullHouse(entity.NewListCard(fullHouseCards)); !ok {
		t.Errorf("wrong check full house card")
	} else {
		t.Logf("check full house ok")
	}
}

func TestThreeOfAKind(t *testing.T) {
	threeOfAKindCards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
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
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckThreeOfAKind(entity.NewListCard(threeOfAKindCards)); !ok {
		t.Errorf("wrong check three of a kind card")
	} else {
		t.Logf("check three of a kind ok")
	}
}

func TestTwoPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckTwoPairs(entity.NewListCard(cards)); !ok {
		t.Errorf("wrong check two pairs")
	} else {
		t.Logf("check two pairs ok")
	}
}

func TestPair(t *testing.T) {
	cards := []*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
	}

	if _, ok := CheckPair(entity.NewListCard(cards)); !ok {
		t.Errorf("wrong check pairs")
	} else {
		t.Logf("check pairs ok")
	}
}

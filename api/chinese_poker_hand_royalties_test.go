package api

import (
	"testing"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

func TestCheckSameSuit(t *testing.T) {
	cards := ListCard{
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
	}

	if _, ok := CheckSameSuit(cards); ok {
		t.Logf("check same suit ok")
	} else {
		t.Logf("check same suit failed")
	}
}

func TestCheckSameColor(t *testing.T) {
	cards := ListCard{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_SPADES,
		},
		{
			Rank: pb.CardRank_RANK_4,
			Suit: pb.CardSuit_SUIT_SPADES,
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
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
	}

	if _, ok := CheckSameColor(cards); ok {
		t.Logf("check same color ok")
	} else {
		t.Logf("check same color failed")
	}
}

func TestCheckDragon(t *testing.T) {
	cards := ListCard{
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
	}

	if _, ok := CheckDragon(cards); ok {
		t.Logf("check dragon ok")
	} else {
		t.Logf("check dragon failed")
	}
}

func TestCheckSixPairs(t *testing.T) {
	cards := ListCard{
		{
			Rank: pb.CardRank_RANK_2,
			Suit: pb.CardSuit_SUIT_CLUBS,
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
			Rank: pb.CardRank_RANK_5,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_7,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_9,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},
		{
			Rank: pb.CardRank_RANK_K,
			Suit: pb.CardSuit_SUIT_SPADES,
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

	if _, ok := CheckSixPairs(cards); ok {
		t.Logf("check six pairs ok")
	} else {
		t.Logf("check six pairs failed")
	}
}

package chinese_poker

// import (
// 	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
// 	"testing"

// 	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
// )

// func TestCheckCleanDragon(t *testing.T) {
// 	cards := []*pb.Card{
// 		{
// 			Rank: pb.CardRank_RANK_2,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_3,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_4,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_5,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_6,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_7,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_8,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_9,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_10,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_J,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_Q,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_K,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_A,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 	}

// 	if _, ok := CheckCleanDragon(entity.NewListCard(cards)); ok {
// 		t.Logf("check clean dragon ok")
// 	} else {
// 		t.Logf("check clean dragon failed")
// 	}
// }

// func TestCheckFullColor(t *testing.T) {
// 	cards := []*pb.Card{
// 		{
// 			Rank: pb.CardRank_RANK_2,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_3,
// 			Suit: pb.CardSuit_SUIT_SPADES,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_4,
// 			Suit: pb.CardSuit_SUIT_SPADES,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_5,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_6,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_7,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_8,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_9,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_10,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_J,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_Q,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_K,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_A,
// 			Suit: pb.CardSuit_SUIT_DIAMONDS,
// 		},
// 	}

// 	if _, ok := CheckFullColor(entity.NewListCard(cards)); ok {
// 		t.Logf("check full color ok")
// 	} else {
// 		t.Logf("check full color failed")
// 	}
// }

// func TestCheckDragon(t *testing.T) {
// 	cards := []*pb.Card{
// 		{
// 			Rank: pb.CardRank_RANK_2,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_3,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_4,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_5,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_6,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_7,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_8,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_9,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_10,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_J,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_Q,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_K,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_A,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 	}

// 	if _, ok := CheckDragon(entity.NewListCard(cards)); ok {
// 		t.Logf("check dragon ok")
// 	} else {
// 		t.Logf("check dragon failed")
// 	}
// }

// func TestCheckSixPairs(t *testing.T) {
// 	cards := []*pb.Card{
// 		{
// 			Rank: pb.CardRank_RANK_2,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_2,
// 			Suit: pb.CardSuit_SUIT_DIAMONDS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_5,
// 			Suit: pb.CardSuit_SUIT_DIAMONDS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_5,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_7,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_7,
// 			Suit: pb.CardSuit_SUIT_DIAMONDS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_9,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_9,
// 			Suit: pb.CardSuit_SUIT_DIAMONDS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_J,
// 			Suit: pb.CardSuit_SUIT_HEARTS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_J,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_K,
// 			Suit: pb.CardSuit_SUIT_SPADES,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_K,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 		{
// 			Rank: pb.CardRank_RANK_A,
// 			Suit: pb.CardSuit_SUIT_CLUBS,
// 		},
// 	}

// 	if _, ok := CheckSixPairs(entity.NewListCard(cards)); ok {
// 		t.Logf("check six pairs ok")
// 	} else {
// 		t.Logf("check six pairs failed")
// 	}
// }

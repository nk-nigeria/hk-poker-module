package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"testing"
)

func TestNewBinListCards(t *testing.T) {
	t.Log("test bin list card")

	cards := NewListCard([]*pb.Card{
		{
			Rank: pb.CardRank_RANK_10,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},

		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},

		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_DIAMONDS,
		},

		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_SPADES,
		},

		{
			Rank: pb.CardRank_RANK_3,
			Suit: pb.CardSuit_SUIT_CLUBS,
		},

		{
			Rank: pb.CardRank_RANK_J,
			Suit: pb.CardSuit_SUIT_HEARTS,
		},
	})

	binCardList := NewBinListCards(cards)
	found, cards := binCardList.GetChain(kCombineFour)
	t.Logf("cards: %s, is four of a kind %v of %v", binCardList, found, cards)
}

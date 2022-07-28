package mockcodegame

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

var MapMockCodeListCard = make(map[int32][]*pb.ListCard)

func init() {
	MapMockCodeListCard[1] = append(MapMockCodeListCard[1], &pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},

			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	},
	)
	MapMockCodeListCard[1] = append(MapMockCodeListCard[1], &pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_2,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	},
	)
}

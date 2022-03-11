package entity

import (
	"fmt"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

type Card int16
type ListCard []*Card

const (
	Rank2  = 0x02
	Rank3  = 0x03
	Rank4  = 0x04
	Rank5  = 0x05
	Rank6  = 0x06
	Rank7  = 0x07
	Rank8  = 0x08
	Rank9  = 0x09
	Rank10 = 0x0A
	RankJ  = 0x0B
	RankQ  = 0x0C
	RankK  = 0x0D
	RankA  = 0x0E

	SuitClubs    = 0x10
	SuitSpides   = 0x20
	SuitDiamonds = 0x30
	SuitHearts   = 0x40
)

var mapRanks = map[pb.CardRank]int16{
	pb.CardRank_RANK_2:  Rank2,
	pb.CardRank_RANK_3:  Rank3,
	pb.CardRank_RANK_4:  Rank4,
	pb.CardRank_RANK_5:  Rank5,
	pb.CardRank_RANK_6:  Rank6,
	pb.CardRank_RANK_7:  Rank7,
	pb.CardRank_RANK_8:  Rank8,
	pb.CardRank_RANK_9:  Rank9,
	pb.CardRank_RANK_10: Rank10,
	pb.CardRank_RANK_J:  RankJ,
	pb.CardRank_RANK_Q:  RankQ,
	pb.CardRank_RANK_K:  RankK,
	pb.CardRank_RANK_A:  RankA,
}

var mapSuits = map[pb.CardSuit]int16{
	pb.CardSuit_SUIT_CLUBS:    SuitClubs,
	pb.CardSuit_SUIT_SPADES:   SuitSpides,
	pb.CardSuit_SUIT_DIAMONDS: SuitDiamonds,
	pb.CardSuit_SUIT_HEARTS:   SuitHearts,
}

var mapStringRanks = map[int16]string{
	Rank2:  "Rank2",
	Rank3:  "Rank3",
	Rank4:  "Rank4",
	Rank5:  "Rank5",
	Rank6:  "Rank6",
	Rank7:  "Rank7",
	Rank8:  "Rank8",
	Rank9:  "Rank9",
	Rank10: "Rank10",
	RankJ:  "RankJ",
	RankQ:  "RankQ",
	RankK:  "RankK",
	RankA:  "RankA",
}

var mapStringSuits = map[int16]string{
	SuitClubs:    "Clubs",
	SuitSpides:   "Spides",
	SuitDiamonds: "Diamonds",
	SuitHearts:   "Hearts",
}

func NewCard(rank pb.CardRank, suit pb.CardSuit) Card {
	card := int16(0)
	card |= mapRanks[rank]
	card |= mapSuits[suit]
	return Card(card)
}

func (c Card) String() string {
	return fmt.Sprintf("Rank %s, Suit %s", mapStringRanks[int16(c&0x0F)], mapStringSuits[int16(c&0xF0)])
}

package entity

import (
	"fmt"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

type Card uint8

const (
	Rank2  uint8 = 0x20
	Rank3  uint8 = 0x30
	Rank4  uint8 = 0x40
	Rank5  uint8 = 0x50
	Rank6  uint8 = 0x60
	Rank7  uint8 = 0x70
	Rank8  uint8 = 0x80
	Rank9  uint8 = 0x90
	Rank10 uint8 = 0xA0
	RankJ  uint8 = 0xB0
	RankQ  uint8 = 0xC0
	RankK  uint8 = 0xD0
	RankA  uint8 = 0xE0

	SuitClubs    uint8 = 0x01
	SuitSpides   uint8 = 0x02
	SuitDiamonds uint8 = 0x03
	SuitHearts   uint8 = 0x04
)

var ranks = []uint8{
	Rank2,
	Rank3,
	Rank4,
	Rank5,
	Rank6,
	Rank7,
	Rank8,
	Rank9,
	Rank10,
	RankJ,
	RankQ,
	RankK,
	RankA,
}

var suits = []uint8{
	SuitClubs,
	SuitSpides,
	SuitDiamonds,
	SuitHearts,
}

var mapRanks = map[pb.CardRank]uint8{
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

var mapSuits = map[pb.CardSuit]uint8{
	pb.CardSuit_SUIT_CLUBS:    SuitClubs,
	pb.CardSuit_SUIT_SPADES:   SuitSpides,
	pb.CardSuit_SUIT_DIAMONDS: SuitDiamonds,
	pb.CardSuit_SUIT_HEARTS:   SuitHearts,
}

var mapStringRanks = map[uint8]string{
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

var mapStringSuits = map[uint8]string{
	SuitClubs:    "Clubs",
	SuitSpides:   "Spides",
	SuitDiamonds: "Diamonds",
	SuitHearts:   "Hearts",
}

func NewCardFromPb(rank pb.CardRank, suit pb.CardSuit) Card {
	card := uint8(0)
	card |= mapRanks[rank]
	card |= mapSuits[suit]
	return Card(card)
}

func NewCard(rank uint8, suit uint8) Card {
	card := uint8(0)
	card |= rank
	card |= suit
	return Card(card)
}

func (c Card) String() string {
	return fmt.Sprintf("Rank: %s, Suit: %s", mapStringRanks[c.GetRank()], mapStringSuits[c.GetSuit()])
}

func (c Card) GetRank() uint8 {
	return uint8(c & 0xF0)
}

func (c Card) GetSuit() uint8 {
	return uint8(c & 0x0F)
}

func IsSameCards(cardsA []*pb.Card, cardsB []*pb.Card) bool {
	mapCardsA := make(map[Card]bool)
	for _, c := range cardsA {
		key := NewCardFromPb(c.GetRank(), c.GetSuit())
		mapCardsA[key] = false
	}
	for _, c := range cardsB {
		key := NewCardFromPb(c.GetRank(), c.GetSuit())
		_, exist := mapCardsA[key]
		if !exist {
			return false
		}
		mapCardsA[key] = true
	}
	for _, v := range mapCardsA {
		if !v {
			return false
		}
	}
	return true
}

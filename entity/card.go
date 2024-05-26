package entity

import (
	"fmt"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type Card uint8

const (
	RankStep uint8 = 0x10

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

	SuitNone     uint8 = 0x00
	SuitSpades   uint8 = 0x01
	SuitClubs    uint8 = 0x02
	SuitDiamonds uint8 = 0x03
	SuitHearts   uint8 = 0x04
)

var Ranks = []uint8{
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

var Suits = []uint8{
	SuitClubs,
	SuitSpades,
	SuitDiamonds,
	SuitHearts,
}

const (
	Card2S = Card(Rank2 | SuitSpades)
	Card2C = Card(Rank2 | SuitClubs)
	Card2D = Card(Rank2 | SuitDiamonds)
	Card2H = Card(Rank2 | SuitHearts)

	Card3S = Card(Rank3 | SuitSpades)
	Card3C = Card(Rank3 | SuitClubs)
	Card3D = Card(Rank3 | SuitDiamonds)
	Card3H = Card(Rank3 | SuitHearts)

	Card4S = Card(Rank4 | SuitSpades)
	Card4C = Card(Rank4 | SuitClubs)
	Card4D = Card(Rank4 | SuitDiamonds)
	Card4H = Card(Rank4 | SuitHearts)

	Card5S = Card(Rank5 | SuitSpades)
	Card5C = Card(Rank5 | SuitClubs)
	Card5D = Card(Rank5 | SuitDiamonds)
	Card5H = Card(Rank5 | SuitHearts)

	Card6S = Card(Rank6 | SuitSpades)
	Card6C = Card(Rank6 | SuitClubs)
	Card6D = Card(Rank6 | SuitDiamonds)
	Card6H = Card(Rank6 | SuitHearts)

	Card7S = Card(Rank7 | SuitSpades)
	Card7C = Card(Rank7 | SuitClubs)
	Card7D = Card(Rank7 | SuitDiamonds)
	Card7H = Card(Rank7 | SuitHearts)

	Card8S = Card(Rank8 | SuitSpades)
	Card8C = Card(Rank8 | SuitClubs)
	Card8D = Card(Rank8 | SuitDiamonds)
	Card8H = Card(Rank8 | SuitHearts)

	Card9S = Card(Rank9 | SuitSpades)
	Card9C = Card(Rank9 | SuitClubs)
	Card9D = Card(Rank9 | SuitDiamonds)
	Card9H = Card(Rank9 | SuitHearts)

	Card10S = Card(Rank10 | SuitSpades)
	Card10C = Card(Rank10 | SuitClubs)
	Card10D = Card(Rank10 | SuitDiamonds)
	Card10H = Card(Rank10 | SuitHearts)

	CardJS = Card(RankJ | SuitSpades)
	CardJC = Card(RankJ | SuitClubs)
	CardJD = Card(RankJ | SuitDiamonds)
	CardJH = Card(RankJ | SuitHearts)

	CardQS = Card(RankQ | SuitSpades)
	CardQC = Card(RankQ | SuitClubs)
	CardQD = Card(RankQ | SuitDiamonds)
	CardQH = Card(RankQ | SuitHearts)

	CardKS = Card(RankK | SuitSpades)
	CardKC = Card(RankK | SuitClubs)
	CardKD = Card(RankK | SuitDiamonds)
	CardKH = Card(RankK | SuitHearts)

	CardAS = Card(RankA | SuitSpades)
	CardAC = Card(RankA | SuitClubs)
	CardAD = Card(RankA | SuitDiamonds)
	CardAH = Card(RankA | SuitHearts)
)

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
	pb.CardSuit_SUIT_SPADES:   SuitSpades,
	pb.CardSuit_SUIT_CLUBS:    SuitClubs,
	pb.CardSuit_SUIT_DIAMONDS: SuitDiamonds,
	pb.CardSuit_SUIT_HEARTS:   SuitHearts,
}

var mapRanksPb = map[uint8]pb.CardRank{
	Rank2:  pb.CardRank_RANK_2,
	Rank3:  pb.CardRank_RANK_3,
	Rank4:  pb.CardRank_RANK_4,
	Rank5:  pb.CardRank_RANK_5,
	Rank6:  pb.CardRank_RANK_6,
	Rank7:  pb.CardRank_RANK_7,
	Rank8:  pb.CardRank_RANK_8,
	Rank9:  pb.CardRank_RANK_9,
	Rank10: pb.CardRank_RANK_10,
	RankJ:  pb.CardRank_RANK_J,
	RankQ:  pb.CardRank_RANK_Q,
	RankK:  pb.CardRank_RANK_K,
	RankA:  pb.CardRank_RANK_A,
}

var mapSuitsPb = map[uint8]pb.CardSuit{
	SuitSpades:   pb.CardSuit_SUIT_SPADES,
	SuitClubs:    pb.CardSuit_SUIT_CLUBS,
	SuitDiamonds: pb.CardSuit_SUIT_DIAMONDS,
	SuitHearts:   pb.CardSuit_SUIT_HEARTS,
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
	SuitSpades:   "Spades",
	SuitClubs:    "Clubs",
	SuitDiamonds: "Diamonds",
	SuitHearts:   "Hearts",
}

func NewCardFromPb(rank pb.CardRank, suit pb.CardSuit) Card {
	card := uint8(0)
	card |= mapRanks[rank]
	card |= mapSuits[suit]
	return Card(card)
}

func NewCardFromUint(c uint) Card {
	return Card(c)
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

func (c Card) ToPB() *pb.Card {
	pbCard := &pb.Card{}
	pbCard.Rank = mapRanksPb[c.GetRank()]
	pbCard.Suit = mapSuitsPb[c.GetSuit()]
	return pbCard
}

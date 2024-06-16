package entity

import (
	"fmt"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type Card uint8

type Rank uint8

const (
	RankStep Rank = 0x10

	Rank2  Rank = 0x20
	Rank3  Rank = 0x30
	Rank4  Rank = 0x40
	Rank5  Rank = 0x50
	Rank6  Rank = 0x60
	Rank7  Rank = 0x70
	Rank8  Rank = 0x80
	Rank9  Rank = 0x90
	Rank10 Rank = 0xA0
	RankJ  Rank = 0xB0
	RankQ  Rank = 0xC0
	RankK  Rank = 0xD0
	RankA  Rank = 0xE0

	SuitNone     uint8 = 0x00
	SuitSpades   uint8 = 0x01
	SuitClubs    uint8 = 0x02
	SuitDiamonds uint8 = 0x03
	SuitHearts   uint8 = 0x04
)

var Ranks = []Rank{
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
	Card2S = Card(uint8(Rank2) | SuitSpades)
	Card2C = Card(uint8(Rank2) | SuitClubs)
	Card2D = Card(uint8(Rank2) | SuitDiamonds)
	Card2H = Card(uint8(Rank2) | SuitHearts)

	Card3S = Card(uint8(Rank3) | SuitSpades)
	Card3C = Card(uint8(Rank3) | SuitClubs)
	Card3D = Card(uint8(Rank3) | SuitDiamonds)
	Card3H = Card(uint8(Rank3) | SuitHearts)

	Card4S = Card(uint8(Rank4) | SuitSpades)
	Card4C = Card(uint8(Rank4) | SuitClubs)
	Card4D = Card(uint8(Rank4) | SuitDiamonds)
	Card4H = Card(uint8(Rank4) | SuitHearts)

	Card5S = Card(uint8(Rank5) | SuitSpades)
	Card5C = Card(uint8(Rank5) | SuitClubs)
	Card5D = Card(uint8(Rank5) | SuitDiamonds)
	Card5H = Card(uint8(Rank5) | SuitHearts)

	Card6S = Card(uint8(Rank6) | SuitSpades)
	Card6C = Card(uint8(Rank6) | SuitClubs)
	Card6D = Card(uint8(Rank6) | SuitDiamonds)
	Card6H = Card(uint8(Rank6) | SuitHearts)

	Card7S = Card(uint8(Rank7) | SuitSpades)
	Card7C = Card(uint8(Rank7) | SuitClubs)
	Card7D = Card(uint8(Rank7) | SuitDiamonds)
	Card7H = Card(uint8(Rank7) | SuitHearts)

	Card8S = Card(uint8(Rank8) | SuitSpades)
	Card8C = Card(uint8(Rank8) | SuitClubs)
	Card8D = Card(uint8(Rank8) | SuitDiamonds)
	Card8H = Card(uint8(Rank8) | SuitHearts)

	Card9S = Card(uint8(Rank9) | SuitSpades)
	Card9C = Card(uint8(Rank9) | SuitClubs)
	Card9D = Card(uint8(Rank9) | SuitDiamonds)
	Card9H = Card(uint8(Rank9) | SuitHearts)

	Card10S = Card(uint8(Rank10) | SuitSpades)
	Card10C = Card(uint8(Rank10) | SuitClubs)
	Card10D = Card(uint8(Rank10) | SuitDiamonds)
	Card10H = Card(uint8(Rank10) | SuitHearts)

	CardJS = Card(uint8(RankJ) | SuitSpades)
	CardJC = Card(uint8(RankJ) | SuitClubs)
	CardJD = Card(uint8(RankJ) | SuitDiamonds)
	CardJH = Card(uint8(RankJ) | SuitHearts)

	CardQS = Card(uint8(RankQ) | SuitSpades)
	CardQC = Card(uint8(RankQ) | SuitClubs)
	CardQD = Card(uint8(RankQ) | SuitDiamonds)
	CardQH = Card(uint8(RankQ) | SuitHearts)

	CardKS = Card(uint8(RankK) | SuitSpades)
	CardKC = Card(uint8(RankK) | SuitClubs)
	CardKD = Card(uint8(RankK) | SuitDiamonds)
	CardKH = Card(uint8(RankK) | SuitHearts)

	CardAS = Card(uint8(RankA) | SuitSpades)
	CardAC = Card(uint8(RankA) | SuitClubs)
	CardAD = Card(uint8(RankA) | SuitDiamonds)
	CardAH = Card(uint8(RankA) | SuitHearts)
)

var mapRanks = map[pb.CardRank]uint8{
	pb.CardRank_RANK_2:  uint8(Rank2),
	pb.CardRank_RANK_3:  uint8(Rank3),
	pb.CardRank_RANK_4:  uint8(Rank4),
	pb.CardRank_RANK_5:  uint8(Rank5),
	pb.CardRank_RANK_6:  uint8(Rank6),
	pb.CardRank_RANK_7:  uint8(Rank7),
	pb.CardRank_RANK_8:  uint8(Rank8),
	pb.CardRank_RANK_9:  uint8(Rank9),
	pb.CardRank_RANK_10: uint8(Rank10),
	pb.CardRank_RANK_J:  uint8(RankJ),
	pb.CardRank_RANK_Q:  uint8(RankQ),
	pb.CardRank_RANK_K:  uint8(RankK),
	pb.CardRank_RANK_A:  uint8(RankA),
}

var mapSuits = map[pb.CardSuit]uint8{
	pb.CardSuit_SUIT_SPADES:   SuitSpades,
	pb.CardSuit_SUIT_CLUBS:    SuitClubs,
	pb.CardSuit_SUIT_DIAMONDS: SuitDiamonds,
	pb.CardSuit_SUIT_HEARTS:   SuitHearts,
}

var mapRanksPb = map[Rank]pb.CardRank{
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

var mapStringRanks = map[Rank]string{
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

func (c Card) GetRank() Rank {
	return Rank(c & 0xF0)
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

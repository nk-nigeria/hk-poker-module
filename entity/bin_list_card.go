package entity

import (
	"fmt"
	"github.com/bits-and-blooms/bitset"
)

var (
	kRed   = uint8(0)
	kBlack = uint8(1)
)

var (
	BitSetRankMap map[uint8]*bitset.BitSet
	BitSetSuitMap map[uint8]*bitset.BitSet
	BitSetColor   map[uint8]*bitset.BitSet
)

var (
	CombinePair      = 1
	CombineThree     = 2
	CombineFour      = 3
	CombineStraight  = 4
	CombineFullHouse = 5
	CombineFlush     = 6
	CombineFullColor = 7
)

func init() {
	BitSetRankMap = make(map[uint8]*bitset.BitSet)
	for _, rank := range ranks {
		BitSet := bitset.New(4)
		for _, suit := range suits {
			BitSet.Set(uint(NewCard(rank, suit)))
		}

		BitSetRankMap[rank] = BitSet
	}

	BitSetSuitMap = make(map[uint8]*bitset.BitSet)
	for _, suit := range suits {
		BitSet := bitset.New(16)
		for _, rank := range ranks {
			BitSet.Set(uint(NewCard(rank, suit)))
		}
		BitSetSuitMap[suit] = BitSet
	}

	BitSetColor = make(map[uint8]*bitset.BitSet)
	BitSetColor[kRed] = BitSetSuitMap[SuitHearts].Union(BitSetSuitMap[SuitDiamonds])
	BitSetColor[kBlack] = BitSetSuitMap[SuitSpades].Union(BitSetSuitMap[SuitClubs])
}

type BinListCard struct {
	b *bitset.BitSet
}

func NewBinListCards(cards ListCard) *BinListCard {
	b := bitset.New(MaxCard)
	for _, card := range cards {
		b.Set(uint(card))
	}
	return &BinListCard{
		b: b,
	}
}

func (b BinListCard) String() string {
	var str = "[\n"

	for i, e := b.b.NextSet(0); e; i, e = b.b.NextSet(i + 1) {
		str += fmt.Sprintf("%d\n", i)
	}
	str += "]"

	return str
}

func (b BinListCard) ToList() ListCard {
	return BitSetToListCard(b.b)
}

func (b BinListCard) GetChain(comb int) (uint, ListCard) {
	switch comb {
	case CombineFour:
		return b.lookupFour()
	case CombineThree:
		return b.lookupThree()
	case CombinePair:
		return b.lookupTwo()
	case CombineStraight:
		return b.lookupStraight()
	case CombineFullHouse:
		return b.lookupFullHouse()
	case CombineFlush:
		return b.lookupFlush()
	case CombineFullColor:
		return b.lookupFullColor()
	}

	return 0, nil
}

func BitSetToListCard(b *bitset.BitSet) ListCard {
	cards := ListCard{}
	for i, e := b.NextSet(0); e; i, e = b.NextSet(i + 1) {
		cards = append(cards, NewCardFromUint(i))
	}
	return cards
}

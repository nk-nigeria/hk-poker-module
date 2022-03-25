package entity

import (
	"fmt"
	"github.com/bits-and-blooms/bitset"
)

var (
	BitSetRankMap map[uint8]*bitset.BitSet
)

var (
	kCombinePair  = 1
	kCombineThree = 2
	kCombineFour  = 3
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

func (b BinListCard) GetChain(comb int) (bool, ListCard) {
	switch comb {
	case kCombineFour:
		for _, rank := range ranks {
			if b.b.IntersectionCardinality(BitSetRankMap[rank]) == 4 {
				//return true, NewCard(rank, 0)
			}
		}
	case kCombineThree:
		for _, rank := range ranks {
			if b.b.IntersectionCardinality(BitSetRankMap[rank]) == 3 {
				//return true, NewCard(rank, 0)
			}
		}
	case kCombinePair:
		for _, rank := range ranks {
			if b.b.IntersectionCardinality(BitSetRankMap[rank]) == 2 {
				//return true, NewCard(rank, 0)
			}
			if b.b.IntersectionCardinality(BitSetRankMap[rank]) == 4 {
				//return true, NewCard(rank, 0)
			}
		}
	}

	return false, ListCard{}
}

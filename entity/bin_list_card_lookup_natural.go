package entity

import "github.com/bits-and-blooms/bitset"

func (b BinListCard) lookupFullColor() (uint, ListCard) {
	var black *bitset.BitSet
	var red *bitset.BitSet

	red = b.b.Intersection(BitSetColor[kRed])
	black = b.b.Intersection(BitSetColor[kBlack])

	if red.Count() >= 12 {
		remain := b.b.Difference(red)
		return 1, createResult(remain.Count()+red.Count(), remain, red)
	}

	if black.Count() >= 12 {
		remain := b.b.Difference(black)
		return 1, createResult(remain.Count()+black.Count(), remain, black)
	}

	return 0, nil
}

package entity

import "github.com/bits-and-blooms/bitset"

func (b BinListCard) lookupFullColor() (uint, ListCard) {
	var black *bitset.BitSet
	var red *bitset.BitSet

	red = b.b.Intersection(BitSetColor[kRed])
	black = b.b.Intersection(BitSetColor[kBlack])

	if red.Count() >= 12 {
		result := ListCard{}
		remain := b.b.Difference(red)
		result = append(result, BitSetToListCard(remain)...)
		result = append(result, BitSetToListCard(red)...)
		return 1, result
	}

	if black.Count() >= 12 {
		result := ListCard{}
		remain := b.b.Difference(black)
		result = append(result, BitSetToListCard(remain)...)
		result = append(result, BitSetToListCard(black)...)
		return 1, result
	}

	return 0, nil
}

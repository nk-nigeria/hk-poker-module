package entity

import "github.com/bits-and-blooms/bitset"

func (b BinListCard) lookupFour() (uint, ListCard) {
	for _, rank := range ranks {
		intersec := b.b.Intersection(BitSetRankMap[rank])
		if intersec.Count() == 4 {
			remain := b.b.Difference(intersec)
			result := ListCard{}
			result = append(result, BitSetToListCard(remain)...)
			result = append(result, BitSetToListCard(intersec)...)

			return 1, result
		}
	}

	return 0, nil
}

func (b BinListCard) lookupThree() (uint, ListCard) {
	for _, rank := range ranks {
		intersec := b.b.Intersection(BitSetRankMap[rank])
		if intersec.Count() == 3 {
			remain := b.b.Difference(intersec)
			result := ListCard{}
			result = append(result, BitSetToListCard(remain)...)
			result = append(result, BitSetToListCard(intersec)...)

			return 1, result
		}
	}

	return 0, nil
}

func (b BinListCard) lookupTwo() (uint, ListCard) {
	count := uint(0)
	var pairs bitset.BitSet
	for _, rank := range ranks {
		intersec := b.b.Intersection(BitSetRankMap[rank])
		if c := intersec.Count(); c >= 2 {
			if c == 3 {
				if i, e := intersec.NextSet(0); e {
					intersec.Clear(i)
				}
			}
			pairs.InPlaceUnion(intersec)
			count += c / 2
		}
	}

	if count > 0 {
		result := ListCard{}
		remain := b.b.Difference(&pairs)
		result = append(result, BitSetToListCard(remain)...)
		result = append(result, BitSetToListCard(&pairs)...)
		return count, result
	}

	return 0, nil
}

func (b BinListCard) lookupFullHouse() (uint, ListCard) {
	count := uint(0)
	var pair *bitset.BitSet
	var threes *bitset.BitSet
	for _, rank := range ranks {
		intersec := b.b.Intersection(BitSetRankMap[rank])
		if c := intersec.Count(); c >= 2 {
			if c == 3 {
				threes = intersec
			} else {
				pair = intersec
			}
		}
	}

	if pair != nil && threes != nil {
		result := ListCard{}
		result = append(result, BitSetToListCard(pair)...)
		result = append(result, BitSetToListCard(threes)...)
		return count, result
	}

	return 0, nil
}

func (b BinListCard) lookupStraight() (uint, ListCard) {
	var j uint
	for i, e := b.b.NextSet(0); e; {
		j, e = b.b.NextSet(i + 1)
		if e {
			c1 := NewCardFromUint(i)
			c2 := NewCardFromUint(j)
			if c2.GetRank()-c1.GetRank() != RankStep {
				return 0, nil
			}
			i = j
		}
	}

	result := ListCard{}
	result = append(result, BitSetToListCard(b.b)...)
	return 1, result
}

func (b BinListCard) lookupFlush() (uint, ListCard) {
	if i, e := b.b.NextSet(0); e {
		card := NewCardFromUint(i)
		suit := card.GetSuit()

		if BitSetSuitMap[suit].IsSuperSet(b.b) {
			return 1, b.ToList()
		}
	}

	return 0, nil
}

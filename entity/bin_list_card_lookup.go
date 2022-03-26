package entity

import "github.com/bits-and-blooms/bitset"

func createResult(size uint, sets ...*bitset.BitSet) ListCard {
	result := NewListCardWithSize(size)
	for _, set := range sets {
		result = append(result, bitSetToListCard(set)...)
	}

	return result
}

func (b BinListCard) lookupFour() (uint, ListCard) {
	for _, rank := range ranks {
		intersec := b.b.Intersection(BitSetRankMap[rank])
		if intersec.Count() == 4 {
			remain := b.b.Difference(intersec)

			return 1, createResult(remain.Count()+intersec.Count(), remain, intersec)
		}
	}

	return 0, nil
}

func (b BinListCard) lookupThree() (uint, ListCard) {
	for _, rank := range ranks {
		intersec := b.b.Intersection(BitSetRankMap[rank])
		if intersec.Count() == 3 {
			remain := b.b.Difference(intersec)

			return 1, createResult(remain.Count()+intersec.Count(), remain, intersec)
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
		remain := b.b.Difference(&pairs)
		return count, createResult(remain.Count()+pairs.Count(), remain, &pairs)
	}

	return 0, nil
}

func (b BinListCard) lookupFullHouse() (uint, ListCard) {
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
		return 1, createResult(pair.Count()+threes.Count(), pair, threes)
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

	return 1, b.ToList()
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

package chinese_poker_bin_list

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
)

func createResult(size uint, sets ...*bitset.BitSet) entity.ListCard {
	result := entity.NewListCardWithSize(size)
	for _, set := range sets {
		result = append(result, entity.BitSetToListCard(set)...)
	}

	return result
}

func (s service) lookupFour(b *entity.BinListCard) (uint, entity.ListCard) {
	for _, rank := range entity.Ranks {
		intersec := b.GetBitSet().Intersection(BitSetRankMap[rank])
		if intersec.Count() == 4 {
			remain := b.GetBitSet().Difference(intersec)

			return 1, createResult(remain.Count()+intersec.Count(), remain, intersec)
		}
	}

	return 0, nil
}

func (s service) lookupThree(b *entity.BinListCard) (uint, entity.ListCard) {
	for _, rank := range entity.Ranks {
		intersec := b.GetBitSet().Intersection(BitSetRankMap[rank])
		if intersec.Count() == 3 {
			remain := b.GetBitSet().Difference(intersec)

			return 1, createResult(remain.Count()+intersec.Count(), remain, intersec)
		}
	}

	return 0, nil
}

func (s service) lookupTwo(b *entity.BinListCard) (uint, entity.ListCard) {
	count := uint(0)
	var pairs bitset.BitSet
	for _, rank := range entity.Ranks {
		intersec := b.GetBitSet().Intersection(BitSetRankMap[rank])
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
		remain := b.GetBitSet().Difference(&pairs)
		return count, createResult(remain.Count()+pairs.Count(), remain, &pairs)
	}

	return 0, nil
}

func (s service) lookupFullHouse(b *entity.BinListCard) (uint, entity.ListCard) {
	var pair *bitset.BitSet
	var threes *bitset.BitSet
	for _, rank := range entity.Ranks {
		intersec := b.GetBitSet().Intersection(BitSetRankMap[rank])
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

func (s service) lookupStraight(b *entity.BinListCard) (uint, entity.ListCard) {
	var j uint
	for i, e := b.GetBitSet().NextSet(0); e; {
		j, e = b.GetBitSet().NextSet(i + 1)
		if e {
			c1 := entity.NewCardFromUint(i)
			c2 := entity.NewCardFromUint(j)
			if c2.GetRank()-c1.GetRank() != entity.RankStep {
				return 0, nil
			}
			i = j
		}
	}

	return 1, b.ToList()
}

func (s service) lookupFlush(b *entity.BinListCard) (uint, entity.ListCard) {
	if i, e := b.GetBitSet().NextSet(0); e {
		card := entity.NewCardFromUint(i)
		suit := card.GetSuit()

		if BitSetSuitMap[suit].IsSuperSet(b.GetBitSet()) {
			return 1, b.ToList()
		}
	}

	return 0, nil
}

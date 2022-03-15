package entity

import pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"

type ListCard []Card

func NewListCard(list []*pb.Card) ListCard {
	newList := ListCard{}
	for _, card := range list {
		newList = append(newList, NewCard(card.GetRank(), card.GetSuit()))
	}

	return newList
}

// GetMaxPointCard
func (ls ListCard) GetMaxPointCard() uint8 {
	isStraight := true
	isContainCardRanK := false
	prevRankPoint := ls[0].GetRank()
	for i := 1; i < len(ls); i++ {
		card := ls[i]
		if isStraight {
			rankPoint := card.GetRank()
			prevRankPoint--
			if rankPoint != prevRankPoint {
				isStraight = false
			}
		}
		isContainCardRanK = card.GetRank() == RankK
	}
	maxCard := ls[len(ls)-1]
	// Chú ý rằng trong Mậu Binh có thể xếp sảnh (hoặc thùng phá sảnh) con A ghép với 2,3,4,5
	// (tuy nhiên đây là bài sảnh hay thùng phá sảnh nhỏ nhất),
	// còn con A ghép với 10,J,Q,K là lá bài lớn nhất.
	if isStraight && !isContainCardRanK {
		return 1
	}
	return maxCard.GetRank()
}

// CompareHighCard
// -1 lower
// 0 equal
// 1 higher
func (ls ListCard) CompareHighCard(other ListCard) int8 {
	if len(ls) != len(other) {
		return 0
	}
	for i := len(ls) - 1; i >= 0; i-- {
		point1 := ls[i].GetRank()
		point2 := other[i].GetRank()
		if point1 > point2 {
			return 1
		}
		if point2 < point1 {
			return -1
		}
	}
	return 0
}

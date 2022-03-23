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

// deep copy, no shawdow
func (lc ListCard) Clone() ListCard {
	ml := make(ListCard, 0, len(lc))
	ml = append(ml, lc...)
	return ml
}

func (lc ListCard) SplitHand() (ListCard, ListCard, ListCard) {
	return lc[8:], lc[3:8], lc[:3]
}

// GetMaxRankPointCard
func (ls ListCard) GetMaxRankPointCard() uint8 {
	isStraight := true
	isContainCardRankK := false
	isContainCardRankA := false
	prevRankPoint := ls[0].GetRank()
	maxCard := ls[0]
	secondMaxCard := ls[0]
	for i := 1; i < len(ls); i++ {
		card := ls[i]
		if isStraight {
			rankPoint := card.GetRank()
			prevRankPoint--
			if rankPoint != prevRankPoint {
				isStraight = false
			}
		}
		isContainCardRankK = card.GetRank() == RankK
		isContainCardRankA = card.GetRank() == RankA
		if maxCard.GetRank() < card.GetRank() {
			secondMaxCard = maxCard
			maxCard = card
		}
	}
	// maxCard := ls[len(ls)-1]
	// Chú ý rằng trong Mậu Binh có thể xếp sảnh (hoặc thùng phá sảnh) con A ghép với 2,3,4,5
	// (tuy nhiên đây là bài sảnh hay thùng phá sảnh nhỏ nhất),
	// còn con A ghép với 10,J,Q,K là lá bài lớn nhất.
	if isStraight && isContainCardRankA {
		if !isContainCardRankK {
			return secondMaxCard.GetRank()
		}
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

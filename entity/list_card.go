package entity

import (
	pb "github.com/nk-nigeria/cgp-common/proto"
)

type ListCard []Card

func (lc ListCard) Clone() ListCard {
	newCard := make(ListCard, len(lc))
	copy(newCard, lc)
	return newCard
}

func NewListCard(list []*pb.Card) ListCard {
	newList := ListCard{}
	for _, card := range list {
		newList = append(newList, NewCardFromPb(card.GetRank(), card.GetSuit()))
	}

	return newList
}

func NewListCardWithSize(size uint) ListCard {
	l := make([]Card, 0, size)
	return l
}

package chinese_poker

import (
	"fmt"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
)

type ChildHand struct {
	Cards entity.ListCard
	Point *HandPoint
}

func (ch ChildHand) String() string {
	return fmt.Sprintf("Cards: %v, Point: %v", ch.Cards, ch.Point)
}

func (ch *ChildHand) calculatePoint() {
	if ch.Point != nil {
		return
	}
	ch.Point = CalculatePoint(ch.Cards)
}

func NewChildHand(cards entity.ListCard) *ChildHand {
	child := ChildHand{
		Cards: cards[:],
	}

	return &child
}

func (h1 *ChildHand) CompareHand(h2 *ChildHand) int {
	h1.calculatePoint()
	h2.calculatePoint()

	resultPoint := 0

	return resultPoint
}

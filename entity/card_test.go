package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"testing"
)

func TestCard_NewCard(t *testing.T) {
	t.Logf("test card")

	card := NewCard(pb.CardRank_RANK_2, pb.CardSuit_SUIT_SPADES)

	t.Logf("%s", card)
}

package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"testing"
)

func TestCard_NewCard(t *testing.T) {
	t.Logf("test card")

	card := NewCardFromPb(pb.CardRank_RANK_2, pb.CardSuit_SUIT_SPADES)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.CardRank_RANK_3, pb.CardSuit_SUIT_SPADES)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.CardRank_RANK_4, pb.CardSuit_SUIT_SPADES)
	t.Logf("%b", card)

	card = NewCardFromPb(pb.CardRank_RANK_5, pb.CardSuit_SUIT_SPADES)
	t.Logf("%b", card)
}

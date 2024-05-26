package entity

import (
	"testing"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
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

func TestCard_ToPB(t *testing.T) {
	name := "TestCard_ToPB"
	t.Run(name, func(t *testing.T) {
		card := NewCard(Rank10, SuitSpades)
		pbCard := card.ToPB()
		cCard := NewCardFromPb(pbCard.GetRank(), pbCard.GetSuit())
		assert.Equal(t, cCard.GetRank(), card.GetRank())
	})
}

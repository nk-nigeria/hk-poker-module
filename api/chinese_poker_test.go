package api

import (
	"testing"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func TestGame(t *testing.T) {
	t.Logf("Test Game")
	processor := NewProcessor()

	// mock presense
	presense := linkedhashmap.New()
	presense.Put("user1", nil)
	presense.Put("user2", nil)
	presense.Put("user3", nil)

	// mock state
	state := &entity.MatchState{
		Presences: presense,
	}

	var err = processor.NewGame(state)
	if err != nil {
		t.Errorf("new game error %v", err)
	}

	t.Logf("new game success")
	processor.Deal(state)

	// check dealt cards
	for u, cards := range state.Cards {
		t.Logf("card %v, %v", u, cards)
	}

	card1 := state.Cards["user1"]
	card2 := state.Cards["user2"]
	card3 := state.Cards["user3"]

	cardOrganize1 := entity.Shuffle(card1)
	cardOrganize2 := entity.Shuffle(card2)
	cardOrganize3 := entity.Shuffle(card3)

	processor.Organize(nil, state, "user1", cardOrganize1)
	processor.Organize(nil, state, "user2", cardOrganize2)
	processor.Organize(nil, state, "user3", cardOrganize3)

	processor.Finish(nil, state)
	// check dealt cards
	for u, cards := range state.OrganizeCards {
		t.Logf("card organize %v, %v", u, cards)
	}

	for u, cards := range state.Cards {
		t.Logf("card2 %v, %v", u, cards)
	}
}

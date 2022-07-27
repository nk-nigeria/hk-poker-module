package engine

import (
	"reflect"
	"testing"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func TestGame(t *testing.T) {
	t.Logf("Test Game")
	processor := NewChinesePokerEngine()

	// mock presense
	presense := linkedhashmap.New()
	presense.Put("user1", nil)
	presense.Put("user2", nil)
	presense.Put("user3", nil)

	// mock state
	state := &entity.MatchState{
		Presences:        presense,
		PlayingPresences: linkedhashmap.New(),
	}

	//var err = processor.NewGame(state)
	//if err != nil {
	//	t.Errorf("new game error %v", err)
	//}

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

	processor.Organize(state, "user1", cardOrganize1)
	processor.Organize(state, "user2", cardOrganize2)
	processor.Organize(state, "user3", cardOrganize3)

	processor.Finish(state)
	// check dealt cards
	for u, cards := range state.OrganizeCards {
		t.Logf("card organize %v, %v", u, cards)
	}

	for u, cards := range state.Cards {
		t.Logf("card2 %v, %v", u, cards)
	}
}

func TestEngine_Finish(t *testing.T) {
	type fields struct {
		deck *entity.Deck
	}
	type args struct {
		s *entity.MatchState
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *pb.UpdateFinish
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Engine{
				deck: tt.fields.deck,
			}
			if got := c.Finish(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Engine.Finish() = %v, want %v", got, tt.want)
			}
		})
	}
}

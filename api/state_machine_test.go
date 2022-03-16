package api

import (
	"fmt"
	"testing"
)

func TestStateMachine(t *testing.T) {
	t.Logf("test game state machine")

	game := NewGameStateMachine()
	fmt.Printf("State is %v\n", game.MustState())

	game.Fire(triggerFinish)
	fmt.Printf("State is %v\n", game.MustState())
}

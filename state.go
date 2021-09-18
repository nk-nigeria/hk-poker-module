package main

import (
	"github.com/ciaolink-game-platform/cgp-blackjack-module/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"math/rand"
)

type MatchState struct {
	random     *rand.Rand
	label      *MatchLabel
	emptyTicks int

	// Currently connected users, or reserved spaces.
	presences map[string]runtime.Presence
	// Number of users currently in the process of connecting to the match.
	joinsInProgress int

	// True if there's a game currently in progress.
	playing bool
	// Mark assignments to player user IDs.
	cards map[string]api.ListCard
	// Whose turn it currently is.
	turn string
	// Ticks until they must submit their move.
	deadlineRemainingTicks int64
	// The winner of the current game.
	result map[string]int
	//// The winner positions.
	//winnerPositions []int32
	// Ticks until the next game starts, if applicable.
	nextGameRemainingTicks int64
}

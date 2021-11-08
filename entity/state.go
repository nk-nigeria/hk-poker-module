package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-blackjack-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"math/rand"
)

type MatchLabel struct {
	Open int32  `json:"open"`
	Bet  int32  `json:"bet"`
	Code string `json:"code"`
}

type MatchState struct {
	Random     *rand.Rand
	Label      *MatchLabel
	EmptyTicks int

	// Currently connected users, or reserved spaces.
	Presences *linkedhashmap.Map
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress int

	// True if there's a game currently in progress.
	Playing bool
	// Mark assignments to player user IDs.
	Cards map[string]*pb.ListCard
	// Mark assignments to player user IDs.
	OrganizeCards map[string]*pb.ListCard
	// Whose turn it currently is.
	Turn string
	// Ticks until they must submit their move.
	DeadlineRemainingTicks int64
	//// The winner positions.
	//winnerPositions []int32
	// Ticks until the next game starts, if applicable.
	NextGameRemainingTicks int64
}

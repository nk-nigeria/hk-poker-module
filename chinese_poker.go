package main

import (
	"github.com/ciaolink-game-platform/cgp-blackjack-module/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type ChinesePokerGame struct {
}

// NewGame
func NewGame() *ChinesePokerGame {
	return &ChinesePokerGame{}
}

// Deal
func (c *ChinesePokerGame) Deal(dispatcher runtime.MatchDispatcher) {
	return
}

// Organize
func (c *ChinesePokerGame) Organize(dispatcher runtime.MatchDispatcher, presence string, cards api.ListCard) error {
	return nil
}

// FinishGame
func (c *ChinesePokerGame) FinishGame(dispatcher runtime.MatchDispatcher) {

}

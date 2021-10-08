package main

import (
	"github.com/ciaolink-game-platform/cgp-blackjack-module/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

const MAX_PRESENCE_CARD = 13

type ChinesePokerGame struct {
	deck *Deck
}

// NewProcessor
func NewProcessor() *ChinesePokerGame {
	return &ChinesePokerGame{}
}

// NewGame
func (c *ChinesePokerGame) NewGame(s *MatchState) error {
	c.deck = NewDeck()
	c.deck.Shuffle()

	s.cards = make(map[string]*api.ListCard)
	for _, k := range s.presences.Keys() {
		presence := k.(string)
		cards, err := c.deck.Deal(MAX_PRESENCE_CARD)
		if err == nil {
			s.cards[presence] = cards
		} else {
			return err
		}
	}

	s.organizeCards = make(map[string]*api.ListCard)

	return nil
}

// Organize
func (c *ChinesePokerGame) Organize(dispatcher runtime.MatchDispatcher, s *MatchState, presence string, cards *api.ListCard) error {
	s.organizeCards[presence] = cards
	return nil
}

// FinishGame
func (c *ChinesePokerGame) FinishGame(dispatcher runtime.MatchDispatcher, s *MatchState) {

}

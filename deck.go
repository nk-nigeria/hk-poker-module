package main

import (
	"errors"
	"github.com/ciaolink-game-platform/cgp-blackjack-module/api"
	"math/rand"
)

const MAX_CARD = 52

type Deck struct {
	cards *api.ListCard
	dealt int
}

func NewDeck() *Deck {
	ranks := []api.CardRank{
		api.CardRank_RANK_A,
		api.CardRank_RANK_2,
		api.CardRank_RANK_3,
		api.CardRank_RANK_4,
		api.CardRank_RANK_5,
		api.CardRank_RANK_6,
		api.CardRank_RANK_7,
		api.CardRank_RANK_8,
		api.CardRank_RANK_9,
		api.CardRank_RANK_10,
		api.CardRank_RANK_J,
		api.CardRank_RANK_Q,
	}

	suits := []api.CardSuit{
		api.CardSuit_SUIT_CLUBS,
		api.CardSuit_SUIT_DIAMONDS,
		api.CardSuit_SUIT_HEARTS,
		api.CardSuit_SUIT_SPADES,
	}

	cards := &api.ListCard{}
	for _, rank := range ranks {
		for _, suit := range suits {
			cards.Cards = append(cards.Cards, &api.Card{
				Rank: rank,
				Suit: suit,
			})
		}
	}

	return &Deck{
		dealt: 0,
		cards: cards,
	}
}

// Shuffle the deck
func (d *Deck) Shuffle() {
	for i := 1; i < len(d.cards.Cards); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d.cards.Cards[r], d.cards.Cards[i] = d.cards.Cards[i], d.cards.Cards[r]
		}
	}
}

// Deal a specified amount of cards
func (d *Deck) Deal(n int) (*api.ListCard, error) {
	if (MAX_CARD - d.dealt) < n {
		return nil, errors.New("deck.deal.error-not-enough")
	}

	var cards api.ListCard
	for i := 0; i < n; i++ {
		cards.Cards = append(cards.Cards, d.cards.Cards[d.dealt])
		d.dealt++
	}

	return &cards, nil
}

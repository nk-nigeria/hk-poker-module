package entity

import (
	"errors"
	"math/rand"

	pb "github.com/ciaolink-game-platform/cgp-blackjack-module/proto"
)

const MaxCard = 52

type Deck struct {
	Cards *pb.ListCard
	Dealt int
}

func NewDeck() *Deck {
	ranks := []pb.CardRank{
		pb.CardRank_RANK_A,
		pb.CardRank_RANK_2,
		pb.CardRank_RANK_3,
		pb.CardRank_RANK_4,
		pb.CardRank_RANK_5,
		pb.CardRank_RANK_6,
		pb.CardRank_RANK_7,
		pb.CardRank_RANK_8,
		pb.CardRank_RANK_9,
		pb.CardRank_RANK_10,
		pb.CardRank_RANK_J,
		pb.CardRank_RANK_Q,
		pb.CardRank_RANK_K,
	}

	suits := []pb.CardSuit{
		pb.CardSuit_SUIT_CLUBS,
		pb.CardSuit_SUIT_DIAMONDS,
		pb.CardSuit_SUIT_HEARTS,
		pb.CardSuit_SUIT_SPADES,
	}

	cards := &pb.ListCard{}
	for _, rank := range ranks {
		for _, suit := range suits {
			cards.Cards = append(cards.Cards, &pb.Card{
				Rank: rank,
				Suit: suit,
			})
		}
	}

	return &Deck{
		Dealt: 0,
		Cards: cards,
	}
}

// Shuffle the deck
func (d *Deck) Shuffle() {
	for i := 1; i < len(d.Cards.Cards); i++ {
		// Create a random int up to the number of Cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d.Cards.Cards[r], d.Cards.Cards[i] = d.Cards.Cards[i], d.Cards.Cards[r]
		}
	}
}

// Deal a specified amount of Cards
func (d *Deck) Deal(n int) (*pb.ListCard, error) {
	if (MaxCard - d.Dealt) < n {
		return nil, errors.New("deck.deal.error-not-enough")
	}

	var cards pb.ListCard
	for i := 0; i < n; i++ {
		cards.Cards = append(cards.Cards, d.Cards.Cards[d.Dealt])
		d.Dealt++
	}

	return &cards, nil
}

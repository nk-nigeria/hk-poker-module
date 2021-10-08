package main

import "testing"

func TestShuffle(t *testing.T) {
	deck := NewDeck()
	t.Logf("deck total %v", len(deck.cards.GetCards()))
	t.Log("ok")
}

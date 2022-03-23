package chinese_poker

func IsMisSets(hand *Hand) bool {
	if hand.backHand.Compare(hand.middleHand) < 0 {
		return true
	}

	if hand.middleHand.Compare(hand.frontHand) < 0 {
		return true
	}

	return false
}

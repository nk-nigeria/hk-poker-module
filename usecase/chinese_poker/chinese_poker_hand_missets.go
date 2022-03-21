package chinese_poker

func IsMisSets(hand *Hand) bool {
	if hand.backHand.CompareHand(hand.middleHand) < 0 {
		return true
	}

	if hand.middleHand.CompareHand(hand.frontHand) < 0 {
		return true
	}

	return false
}

package hand

import "github.com/nakamaFramework/cgp-chinese-poker-module/pkg/log"

func IsMisSets(hand *Hand) bool {
	if hand.backHand.Compare(hand.middleHand) < 0 {
		log.GetLogger().Warn("missets %s", hand)
		return true
	}

	if hand.middleHand.Compare(hand.frontHand) < 0 {
		log.GetLogger().Warn("missets %s", hand)
		return true
	}

	return false
}

package mock

import (
	log "github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/log"
	"github.com/heroiclabs/nakama-common/runtime"
)

type MockDispatcher struct {
}

func (m MockDispatcher) BroadcastMessage(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	log.GetLogger().Info("broadcast opcode: %v data %v to presences %v", opCode, string(data), presences)
	return nil
}

func (m MockDispatcher) BroadcastMessageDeferred(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	log.GetLogger().Info("broadcast defer opcode: %v data %v to presences %v", opCode, string(data), presences)
	return nil
}

func (m MockDispatcher) MatchKick(presences []runtime.Presence) error {
	log.GetLogger().Info("kick ", presences)
	return nil
}

func (m MockDispatcher) MatchLabelUpdate(label string) error {
	log.GetLogger().Info("label update ", label)
	return nil
}

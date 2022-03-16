package mock

import (
	"github.com/heroiclabs/nakama-common/runtime"
	log "github.com/sirupsen/logrus"
)

type MockDispatcher struct {
}

func (m MockDispatcher) BroadcastMessage(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	log.Info("broadcast opcode", opCode, " data ", string(data), " to presences ", presences)
	return nil
}

func (m MockDispatcher) BroadcastMessageDeferred(opCode int64, data []byte, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	log.Info("broadcast opcode", opCode, " data ", string(data), " to presences ", presences)
	return nil
}

func (m MockDispatcher) MatchKick(presences []runtime.Presence) error {
	log.Info("kick ", presences)
	return nil
}

func (m MockDispatcher) MatchLabelUpdate(label string) error {
	log.Info("label update ", label)
	return nil
}

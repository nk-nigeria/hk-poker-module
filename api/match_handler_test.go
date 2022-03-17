package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/mock"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"testing"
	"time"
)

func TestMatch(t *testing.T) {
	t.Logf("test match")

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}

	m := NewMatchHandler(marshaler, unmarshaler)
	var params = make(map[string]interface{})
	params["bet"] = int32(100)
	params["name"] = "name"
	params["password"] = "password"

	logger := &mock.MockLog{}
	dispatcher := mock.MockDispatcher{}
	s, _, _ := m.MatchInit(nil, logger, nil, nil, params)

	// mock event routine
	var stop = make(chan bool)
	go func() {
		t.Logf("start mock loop")
		for i := 0; i < 50; i++ {
			t.Logf("log %d", i)
			time.Sleep(time.Millisecond * 500)
			m.MatchLoop(nil, logger, nil, nil, dispatcher, 0, s, nil)
		}

		t.Logf("current state %v", m.GetState())

		stop <- true
	}()

	go func() {
		t.Logf("start mock join leave")
		var presences []runtime.Presence
		presences = make([]runtime.Presence, 1)
		presences[0] = &mock.MockPresence{
			UserId: "user1",
		}

		m.MatchJoin(nil, logger, nil, nil, dispatcher, 0, s, presences)

		time.Sleep(time.Second * 2)
		presences[0] = &mock.MockPresence{
			UserId: "user2",
		}
		m.MatchJoin(nil, logger, nil, nil, dispatcher, 0, s, presences)
	}()

	t.Logf("wait for finish")
	<-stop
	t.Logf("wait for finish done")
}

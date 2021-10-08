// Copyright 2020 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"time"

	"github.com/ciaolink-game-platform/cgp-blackjack-module/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	moduleName = "blackjack"

	tickRate = 5

	maxEmptySec = 60 * 5

	delayBetweenGamesSec = 5
	turnTimeFastSec      = 10
	turnTimeNormalSec    = 20
	maxPlayer            = 4
)

var winningPositions = [][]int32{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
}

// Compile-time check to make sure all required functions are implemented.
var _ runtime.Match = &MatchHandler{}

type MatchLabel struct {
	Open int32  `json:"open"`
	Bet  int32  `json:"bet"`
	Code string `json:"code"`
}

type MatchHandler struct {
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions
	processor   *ChinesePokerGame
}

func NewMatchHandler(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) *MatchHandler {
	return &MatchHandler{
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
	}
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("match init: %v", params)
	bet, ok := params["bet"].(int32)
	if !ok {
		logger.Error("invalid match init parameter \"bet\"")
		return nil, 0, ""
	}

	label := &MatchLabel{
		Open: 1,
		Bet:  bet,
		Code: "blackjack",
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.WithField("error", err).Error("match init failed")
		labelJSON = []byte("{}")
	}

	logger.Info("match init label=", string(labelJSON))

	m.processor = NewProcessor()

	return &MatchState{
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
		label:   label,
		playing: false,
		//presences: make(map[string]runtime.Presence, maxPlayer),
		presences: linkedhashmap.New(),
	}, tickRate, string(labelJSON)
}

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*MatchState)
	logger.Info("match join attempt, state=%v, meta=%v", s, metadata)

	// Check if it's a user attempting to rejoin after a disconnect.
	if presence, ok := s.presences.Get(presence.GetUserId()); ok {
		if presence == nil {
			// User rejoining after a disconnect.
			s.joinsInProgress++
			return s, true, ""
		} else {
			// User attempting to join from 2 different devices at the same time.
			return s, false, "already joined"
		}
	}

	// Check if match is full.
	if s.presences.Size()+s.joinsInProgress >= maxPlayer {
		return s, false, "match full"
	}

	// New player attempting to connect.
	s.joinsInProgress++
	return s, true, ""
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*MatchState)
	logger.Info("match join, state=%v, presences=%v", s, presences)

	for _, presence := range presences {
		s.emptyTicks = 0
		s.presences.Put(presence.GetUserId(), presence)
		s.joinsInProgress--

		// Check if we must send a message to this user to update them on the current game state.
		var msg proto.Message
		var currentPresences []string
		for _, p := range s.presences.Keys() {
			currentPresences = append(currentPresences, p.(string))
		}
		msg = &api.UpdatePresence{
			JoinPresence: presence.GetUserId(),
			Presences:    currentPresences,
		}

		// Send a message to the user that just joined, if one is needed based on the logic above.
		if msg != nil {
			buf, err := m.marshaler.Marshal(msg)
			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				_ = dispatcher.BroadcastMessage(int64(api.OpCodeUpdate_OPCODE_UPDATE_PRESENCE), buf, nil, nil, true)
			}
		}
	}

	// Check if match was open to new players, but should now be closed.
	if s.presences.Size() >= 2 && s.label.Open != 0 {
		s.label.Open = 0
		if labelJSON, err := json.Marshal(s.label); err != nil {
			logger.Error("error encoding label: %v", err)
		} else {
			if err := dispatcher.MatchLabelUpdate(string(labelJSON)); err != nil {
				logger.Error("error updating label: %v", err)
			}
		}
	}

	return s
}

func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*MatchState)
	logger.Info("match leave, state=%v, presences=%v", s, presences)

	// Check if we must send a message to this user to update them on the current game state.
	var msg proto.Message
	for _, presence := range presences {
		//s.presences[presence.GetUserId()] = nil
		_, found := s.presences.Get(presence.GetUserId())
		if found {
			s.presences.Remove(presence.GetUserId())

			var currentPresences []string
			for _, p := range s.presences.Keys() {
				currentPresences = append(currentPresences, p.(string))
			}
			msg = &api.UpdatePresence{
				LeavePresence: presence.GetUserId(),
				Presences:     currentPresences,
			}

			// Send a message to the user that just joined, if one is needed based on the logic above.
			if msg != nil {
				buf, err := m.marshaler.Marshal(msg)
				if err != nil {
					logger.Error("error encoding message: %v", err)
				} else {
					_ = dispatcher.BroadcastMessage(int64(api.OpCodeUpdate_OPCODE_UPDATE_PRESENCE), buf, nil, nil, true)
				}
			}
		}
	}

	return s
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*MatchState)
	logger.Info("match loop, state=%v, messages=%v", s, messages)

	if s.presences.Size()+s.joinsInProgress == 0 {
		s.emptyTicks++
		if s.emptyTicks >= maxEmptySec*tickRate {
			// Match has been empty for too long, close it.
			logger.Info("closing idle match")
			return nil
		}
	}

	if m.checkAutoNewGame(logger, dispatcher, s) {
		return s
	}

	// There's a game in progress. Check for input, update match state, and send messages to clients.
	for _, message := range messages {
		switch api.OpCodeRequest(message.GetOpCode()) {
		case api.OpCodeRequest_OPCODE_REQUEST_NEW_GAME:
			m.processNewGame(logger, dispatcher, s)
		case api.OpCodeRequest_OPCODE_REQUEST_ORGANIZE:
			msg := &api.Organize{}
			err := m.unmarshaler.Unmarshal(message.GetData(), msg)
			if err != nil {
				// Client sent bad data.
				_ = dispatcher.BroadcastMessage(int64(api.OpCodeUpdate_OPCODE_UPDATE_REJECTED), nil, []runtime.Presence{message}, nil, true)
				continue
			}

			m.processOrganize(dispatcher, s, message.GetUserId(), msg)
		default:
			// No other opcodes are expected from the client, so automatically treat it as an error.
			_ = dispatcher.BroadcastMessage(int64(api.OpCodeUpdate_OPCODE_UPDATE_REJECTED), nil, []runtime.Presence{message}, nil, true)
		}
	}

	// Keep track of the time remaining for the player to submit their move. Idle players forfeit.
	m.checkFinishGame(logger, dispatcher, s)

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	return state
}

func (m *MatchHandler) checkAutoNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *MatchState) bool {
	// If there's no game in progress check if we can (and should) start one!
	if !s.playing {
		// Between games any disconnected users are purged, there's no in-progress game for them to return to anyway.
		//for userID, presence := range s.presences. {
		//	if presence == nil {
		//		delete(s.presences, userID)
		//	}
		//}
		// Check if we need to update the label so the match now advertises itself as open to join.
		if s.presences.Size() < 2 && s.label.Open != 1 {
			s.label.Open = 1
			if labelJSON, err := json.Marshal(s.label); err != nil {
				logger.Error("error encoding label: %v", err)
			} else {
				if err := dispatcher.MatchLabelUpdate(string(labelJSON)); err != nil {
					logger.Error("error updating label: %v", err)
				}
			}
		}

		// Check if we have enough players to start a game.
		if s.presences.Size() < 2 {
			return false
		}

		// Check if enough time has passed since the last game.
		if s.nextGameRemainingTicks > 0 {
			s.nextGameRemainingTicks--
			return false
		}

		// We can start a game! Set up the game state and assign the marks to each player.
		s.playing = true
		m.processNewGame(logger, dispatcher, s)
		return true
	}

	return false
}

// Call when client request or timeout
func (m *MatchHandler) processNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *MatchState) {
	err := m.processor.NewGame(s)
	if err == nil {
		for k, v := range s.cards {
			buf, err := m.marshaler.Marshal(&api.UpdateDeal{
				PresenceCard: &api.PresenceCards{
					Presence: k,
					Cards:    v.Cards,
				},
			})

			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				presence, found := s.presences.Get(k)
				if found {
					_ = dispatcher.BroadcastMessage(int64(api.OpCodeUpdate_OPCODE_UPDATE_DEAL), buf, []runtime.Presence{presence.(runtime.Presence)}, nil, true)
				}
			}
		}
	}
}

// Call when client request organize
func (m *MatchHandler) processOrganize(dispatcher runtime.MatchDispatcher, s *MatchState, presence string, msg *api.Organize) {
	err := m.processor.Organize(dispatcher, s, presence, msg.Cards)
	if err == nil {

	} else {

	}
}

// Check should finish game due to enough organize or timeout
func (m *MatchHandler) checkFinishGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *MatchState) {
	m.processor.FinishGame(dispatcher, s)
}

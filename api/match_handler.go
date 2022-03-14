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

package api

import (
	"context"
	"database/sql"
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
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

	name, ok := params["name"].(string)
	if !ok {
		logger.Error("invalid match init parameter \"name\"")
		return nil, 0, ""
	}

	password, ok := params["password"].(string)
	if !ok {
		logger.Error("invalid match init parameter \"password\"")
		return nil, 0, ""
	}

	label := &entity.MatchLabel{
		Open:     1,
		Bet:      bet,
		Code:     entity.ModuleName,
		Name:     name,
		Password: password,
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.WithField("error", err).Error("match init failed")
		labelJSON = []byte("{}")
	}

	logger.Info("match init label=", string(labelJSON))

	m.processor = NewProcessor()

	matchState := entity.NewMathState(label)

	matchState.SetGameState(pb.GameState_GameStateLobby, logger)
	return &matchState, entity.TickRate, string(labelJSON)
}

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*entity.MatchState)
	logger.Info("match join attempt, state=%v, meta=%v", s, metadata)

	// Check if it's a user attempting to rejoin after a disconnect.
	if presence, ok := s.Presences.Get(presence.GetUserId()); ok {
		if presence == nil {
			// User rejoining after a disconnect.
			s.JoinsInProgress++
			return s, true, ""
		} else {
			// User attempting to join from 2 different devices at the same time.
			return s, false, "already joined"
		}
	}

	// Check if match is full.
	if s.Presences.Size()+s.JoinsInProgress >= maxPlayer {
		return s, false, "match full"
	}

	// New player attempting to connect.
	s.JoinsInProgress++
	return s, true, ""
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*entity.MatchState)
	logger.Info("match join, state=%v, presences=%v", s, presences)

	s = s.ProcessEvent(entity.MatchJoin, logger, presences)
	for _, presence := range presences {
		// s.EmptyTicks = 0
		// s.Presences.Put(presence.GetUserId(), presence)
		// s.JoinsInProgress--

		// Check if we must send a message to this user to update them on the current game state.
		var msg proto.Message
		var currentPresences []string
		for _, p := range s.Presences.Keys() {
			currentPresences = append(currentPresences, p.(string))
		}
		msg = &pb.UpdatePresence{
			JoinPresence: presence.GetUserId(),
			Presences:    currentPresences,
		}

		// Send a message to the user that just joined, if one is needed based on the logic above.
		if msg != nil {
			buf, err := m.marshaler.Marshal(msg)
			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_PRESENCE), buf, nil, nil, true)
			}
		}
	}

	// Check if matchwas open to new players, but should now be closed.
	if s.Presences.Size() >= 2 && s.Label.LastOpenValueNoti == 1 {
		s.Label.Open = 0
		s.Label.LastOpenValueNoti = 0
		if labelJSON, err := json.Marshal(s.Label); err != nil {
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
	s := state.(*entity.MatchState)
	logger.Info("match leave, state=%v, presences=%v", s, presences)

	s = s.ProcessEvent(entity.MatchLeave, logger, presences)

	// Check if we must send a message to this user to update them on the current game state.
	var msg proto.Message
	for _, presence := range presences {
		//s.presences[presence.GetUserId()] = nil
		_, found := s.Presences.Get(presence.GetUserId())
		if found {
			var currentPresences []string
			for _, p := range s.Presences.Keys() {
				currentPresences = append(currentPresences, p.(string))
			}
			msg = &pb.UpdatePresence{
				LeavePresence: presence.GetUserId(),
				Presences:     currentPresences,
			}

			// Send a message to the user that just joined, if one is needed based on the logic above.
			if msg != nil {
				buf, err := m.marshaler.Marshal(msg)
				if err != nil {
					logger.Error("error encoding message: %v", err)
				} else {
					_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_PRESENCE), buf, nil, nil, true)
				}
			}
		}
	}

	return s
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*entity.MatchState)
	// logger.Info("match loop, state=%v, messages=%v, game state: %s", s, messages, s.GetGameState().String())

	s = s.ProcessEvent(entity.MathLoop, logger, nil)
	if s.GetGameState() == pb.GameState_GameStateLobby && s.EmptyTicks > entity.MaxEmptySec {
		logger.Info("closing idle match id")
		return nil
	}
	m.checkAndSendUpdateGameState(logger, s, dispatcher)

	if s.GetGameState() == pb.GameState_GameStateFinish {
		return s
	}

	if s.GetGameState() == pb.GameState_GameStateCountdown {
	}

	if s.GetGameState() == pb.GameState_GameStateReward {
	}

	// only accept command from client when
	// game in state run
	if s.GetGameState() != pb.GameState_GameStateRun {
		return s
	}

	if !s.Playing {
		m.processNewGame(logger, dispatcher, s)
		s.Playing = true
	}

	// send time remain before end game
	if s.CountDown.IsUpdate {

		if s.CountDown.Sec == 0 {
			logger.Info("Send notification all card of all user ")
			pbGameState := pb.UpdateGameState{
				State:     s.GetGameState(),
				CountDown: s.CountDown.Sec,
			}
			pbGameState.PresenceCards = make([]*pb.PresenceCards, len(s.Cards))
			for k, v := range s.Cards {
				presenceCards := pb.PresenceCards{
					Presence: k,
					Cards:    v.GetCards(),
				}
				pbGameState.State = pb.GameState_GameStateReward
				pbGameState.PresenceCards = append(pbGameState.PresenceCards, &presenceCards)
				m.broadcastMessage(
					logger, dispatcher,
					int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
					&pbGameState, nil, nil, true)
			}

			// update finish
			updateFinish := m.processor.Finish(dispatcher, s)
			m.broadcastMessage(
				logger, dispatcher,
				int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH),
				updateFinish, nil, nil, true)
		}
	}
	// There's a game in progress. Check for input, update match state, and send messages to clients.
	for _, message := range messages {
		switch pb.OpCodeRequest(message.GetOpCode()) {
		case pb.OpCodeRequest_OPCODE_REQUEST_NEW_GAME:
			m.processNewGame(logger, dispatcher, s)
		case pb.OpCodeRequest_OPCODE_REQUEST_ORGANIZE:
			msg := &pb.Organize{}
			err := m.unmarshaler.Unmarshal(message.GetData(), msg)
			if err != nil {
				// Client sent bad data.
				_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_REJECTED), nil, []runtime.Presence{message}, nil, true)
				continue
			}

			m.processOrganize(dispatcher, s, message.GetUserId(), msg)
		case pb.OpCodeRequest_OPCODE_REQUEST_LEAVE_GAME:
			m.checkLeaveGame(logger, dispatcher, s)
		case pb.OpCodeRequest_OPCODE_REQUEST_COMBINE_CARDS:
			{
				m.combineCard(logger, dispatcher, s, message)
			}
		case pb.OpCodeRequest_OPCODE_REQUEST_SHOW_CARDS:
			{
				m.showCard(logger, dispatcher, s, message)
			}
		case pb.OpCodeRequest_OPCODE_REQUEST_DECLARE_CARDS:
			{
				m.declareCard(logger, dispatcher, s, message)
			}
		default:
			// No other opcodes are expected from the client, so automatically treat it as an error.
			_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_REJECTED), nil, []runtime.Presence{message}, nil, true)
		}
	}

	// // Keep track of the time remaining for the player to submit their move. Idle players forfeit.
	// m.checkFinishGame(logger, dispatcher, s)

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	return state
}

func (m *MatchHandler) checkAutoNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) bool {
	// If there's no game in progress check if we can (and should) start one!
	if !s.Playing {
		// Between games any disconnected users are purged, there's no in-progress game for them to return to anyway.
		//for userID, presence := range s.presences. {
		//	if presence == nil {
		//		delete(s.presences, userID)
		//	}
		//}
		// Check if we need to update the label so the match now advertises itself as open to join.
		if s.Presences.Size() < 2 && s.Label.LastOpenValueNoti == 0 {
			s.Label.Open = 1
			s.Label.LastOpenValueNoti = 1
			if labelJSON, err := json.Marshal(s.Label); err != nil {
				logger.Error("error encoding label: %v", err)
			} else {
				if err := dispatcher.MatchLabelUpdate(string(labelJSON)); err != nil {
					logger.Error("error updating label: %v", err)
				}
			}
		}

		// Check if we have enough players to start a game.
		if s.Presences.Size() < 2 {
			return false
		}

		// Check if enough time has passed since the last game.
		if s.NextGameRemainingTicks > 0 {
			s.NextGameRemainingTicks--
			return false
		}

		// We can start a game! Set up the game state and assign the marks to each player.
		s.Playing = true
		m.processNewGame(logger, dispatcher, s)
		return true
	}

	return false
}

// Call when client request or timeout
func (m *MatchHandler) processNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	err := m.processor.NewGame(s)
	m.processor.Deal(s)
	if err == nil {
		for k, v := range s.Cards {
			buf, err := m.marshaler.Marshal(&pb.UpdateDeal{
				PresenceCard: &pb.PresenceCards{
					Presence: k,
					Cards:    v.Cards,
				},
			})

			if err != nil {
				logger.Error("error encoding message: %v", err)
			} else {
				presence, found := s.Presences.Get(k)
				if found {
					_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), buf, []runtime.Presence{presence.(runtime.Presence)}, nil, true)
				}
			}
		}
	}
}

// Call when client request organize
func (m *MatchHandler) processOrganize(dispatcher runtime.MatchDispatcher, s *entity.MatchState, presence string, msg *pb.Organize) {
	err := m.processor.Organize(dispatcher, s, presence, msg.Cards)
	if err == nil {

	} else {

	}
}

// Check should finish game due to enough organize or timeout
func (m *MatchHandler) checkFinishGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	m.processor.Finish(dispatcher, s)
}

func (m *MatchHandler) checkLeaveGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
}

func (m *MatchHandler) broadcastMessage(logger runtime.Logger, dispatcher runtime.MatchDispatcher, opCode int64, data proto.Message, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	dataJson, err := m.marshaler.Marshal(data)
	if err != nil {
		logger.Error("Error when marshaler data for broadcastMessage")
		return err
	}
	err = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), dataJson, nil, nil, true)
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s", string(dataJson))
		return err
	}
	return nil
}

func (m *MatchHandler) combineCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %d request combineCard", message.GetUserId())
	msg := pb.UpdateGameState{
		State: s.GetGameState(),
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_COMBINE,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), &msg, nil, nil, true)
}

func (m *MatchHandler) showCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %d request showCard", message.GetUserId())

	msg := pb.UpdateGameState{
		State: s.GetGameState(),
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_SHOW,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), &msg, nil, nil, true)
	m.saveCard(logger, s, message)
}

func (m *MatchHandler) declareCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	// msg := &pb.UpdatePresence{
	// 	Presences: []string{message.GetUserId()},
	// }
	//m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeRequest_OPCODE_REQUEST_SHOW_CARDS), msg, nil, nil, true)
	logger.Info("User %d request declareCard", message.GetUserId())
	m.saveCard(logger, s, message)
}

func (m *MatchHandler) saveCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
	cards := s.Cards[message.GetUserId()]
	organize := &pb.Organize{}
	err := m.unmarshaler.Unmarshal(message.GetData(), organize)
	if err != nil {
		logger.Error("Parse organize cards from client error %s", err.Error())
		return
	}
	cardsByClient := organize.Cards
	// check len card
	if len(cardsByClient.GetCards()) != len(cards.GetCards()) {
		logger.Error("Amount cards from client [%d] diffirrent amount card in server [%d]",
			len(cardsByClient.GetCards()), len(cards.GetCards()))
		return
	}
	// check card send by client is the same card in server
	if !entity.IsSameCards(cards.GetCards(), cardsByClient.GetCards()) {
		logger.Error("cards from client not the same card in server, invalid action",
			len(cardsByClient.GetCards()), len(cards.GetCards()))
		return
	}

	// TODO: recheck
	s.Cards[message.GetUserId()] = cardsByClient
	s.OrganizeCards[message.GetUserId()] = cardsByClient
}

func (m *MatchHandler) addChip(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, amountChip int) {
	changeset := map[string]int64{
		"chips": int64(amountChip), // Add amountChip coins to the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_topup": "topup",
	}

	_, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallet update error.")
	}

}

func (m *MatchHandler) subtractChip(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, amountChip int) {
	changeset := map[string]int64{
		"chips": -int64(amountChip), // Substract amountChip coins to the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_topup": "topup",
	}

	_, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallet update error.")
	}
}

func (m *MatchHandler) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, resultGame *pb.UpdateFinish) {
	walletUpdates := make([]*runtime.WalletUpdate, len(resultGame.Results))
	for _, result := range resultGame.Results {
		amountChip := int64(0)
		amountChip = 200*(result.FrontFactor+result.MiddleFactor+result.BackFactor) +
			(result.FrontBonus + result.MiddleBonus + result.BackBonus)
		changeset := map[string]int64{
			"chips": amountChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_topup": "topup",
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}

	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallets update error.")
	}
}

func (m *MatchHandler) checkAndSendUpdateGameState(logger runtime.Logger, s *entity.MatchState, dispatcher runtime.MatchDispatcher) {
	if s.CountDown.IsUpdate {
		pbGameState := pb.UpdateGameState{
			State:     s.GetGameState(),
			CountDown: s.CountDown.Sec,
		}
		// data, err := m.marshaler.Marshal(&pbGameState)
		logger.Info("Send notification countdown from %s --> %s, %d", s.GetGameState().String(), pb.GameState_GameStateFinish.String(), s.CountDown.Sec)
		err := m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), &pbGameState, nil, nil, true)
		if err == nil {
			s.CountDown.IsUpdate = false
		}
	}
}

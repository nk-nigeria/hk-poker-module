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
	"github.com/qmuntal/stateless"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	maxPlayer = 4
)

// Compile-time check to make sure all required functions are implemented.
var _ runtime.Match = &MatchHandler{}

type MatchHandler struct {
	marshaler    *protojson.MarshalOptions
	unmarshaler  *protojson.UnmarshalOptions
	processor    *ChinesePokerGame
	stateMachine *GameStateMachine
}

func NewMatchHandler(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) *MatchHandler {
	return &MatchHandler{
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
	}
}

func (m *MatchHandler) GetState() stateless.State {
	return m.stateMachine.MustState()
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

	// State machine created with state wait
	m.stateMachine = NewGameStateMachine()

	matchState.SetGameState(pb.GameState_GameStateLobby, logger)
	return &matchState, entity.TickRate, string(labelJSON)
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*entity.MatchState)
	// logger.Info("match loop, state=%v, messages=%v, game state: %s", s, messages, s.GetGameState().String())

	m.stateMachine.FireProcessEvent(logger, dispatcher, messages)

	//s = s.ProcessEvent(entity.MathLoop, logger, nil)
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
	if s.GetGameState() != pb.GameState_GameStatePlay {
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

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	return state
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
	err = dispatcher.BroadcastMessage(opCode, dataJson, nil, nil, true)
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s", string(dataJson))
		return err
	}
	return nil
}

func (m *MatchHandler) combineCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request combineCard", message.GetUserId())
	msg := pb.UpdateGameState{
		State: s.GetGameState(),
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_COMBINE,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
}

func (m *MatchHandler) showCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request showCard", message.GetUserId())

	msg := pb.UpdateGameState{
		State: s.GetGameState(),
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_SHOW,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.saveCard(logger, s, message)
}

func (m *MatchHandler) declareCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
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
		logger.Error("Amount cards from client [%d] different amount card in server [%d]",
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
		if s.CountDown.Sec == 0 {
			logger.Info("Send notification countdown from %s --> %s, %d", s.GetGameState().String(), s.GetNextGameState().String(), s.CountDown.Sec)
		}
		err := m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), &pbGameState, nil, nil, true)
		if err == nil {
			s.CountDown.IsUpdate = false
		}
	}
}

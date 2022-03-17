package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const processorKey = "pd"

type MatchProcessor struct {
	gameEngine   *ChinesePokerGame
	stateMachine *GameStateMachine
	marshaler    *protojson.MarshalOptions
	unmarshaler  *protojson.UnmarshalOptions
}

// Call when client request or timeout
func (m *MatchProcessor) processNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	// clean up game state
	m.gameEngine.NewGame(s)

	if err := m.gameEngine.Deal(s); err == nil {
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
				presence, found := s.PlayingPresences.Get(k)
				if found {
					_ = dispatcher.BroadcastMessage(int64(pb.OpCodeUpdate_OPCODE_UPDATE_DEAL), buf, []runtime.Presence{presence.(runtime.Presence)}, nil, true)
				}
			}
		}
	}
}

func (m *MatchProcessor) processFinishGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	logger.Info("process finish game len cards %v", len(s.Cards))
	// send organize card to all
	pbGameState := pb.UpdateGameState{
		State: pb.GameState_GameStateReward,
	}
	pbGameState.PresenceCards = make([]*pb.PresenceCards, len(s.Cards))
	for k, v := range s.Cards {
		organizeCards := s.OrganizeCards[k]
		if organizeCards == nil {
			logger.Warn("user %s not submit cards use deal cards", k)
			organizeCards = v
			s.OrganizeCards[k] = v
		}
		presenceCards := pb.PresenceCards{
			Presence: k,
			Cards:    organizeCards.GetCards(),
		}
		pbGameState.PresenceCards = append(pbGameState.PresenceCards, &presenceCards)
	}

	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
		&pbGameState, nil, nil, true)

	// update finish
	updateFinish := m.gameEngine.Finish(dispatcher, s)
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH),
		updateFinish, nil, nil, true)

	logger.Info("process finish game done %v", updateFinish)
}

func (m *MatchProcessor) combineCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request combineCard", message.GetUserId())
	msg := pb.UpdateGameState{
		State: m.stateMachine.GetPbState(),
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_COMBINE,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.removeShowCard(logger, s, message)
}

func (m *MatchProcessor) showCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request showCard", message.GetUserId())

	msg := pb.UpdateGameState{
		State: m.stateMachine.GetPbState(),
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_SHOW,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.saveCard(logger, s, message)
}

func (m *MatchProcessor) declareCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %d request declareCard", message.GetUserId())
	// TODO: check royalties
	m.saveCard(logger, s, message)
}

func (m *MatchProcessor) saveCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
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

	m.gameEngine.Organize(s, message.GetUserId(), cardsByClient)
}

func (m *MatchProcessor) removeShowCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
	m.gameEngine.Combine(s, message.GetUserId())
}

func (m *MatchProcessor) broadcastMessage(logger runtime.Logger, dispatcher runtime.MatchDispatcher, opCode int64, data proto.Message, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
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

func (m *MatchProcessor) notifyUpdateGameState(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message) {
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
		updateState, nil, nil, true)
}

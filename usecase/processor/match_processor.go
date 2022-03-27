package processor

import (
	"context"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/chinese_poker"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type processor struct {
	engine      chinese_poker.UseCase
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions
}

func NewMatchProcessor(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions, engine chinese_poker.UseCase) UseCase {
	return &processor{
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
		engine:      engine,
	}
}

// Call when client request or timeout
func (m *processor) ProcessNewGame(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	// clean up game state
	m.engine.NewGame(s)

	if err := m.engine.Deal(s); err == nil {
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

func (m *processor) ProcessFinishGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	logger.Info("process finish game len cards %v", len(s.Cards))
	// send organize card to all
	pbGameState := pb.UpdateGameState{
		State: pb.GameState_GameStateReward,
	}
	pbGameState.PresenceCards = make([]*pb.PresenceCards, 0, len(s.Cards))
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
	updateFinish := m.engine.Finish(s)
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH),
		updateFinish, nil, nil, true)

	m.updateWallet(ctx, nk, logger, dispatcher, s, updateFinish)
	logger.Info("process finish game done %v", updateFinish)
}

func (m *processor) CombineCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request combineCard", message.GetUserId())
	msg := pb.UpdateGameState{
		State: pb.GameState_GameStatePlay,
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_COMBINE,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.removeShowCard(logger, s, message)
}

func (m *processor) ShowCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request showCard", message.GetUserId())

	msg := pb.UpdateGameState{
		State: pb.GameState_GameStatePlay,
		ArrangeCard: &pb.ArrangeCard{
			Presence:  message.GetUserId(),
			CardEvent: pb.CardEvent_SHOW,
		},
	}
	m.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_CARD_STATE), &msg, nil, nil, true)
	m.saveCard(logger, s, message)
}

func (m *processor) DeclareCard(logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, message runtime.MatchData) {
	logger.Info("User %s request declareCard", message.GetUserId())
	// TODO: check royalties
	m.saveCard(logger, s, message)
}

func (m *processor) saveCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
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

	logger.Info("update save card %v, %v", message.GetUserId(), cardsByClient)

	m.engine.Organize(s, message.GetUserId(), cardsByClient)
}

func (m *processor) removeShowCard(logger runtime.Logger, s *entity.MatchState, message runtime.MatchData) {
	m.engine.Combine(s, message.GetUserId())
}

func (m *processor) broadcastMessage(logger runtime.Logger, dispatcher runtime.MatchDispatcher, opCode int64, data proto.Message, presences []runtime.Presence, sender runtime.Presence, reliable bool) error {
	dataJson, err := m.marshaler.Marshal(data)
	if err != nil {
		logger.Error("Error when marshaler data for broadcastMessage")
		return err
	}
	err = dispatcher.BroadcastMessage(opCode, dataJson, nil, nil, true)

	logger.Info("broadcast message opcode %v, to %v, data %v", opCode, presences, string(dataJson))
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s", string(dataJson))
		return err
	}
	return nil
}

func (m *processor) NotifyUpdateGameState(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message) {
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
		updateState, nil, nil, true)
}

func (m *processor) NotifyUpdatePresences(s *entity.MatchState, logger runtime.Logger, dispatcher runtime.MatchDispatcher, updateState proto.Message) {
	m.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_PRESENCE),
		updateState, nil, nil, true)
}

func (m *processor) updateWallet(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, dispatcher runtime.MatchDispatcher, s *entity.MatchState, updateFinish *pb.UpdateFinish) {
	listUserId := make([]string, 0, len(updateFinish.Results))
	for _, uf := range updateFinish.Results {
		listUserId = append(listUserId, uf.UserId)
	}
	wallets, err := m.readWalletUsers(ctx, nk, logger, listUserId...)
	if err != nil {
		return
	}
	mapUserWallet := make(map[string]entity.Wallet)
	for _, w := range wallets {
		mapUserWallet[w.UserId] = w
	}

	balanceResult := pb.BalanceResult{}
	for _, uf := range updateFinish.Results {
		balance := pb.BalanceUpdate{
			UserId:           uf.UserId,
			AmountChipBefore: mapUserWallet[uf.UserId].Chips,
		}

		percentFee := 0.05
		fee := int64(percentFee*float64(uf.ScoreResult.NumHandWin)) * int64(s.Label.Bet)
		balance.AmountChipAdd = uf.ScoreResult.TotalFactor * int64(s.Label.Bet)
		balance.AmountChipCurrent = balance.AmountChipCurrent + balance.AmountChipAdd - fee
		balanceResult.Updates = append(balanceResult.Updates, &balance)
	}
	m.updateChipByResultGameFinish(ctx, logger, nk, &balanceResult)
	m.broadcastMessage(
		logger,
		dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_WALLET),
		&balanceResult,
		nil,
		nil,
		true,
	)
}
func (m *processor) readWalletUsers(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger, userIds ...string) ([]entity.Wallet, error) {
	return entity.ReadWalletUsers(ctx, nk, logger, userIds...)
}

func (m *processor) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, balanceResult *pb.BalanceResult) {
	walletUpdates := make([]*runtime.WalletUpdate, 0, len(balanceResult.Updates))
	for _, result := range balanceResult.Updates {
		amoutChip := result.AmountChipCurrent - result.AmountChipBefore
		changeset := map[string]int64{
			"chips": amoutChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_reward": entity.ModuleName,
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}

	logger.Info("wallet update ctx %v, walletUpdates %v", ctx, walletUpdates)

	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallets update error.")
	}
}

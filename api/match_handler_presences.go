package api

import (
	"context"
	"database/sql"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
)

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
	if s.Presences.Size()+s.JoinsInProgress >= entity.MaxPresences {
		return s, false, "match full"
	}

	// New player attempting to connect.
	s.JoinsInProgress++
	return s, true, ""
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*entity.MatchState)
	logger.Info("match join, state=%v, presences=%v", s, presences)

	s.AddPresence(presences)

	for _, presence := range presences {
		// Check if we must send a message to this user to update them on the current game state.
		// var msg proto.Message
		var currentPresences []string
		for _, p := range s.Presences.Keys() {
			currentPresences = append(currentPresences, p.(string))
		}
		msg := &pb.UpdateTable{
			Vip:          0,
			Bet:          int64(s.Label.Bet),
			JoinPresence: entity.NewPlayer(presence),
		}
		listPlayerP := entity.NewListPlayer(s.GetPPresence())
		listPlayerP.ReadWallet(ctx, nk, logger)
		msg.PlayersP = listPlayerP

		listPlayerV := entity.NewListPlayer(s.GetVPresence())
		listPlayerV.ReadWallet(ctx, nk, logger)
		msg.PlayersV = listPlayerV

		// Send a message to the user that just joined, if one is needed based on the logic above.
		m.processor.NotifyUpdateTable(s, logger, dispatcher, msg)
	}

	return s
}

func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*entity.MatchState)

	logger.Info("match leave, state=%v, presences=%v", s, presences)

	if m.machine.IsPlayingState() {
		s.AddLeavePresence(presences)
		return nil
	}

	s.RemovePresence(presences)

	// Check if we must send a message to this user to update them on the current game state.
	for _, presence := range presences {
		_, found := s.Presences.Get(presence.GetUserId())
		if found {
			var currentPresences []string
			for _, p := range s.Presences.Keys() {
				currentPresences = append(currentPresences, p.(string))
			}
			msg := &pb.UpdateTable{
				Vip:           0,
				Bet:           int64(s.Label.Bet),
				LeavePresence: entity.NewPlayer(presence),
			}
			listPlayerP := entity.NewListPlayer(s.GetPPresence())
			listPlayerP.ReadWallet(ctx, nk, logger)
			msg.PlayersP = listPlayerP

			listPlayerv := entity.NewListPlayer(s.GetVPresence())
			listPlayerv.ReadWallet(ctx, nk, logger)
			msg.PlayersV = listPlayerv

			// Send a message to the user that just joined, if one is needed based on the logic above.
			m.processor.NotifyUpdateTable(s, logger, dispatcher, msg)
		}
	}

	return s
}

package api

import (
	"context"
	"database/sql"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/heroiclabs/nakama-common/runtime"
)

func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	s := state.(*entity.MatchState)
	logger.Info("match join attempt, state=%v, meta=%v", s, metadata)

	// check password
	if s.Label.Password != "" {
		logger.Info("match protect with password, check password")
		joinPassword := metadata["password"]
		if joinPassword != s.Label.Password {
			return s, false, "wrong password"
		}
	}

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

	m.processor.ProcessPresencesJoin(ctx, logger, nk, dispatcher, s, presences)

	return s
}

func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	s := state.(*entity.MatchState)

	logger.Info("match leave, state=%v, presences=%v", s, presences)

	if m.machine.IsPlayingState() || m.machine.IsReward() {
		m.processor.ProcessPresencesLeavePending(ctx, logger, nk, dispatcher, s, presences)
		return s
	}

	m.processor.ProcessPresencesLeave(ctx, logger, nk, dispatcher, s, presences)

	return s
}

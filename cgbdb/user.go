package cgbdb

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func UpdateUserPlayingInMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, userIdD string, matchId string) error {
	query := `UPDATE
					users AS u
				SET
					metadata
						= u.metadata
						|| jsonb_build_object('on_playing_in_match','` + matchId + `' )
				WHERE	
					id = $1;`
	_, err := db.ExecContext(ctx, query, userIdD)
	if err != nil {
		logger.WithField("err", err).Error("db.ExecContext match update error.")
	}
	return err
}

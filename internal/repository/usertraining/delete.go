package usertraining

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteUserTraining(ctx context.Context, trainingID int64) error {
	_, err := sq.Delete(tableName).
		Where("training_id = ?", trainingID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("delete user relation for training %d: %w", trainingID, err)
	}

	return nil
}

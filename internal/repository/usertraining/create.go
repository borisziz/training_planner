package usertraining

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) CreateUserTraining(
	ctx context.Context,
	relation model.UserTraining,
) error {
	_, err := sq.Insert(tableName).
		Columns("training_id", "user_id").
		Values(relation.TrainingID, relation.UserID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("create user training relation: %w", err)
	}

	return nil
}

package training

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) Update(ctx context.Context, training *model.Training) error {
	_, err := sq.Update(tableName).
		Set("date", training.Date).
		Set("comment", training.Comment).
		Set("result", training.Result).
		Set("result_files", training.ResultFiles).
		Where("id = ?", training.ID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("update training %d: %w", training.ID, err)
	}

	return nil
}

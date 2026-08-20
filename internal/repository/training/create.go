package training

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) CreateTraining(ctx context.Context, training *model.Training) (int64, error) {
	query, args, err := sq.Insert(tableName).
		Columns("date", "comment", "result", "result_files").
		Values(training.Date, training.Comment, training.Result, training.ResultFiles).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build create training query: %w", err)
	}

	var id int64
	err = r.trm.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("create training: %w", err)
	}

	return id, err
}

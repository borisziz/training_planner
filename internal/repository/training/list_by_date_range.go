package training

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) ListByDateRange(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]model.Training, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		Where("date >= ?", from).
		Where("date < ?", to).
		OrderBy("date", "id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list trainings by date range query: %w", err)
	}

	var trainings []model.Training
	err = r.trm.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &trainings, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list trainings by date range: %w", err)
	}

	return trainings, nil
}

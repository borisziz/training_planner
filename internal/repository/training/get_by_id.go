package training

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) GetTrainingByID(ctx context.Context, trainingID int64) (*model.Training, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		Where("id = ?", trainingID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get training query: %w", err)
	}

	var training model.Training
	err = r.trm.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &training, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get training %d: %w", trainingID, err)
	}

	return &training, err
}

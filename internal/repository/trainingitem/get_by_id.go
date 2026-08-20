package trainingitem

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

// todo fix
func (r *Repository) GetTrainingItemsByTrainingID(
	ctx context.Context,
	trainingItemID int64,
) ([]model.TrainingItem, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		Where("training_id = ?", trainingItemID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get training item query: %w", err)
	}

	var items []model.TrainingItem
	err = r.trm.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get training item %d: %w", trainingItemID, err)
	}

	return items, nil
}

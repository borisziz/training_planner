package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		Where("id = ?", userID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get user query: %w", err)
	}

	var result model.User
	err = r.trm.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &result, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user %d: %w", userID, err)
	}

	return &result, err
}

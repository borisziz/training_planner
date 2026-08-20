package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) List(ctx context.Context) ([]model.User, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		OrderBy("id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list users query: %w", err)
	}

	var users []model.User
	err = r.trm.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, err
}

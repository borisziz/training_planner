package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) CreateUser(ctx context.Context, user model.User) (int64, error) {
	query, args, err := sq.Insert(tableName).
		Columns("name", "surname").
		Values(user.Name, user.Surname).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build create user query: %w", err)
	}

	var id int64
	err = r.trm.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, err
}

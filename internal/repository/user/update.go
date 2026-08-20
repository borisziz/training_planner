package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) Update(ctx context.Context, value model.User) error {
	_, err := sq.Update(tableName).
		Set("name", value.Name).
		Set("surname", value.Surname).
		Where("id = ?", value.ID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("update user %d: %w", value.ID, err)
	}

	return nil
}

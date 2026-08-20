package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteUser(ctx context.Context, userID int64) error {
	_, err := sq.Delete(tableName).
		Where("id = ?", userID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("delete user %d: %w", userID, err)
	}

	return nil
}

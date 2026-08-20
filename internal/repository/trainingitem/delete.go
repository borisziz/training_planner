package trainingitem

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteTrainingItem(ctx context.Context, trainingItemID int64) error {
	_, err := sq.Delete(tableName).
		Where("id = ?", trainingItemID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("delete training item %d: %w", trainingItemID, err)
	}

	return nil
}

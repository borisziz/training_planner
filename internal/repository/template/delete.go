package template

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteTemplate(ctx context.Context, templateID int64) error {
	_, err := sq.Delete(tableName).
		Where("id = ?", templateID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("delete template %d: %w", templateID, err)
	}
	return nil
}

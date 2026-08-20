package template

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) UpdateTemplate(ctx context.Context, value *model.Template) error {
	_, err := sq.Update(tableName).
		Set("type", value.Type).
		Set("name", value.Name).
		Set("description", value.Description).
		Set("video_link", value.VideoLinks).
		Where("id = ?", value.ID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("update template %d: %w", value.ID, err)
	}

	return nil
}

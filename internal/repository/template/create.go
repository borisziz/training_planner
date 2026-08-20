package template

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) CreateTemplate(ctx context.Context, template *model.Template) (int64, error) {
	query, args, err := sq.Insert(tableName).
		Columns("type", "name", "description", "video_links").
		Values(template.Type, template.Name, template.Description, template.VideoLinks).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build create template query: %w", err)
	}

	var id int64

	err = r.trm.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("create template: %w", err)
	}

	return id, err
}

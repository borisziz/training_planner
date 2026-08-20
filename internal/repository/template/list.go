package template

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) ListTemplates(ctx context.Context) ([]model.Template, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		OrderBy("id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list templates query: %w", err)
	}

	var templates []model.Template

	err = r.trm.DefaultTrOrDB(ctx, r.db).SelectContext(ctx, &templates, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}

	return templates, err
}

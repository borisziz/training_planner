package template

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) GetTemplateByID(ctx context.Context, templateID int64) (*model.Template, error) {
	query, args, err := sq.Select(columns...).
		From(tableName).
		Where("id = ?", templateID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get template query: %w", err)
	}

	var result model.Template
	err = r.trm.DefaultTrOrDB(ctx, r.db).GetContext(ctx, &result, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get template %d: %w", templateID, err)
	}

	return &result, err
}

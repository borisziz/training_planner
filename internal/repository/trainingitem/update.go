package trainingitem

import (
	"context"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"training_planner/internal/model"
)

func (r *Repository) Update(ctx context.Context, item model.TrainingItem) error {
	params, err := json.Marshal(item.Params)
	if err != nil {
		return fmt.Errorf("marshal training item params: %w", err)
	}

	_, err = sq.Update(tableName).
		Set("template_id", item.TemplateID).
		Set("training_params", params).
		Where("id = ?", item.ID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("update training item %d: %w", item.ID, err)
	}

	return nil
}

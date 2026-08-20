package trainingitemrelation

import (
	"context"
	"fmt"

	"training_planner/internal/model"
)

func (r *Repository) CreateTrainingItemsRelations(
	ctx context.Context,
	relations []model.TrainingItemRelation,
) error {
	query := `
		INSERT INTO training_item_relations(training_id, training_item_id)
		VALUES (:training_id, :training_item_id)
	`

	_, err := r.trm.DefaultTrOrDB(ctx, r.db).NamedExecContext(ctx, query, relations)
	if err != nil {
		return fmt.Errorf("create training item relations: %w", err)
	}

	return nil
}

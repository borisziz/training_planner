package trainingitem

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"training_planner/internal/model"
)

func (r *Repository) CreateTrainingItems(ctx context.Context, items []model.TrainingItem) ([]int64, error) {

	query := `
		INSERT INTO training_items(template_id, params)
		VALUES (:template_id, :params)
		RETURNING id
	`

	rows, err := r.trm.DefaultTrOrDB(ctx, r.db).NamedQuery(query, items)
	if err != nil {
		return nil, fmt.Errorf("insert items: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("close training item rows")
		}
	}()
	var itemIDs []int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("scan training item row: %w", err)
		}
		itemIDs = append(itemIDs, id)
	}
	return itemIDs, err
}

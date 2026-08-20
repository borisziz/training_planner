package trainingitemrelation

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteTrainingItemRelation(
	ctx context.Context,
	trainingID int64,
	trainingItemID int64,
) error {
	_, err := sq.Delete(tableName).
		Where("training_id = ?", trainingID).
		Where("training_item_id = ?", trainingItemID).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.trm.DefaultTrOrDB(ctx, r.db)).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"delete relation between training %d and item %d: %w",
			trainingID,
			trainingItemID,
			err,
		)
	}

	return nil
}

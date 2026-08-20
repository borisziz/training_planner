package training

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"training_planner/internal/model"
)

func (s *Service) CreateTraining(ctx context.Context, userID int64, training *model.Training) (int64, error) {
	var trainingID int64
	err := s.trm.Do(ctx, func(ctx context.Context) error {
		var err error
		trainingID, err = s.trainingRepo.CreateTraining(ctx, training)
		if err != nil {
			return fmt.Errorf("create training in db: %w", err)
		}
		itemIDs, err := s.trainingItemsRepo.CreateTrainingItems(ctx, training.Items)
		if err != nil {
			return fmt.Errorf("create training items in db: %w", err)
		}
		err = s.trainingItemsRelationsRepo.CreateTrainingItemsRelations(ctx, lo.Map(itemIDs, func(itemID int64, _ int) model.TrainingItemRelation {
			return model.TrainingItemRelation{TrainingID: trainingID, TrainingItemID: itemID}
		}))
		if err != nil {
			return fmt.Errorf("create training item relations in db: %w", err)
		}

		err = s.userTrainingsRepo.CreateUserTraining(ctx, model.UserTraining{TrainingID: trainingID, UserID: userID})
		if err != nil {
			return fmt.Errorf("create user training in db: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return trainingID, nil
}

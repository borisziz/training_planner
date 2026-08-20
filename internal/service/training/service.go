package training

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	"training_planner/internal/model"
)

type (
	TemplateRepo interface {
		GetTemplateByID(ctx context.Context, templateID int64) (*model.Template, error)
	}

	TrainingRepo interface {
		CreateTraining(ctx context.Context, tr *model.Training) (int64, error)
	}

	UserTrainingsRepo interface {
		CreateUserTraining(ctx context.Context, relation model.UserTraining) error
	}

	TrainingItemsRepo interface {
		CreateTrainingItems(ctx context.Context, items []model.TrainingItem) ([]int64, error)
	}

	TrainingItemsRelationsRepo interface {
		CreateTrainingItemsRelations(ctx context.Context, relation []model.TrainingItemRelation) error
	}

	Service struct {
		trainingRepo               TrainingRepo
		trainingItemsRepo          TrainingItemsRepo
		trainingItemsRelationsRepo TrainingItemsRelationsRepo
		userTrainingsRepo          UserTrainingsRepo
		trm                        trm.Manager
	}
)

func New(
	trainingRepo TrainingRepo,
	trainingItemsRepo TrainingItemsRepo,
	trainingItemsRelationsRepo TrainingItemsRelationsRepo,
	userTrainingsRepo UserTrainingsRepo,
	transactionManager trm.Manager,
) *Service {
	return &Service{
		trainingRepo:               trainingRepo,
		trainingItemsRepo:          trainingItemsRepo,
		trainingItemsRelationsRepo: trainingItemsRelationsRepo,
		userTrainingsRepo:          userTrainingsRepo,
		trm:                        transactionManager,
	}
}

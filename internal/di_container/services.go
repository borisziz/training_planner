package di_container

import (
	"context"

	"training_planner/internal/service/training"
)

func (c *cont) GetTrainingService(ctx context.Context) *training.Service {
	return makeSingleton(&c.trainingSrv, func() *training.Service {
		return training.New(
			c.GetTrainingRepo(ctx),
			c.GetTrainingItemRepo(ctx),
			c.GetTrainingItemRelationRepo(ctx),
			c.GetUserTrainingsRepo(ctx),
			c.GetTrmManager(ctx),
		)
	})
}

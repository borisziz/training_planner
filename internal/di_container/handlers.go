package di_container

import (
	"context"

	"training_planner/internal/api/create_training"
)

func (c *cont) GetCreateTrainingHandler(ctx context.Context) *create_training.Handler {
	return makeSingleton(&c.createTrainingHandler, func() *create_training.Handler {
		return create_training.New(c.GetTrainingService(ctx))
	})
}

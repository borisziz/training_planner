package di_container

import (
	"context"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"training_planner/internal/postgres"
	"training_planner/internal/repository/training"
	"training_planner/internal/repository/trainingitem"
	"training_planner/internal/repository/trainingitemrelation"
	"training_planner/internal/repository/usertraining"
)

func (c *cont) GetPostgres(ctx context.Context) postgres.Client {
	return makeSingleton(&c.pg, func() postgres.Client {
		client, err := postgres.New(ctx, c.cfg.DBConnectURL)
		if err != nil {
			panic(err)
		}
		c.addCancel(func() {
			_ = client.Close()
		})
		return client
	})
}

func (c *cont) GetTrmManager(ctx context.Context) trm.Manager {
	return makeSingleton(&c.trm, func() trm.Manager {
		tm, err := manager.New(trmsqlx.NewDefaultFactory(c.GetPostgres(ctx)))
		if err != nil {
			panic(fmt.Errorf("failed to create transaction manager: %w", err))
		}
		return tm
	})
}

func (c *cont) GetTrainingRepo(ctx context.Context) *training.Repository {
	return makeSingleton(&c.trainingRepo, func() *training.Repository {
		return training.New(c.GetPostgres(ctx), trmsqlx.DefaultCtxGetter)
	})
}

func (c *cont) GetTrainingItemRepo(ctx context.Context) *trainingitem.Repository {
	return makeSingleton(&c.trainingItemRepo, func() *trainingitem.Repository {
		return trainingitem.New(c.GetPostgres(ctx), trmsqlx.DefaultCtxGetter)
	})
}

func (c *cont) GetTrainingItemRelationRepo(ctx context.Context) *trainingitemrelation.Repository {
	return makeSingleton(&c.trainingItemRelationsRepo, func() *trainingitemrelation.Repository {
		return trainingitemrelation.New(c.GetPostgres(ctx), trmsqlx.DefaultCtxGetter)
	})
}

func (c *cont) GetUserTrainingsRepo(ctx context.Context) *usertraining.Repository {
	return makeSingleton(&c.userTrainingsRepo, func() *usertraining.Repository {
		return usertraining.New(c.GetPostgres(ctx), trmsqlx.DefaultCtxGetter)
	})
}

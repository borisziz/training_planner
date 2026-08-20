package di_container

import (
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2"

	"training_planner/internal/api/create_training"
	"training_planner/internal/config"
	"training_planner/internal/postgres"
	"training_planner/internal/repository/training"
	"training_planner/internal/repository/trainingitem"
	"training_planner/internal/repository/trainingitemrelation"
	"training_planner/internal/repository/user"
	"training_planner/internal/repository/usertraining"
	training_srv "training_planner/internal/service/training"
)

type cont struct {
	cfg *config.Config

	// repo
	pg                        postgres.Client
	trm                       trm.Manager
	trainingRepo              *training.Repository
	trainingItemRepo          *trainingitem.Repository
	trainingItemRelationsRepo *trainingitemrelation.Repository
	usersRepo                 *user.Repository
	userTrainingsRepo         *usertraining.Repository

	// services
	trainingSrv *training_srv.Service

	// handlers
	createTrainingHandler *create_training.Handler

	// server
	server *http.Server

	closeFuncs []func()
}

func New() *cont {
	c := &cont{}
	var err error
	c.cfg, err = config.NewConfig()
	if err != nil {
		panic(err)
	}
	return c
}

func (c *cont) GetConfig() *config.Config {
	return c.cfg
}

func (c *cont) Close() {
	for i := len(c.closeFuncs) - 1; i >= 0; i-- {
		c.closeFuncs[i]()
	}
}

func (c *cont) addCancel(f func()) {
	c.closeFuncs = append(c.closeFuncs, f)
}

func makeSingleton[T comparable](ptr *T, constructor func() T) T {
	if ptr == nil {
		panic("nil reference instead of container field")
	}
	var zero T
	if *ptr == zero {
		*ptr = constructor()
	}
	return *ptr
}

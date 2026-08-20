package training

import trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

const tableName = "trainings"

var columns = []string{"id", "date", "comment", "result", "result_video_links"}

type Repository struct {
	db  trmsqlx.Tr
	trm *trmsqlx.CtxGetter
}

func New(
	db trmsqlx.Tr,
	trm *trmsqlx.CtxGetter,
) *Repository {
	return &Repository{
		db:  db,
		trm: trm,
	}
}

package trainingitem

import trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

const tableName = "training_items"

var columns = []string{"id", "template_id", "params"}

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

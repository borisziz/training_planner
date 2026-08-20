package trainingitemrelation

import trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

const tableName = "training_item_relations"

var columns = []string{"training_id", "training_item_id"}

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

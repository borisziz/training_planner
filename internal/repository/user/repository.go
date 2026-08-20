package user

import trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"

const tableName = "users"

var columns = []string{"id", "name", "surname"}

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

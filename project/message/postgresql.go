package message

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresqlClient(dsn string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

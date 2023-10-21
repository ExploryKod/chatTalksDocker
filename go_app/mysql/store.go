package database

import (
	"chatHTTP"
	"database/sql"
)

func CreateStore(db *sql.DB) *Store {
	return &Store{
		NewUserStore(db),
	}
}

type Store struct {
	chatHTTP.UserStoreInterface
}

package sqlRepo

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

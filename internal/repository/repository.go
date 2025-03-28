package repository

import (
	"database/sql"
	"prosto-delaj-api/models"
	"prosto-delaj-api/server/queries"
)

type Sources struct {
	BusinessDB *sql.DB
}

type Test interface {
	Create(input models.CreateInput) error
}

type Repository struct {
	Test
}

func NewRepository(queries *queries.Queries) *Repository {
	return &Repository{
		Test: NewTestRepository(queries),
	}
}

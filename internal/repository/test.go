package repository

import (
	"prosto-delaj-api/models"
	"prosto-delaj-api/server/queries"
)

type TestRepository struct {
	queries *queries.Queries
}

func NewTestRepository(queries *queries.Queries) *TestRepository {
	return &TestRepository{queries: queries}
}

func (r *TestRepository) Create(input models.CreateInput) error {
	panic("implement me")
	return nil
}

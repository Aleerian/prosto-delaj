package service

import (
	"prosto-delaj-api/internal/repository"
	"prosto-delaj-api/models"
)

type Test interface {
	Create(input models.CreateInput) error
}

type Service struct {
	Test
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Test: NewTestService(repo.Test),
	}
}

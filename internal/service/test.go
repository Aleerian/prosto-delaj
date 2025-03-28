package service

import (
	"prosto-delaj-api/internal/repository"
	"prosto-delaj-api/models"
)

type TestService struct {
	repo repository.Test
}

func NewTestService(repo repository.Test) *TestService {
	return &TestService{repo: repo}
}

func (s *TestService) Create(input models.CreateInput) error {
	return s.repo.Create(input)
}

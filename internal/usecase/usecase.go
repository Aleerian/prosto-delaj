package usecase

import "prosto-delaj-api/internal/service"

type Usecase struct {
	services *service.Service
}

func NewUsecase(service *service.Service) *Usecase {
	return &Usecase{services: service}
}

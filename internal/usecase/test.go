package usecase

import (
	"github.com/sirupsen/logrus"
	"prosto-delaj-api/internal/exceptions"
	"prosto-delaj-api/models"
)

func (u *Usecase) Create(input models.CreateInput) exceptions.ErrorCode {
	err := u.services.Create(input)
	if err != nil {
		logrus.Error(err)
		return exceptions.InternalServerError
	}
	return exceptions.Success
}

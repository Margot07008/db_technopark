package usecase

import (
	"db_technopark/application/models"
	"db_technopark/application/service"
)

type serviceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(serviceRepo service.Repository) service.Usecase {
	return &serviceUsecase{serviceRepo: serviceRepo}
}

func (u serviceUsecase) GetDBStatus() (models.Status, *models.Error) {
	return u.serviceRepo.GetStatus()
}

func (u serviceUsecase) ClearDB() *models.Error {
	return u.serviceRepo.Clear()
}

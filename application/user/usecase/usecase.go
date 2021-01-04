package usecase

import (
	"db_technopark/application/models"
	"db_technopark/application/user"
)

type userUsecase struct {
	userRepo user.Repository
}

func NewUserUsecase(userRepo user.Repository) user.Usecase {
	return &userUsecase{userRepo: userRepo}
}

func (u userUsecase) CreateUser(userNew models.User) (models.User, *models.Error) {
	return u.userRepo.Create(userNew)
}
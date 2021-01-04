package user

import "db_technopark/application/models"

type Usecase interface {
	CreateUser(userNew models.User) (models.User, *models.Error)
	GetUserByNickname(nickname string) (models.User, *models.Error)
	GetUserByEmail(email string) (models.User, *models.Error)
	UpdateUser(userUpd models.User) (models.User, *models.Error)
}
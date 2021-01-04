package user

import "db_technopark/application/models"

type Usecase interface {
	CreateUser(userNew models.User) (models.User, *models.Error)
}
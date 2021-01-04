package repository

import (
	"db_technopark/application/models"
	"db_technopark/application/user"
	"github.com/jackc/pgx"
)

type pgUserRepository struct {
	conn *pgx.ConnPool
}

func NewPgUserRepository(db *pgx.ConnPool) user.Repository {
	return &pgUserRepository{conn: db}
}

func (p pgUserRepository) Create(userNew models.User) (models.User, *models.Error) {
	res, err := p.conn.Exec(`INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		userNew.Nickname, userNew.Fullname, userNew.Email, userNew.About)
	if err != nil {
		return models.User{}, models.NewError(409, models.CreateError)
	}

	if res.RowsAffected() == 0 {
		return models.User{}, models.NewError(409, models.CreateError)
	}

	return userNew, nil
}

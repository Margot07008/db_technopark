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

func (p pgUserRepository) GetByNickname(nickname string) (models.User, *models.Error) {
	modelUser := models.User{}
	res, err := p.conn.Query(`SELECT nickname, email, fullname, about FROM main.users WHERE nickname = $1`, nickname)
	if err != nil {
		return models.User{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()
	if res.Next() {
		err = res.Scan(&modelUser.Nickname, &modelUser.Email, &modelUser.Fullname, &modelUser.About)
		if err != nil {
			return models.User{}, models.NewError(404, models.DBParsingError)
		}
		return modelUser, nil
	}
	return models.User{}, models.NewError(404, models.NotFoundError)
}

func (p pgUserRepository) GetByEmail(email string) (models.User, *models.Error) {
	modelUser := models.User{}
	res, err := p.conn.Query(`SELECT nickname, email, fullname, about FROM main.users WHERE email = $1`, email)
	if err != nil {
		return models.User{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()
	if res.Next() {
		err = res.Scan(&modelUser.Nickname, &modelUser.Email, &modelUser.Fullname, &modelUser.About)
		if err != nil {
			return models.User{}, models.NewError(404, models.DBParsingError)
		}
		return modelUser, nil
	}
	return models.User{}, models.NewError(404, models.NotFoundError)
}

func (p pgUserRepository) Create(userNew models.User) (models.User, *models.Error) {
	res, err := p.conn.Exec(`INSERT INTO main.users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		userNew.Nickname, userNew.Fullname, userNew.Email, userNew.About)
	if err != nil {
		return models.User{}, models.NewError(409, models.CreateError)
	}

	if res.RowsAffected() == 0 {
		return models.User{}, models.NewError(409, models.CreateError)
	}

	return userNew, nil
}

func (p pgUserRepository) Update(userUpd models.User) (models.User, *models.Error) {
	if userUpd.Nickname == "" && userUpd.Fullname == "" && userUpd.Email == "" && userUpd.About == "" {
		updatedUser, _ := p.GetByNickname(userUpd.Nickname)
		return updatedUser, nil
	}

	querySQL := "UPDATE main.users SET"

	if userUpd.Email != "" {
		querySQL += " email = '" + userUpd.Email + "',"
	} else {
		querySQL += " email = email,"
	}

	if userUpd.Fullname != "" {
		querySQL += " fullname = '" + userUpd.Fullname + "',"
	} else {
		querySQL += " fullname = fullname,"
	}

	if userUpd.About != "" {
		querySQL += " about = '" + userUpd.About + "'"
	} else {
		querySQL += " about = about,"
	}

	querySQL += " where nickname = '" + userUpd.Nickname + "'"
	res, err := p.conn.Exec(querySQL)
	if err != nil {
		return models.User{}, models.NewError(404, models.NotFoundError)
	}
	if res.RowsAffected() == 0 {
		return models.User{}, models.NewError(404, models.NotFoundError)
	}
	updatedUser, _ := p.GetByNickname(userUpd.Nickname)
	return updatedUser, nil
}

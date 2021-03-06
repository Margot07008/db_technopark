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

func (p pgUserRepository) GetByForum(forum models.Forum, query models.PostsRequestQuery) (models.Users, *models.Error) {
	baseSQL := `select about, email, fullname, u.nickname from main.users_forum m join main.users u on u.nickname = m.nickname`

	baseSQL += ` where slug = '` + forum.Slug + `'`
	if query.Since != "" {
		if query.Desc {
			baseSQL += ` and lower(u.nickname)::bytea < lower('` + query.Since + `')::bytea`
		} else {
			baseSQL += ` and lower(u.nickname)::bytea > lower('` + query.Since + `')::bytea`
		}
	}

	if query.Desc {
		baseSQL += " order by lower(u.nickname) COLLATE \"C\" desc"
	} else {
		baseSQL += " order by lower(u.nickname) COLLATE \"C\" asc"
	}

	if query.Limit > 0 {
		baseSQL += " limit " + query.GetStringLimit()
	}

	res, err := p.conn.Query(baseSQL)
	if err != nil {
		return models.Users{}, models.NewError(500, models.DBParsingError)
	}
	defer res.Close()

	foundUsers := models.Users{}
	buffer := models.User{}

	for res.Next() {
		err = res.Scan(&buffer.About, &buffer.Email, &buffer.Fullname, &buffer.Nickname)

		if err != nil {
			return models.Users{}, models.NewError(500, models.DBParsingError)
		}
		foundUsers = append(foundUsers, buffer)
	}

	return foundUsers, nil
}

func (p pgUserRepository) GetByNickname(nickname string) (models.User, *models.Error) {
	modelUser := models.User{}
	res, err := p.conn.Query(`select nickname, email, fullname, about from main.users where lower(nickname) = lower($1)`, nickname)
	if err != nil {
		return models.User{}, models.NewError(500, models.InternalError)
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
	res, err := p.conn.Query(`select nickname, email, fullname, about from main.users where lower(email) = lower($1)`, email)
	if err != nil {
		return models.User{}, models.NewError(500, models.InternalError)
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
	res, err := p.conn.Exec(`insert into main.users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
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
	if userUpd.Fullname == "" && userUpd.Email == "" && userUpd.About == "" {
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
		querySQL += " about = about"
	}

	querySQL += " where nickname = '" + userUpd.Nickname + "'"
	res, err := p.conn.Exec(querySQL)
	if err != nil {
		return models.User{}, models.NewError(409, models.UpdateError)
	}
	if res.RowsAffected() == 0 {
		_, err := p.GetByEmail(userUpd.Email)
		if err != nil {
			return models.User{}, models.NewError(404, models.NotFoundError)
		}
		return models.User{}, models.NewError(409, models.UpdateError)

	}
	updatedUser, _ := p.GetByNickname(userUpd.Nickname)
	return updatedUser, nil
}

func (p pgUserRepository) AddUserToForum(nickname string, forum string) {
	_, _ = p.conn.Exec(`insert into main.users_forum (nickname, slug) values ($1, $2) on conflict do nothing`,
		nickname, forum)
	return
}

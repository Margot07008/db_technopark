package repository

import (
	"db_technopark/application/models"
	"db_technopark/application/service"
	"github.com/jackc/pgx"
)

type pgServiceRepository struct {
	conn *pgx.ConnPool
}

func NewPgServiceRepository(conn *pgx.ConnPool) service.Repository {
	return &pgServiceRepository{conn: conn}
}

func (p pgServiceRepository) Clear() *models.Error {
	res, err := p.conn.Query("truncate table main.users, main.forums, main.threads, main.posts, main.votes, main.users_forum cascade")
	if err != nil {
		return models.NewError(500, models.InternalError)
	}
	defer res.Close()
	return nil
}

func (p pgServiceRepository) GetStatus() (models.Status, *models.Error) {
	res, err := p.conn.Query("select * from (select count(posts) from main.forums) as f" +
		" cross join (select count(id) from main.posts) as p" +
		" cross join (select count(id) from main.threads) as t" +
		" cross join (select count(nickname) from main.users) as u")

	if err != nil {
		return models.Status{}, models.NewError(500, models.InternalError)
	}
	defer res.Close()

	s := models.Status{}
	for res.Next() {
		err = res.Scan(&s.Forum, &s.Post, &s.Thread, &s.User)
		if err != nil {
			return models.Status{}, models.NewError(500, models.InternalError)
		}
		return s, nil
	}
	return models.Status{}, models.NewError(500, models.InternalError)
}

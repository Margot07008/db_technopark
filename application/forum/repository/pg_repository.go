package repository

import (
	"db_technopark/application/forum"
	"db_technopark/application/models"
	"github.com/jackc/pgx"
)

type pgForumRepository struct {
	conn *pgx.ConnPool
}

func NewPgForumRepository(db *pgx.ConnPool) forum.Repository {
	return &pgForumRepository{conn: db}
}

func (p pgForumRepository) CreateForum(userNick string, forumNew models.Forum) (models.Forum, *models.Error) {
	forumNew.User = userNick
	_, err := p.conn.Exec(`INSERT INTO main.forums (slug, title, "user", posts, threads) VALUES ($1, $2, $3, $4, $5)`, forumNew.Slug, forumNew.Title, forumNew.User, forumNew.Posts, forumNew.Threads)
	if err != nil {
		return models.Forum{}, models.NewError(409, models.CreateError)
	}
	return forumNew, nil
}

func (p pgForumRepository) GetForumBySlug(slug string) (models.Forum, *models.Error) {

	res, err := p.conn.Query(`select title, "user", slug, posts, threads from main.forums where lower(slug) = lower($1)`, slug)
	if err != nil {
		return models.Forum{}, models.NewError(500, models.InternalError)
	}
	defer res.Close()
	forumModel := models.Forum{}
	if res.Next() {
		err = res.Scan(&forumModel.Title, &forumModel.User, &forumModel.Slug, &forumModel.Posts, &forumModel.Threads)
		if err != nil {
			return models.Forum{}, models.NewError(500, models.DBParsingError)
		}
		return forumModel, nil
	}
	return models.Forum{}, models.NewError(404, models.NotFoundError)

}

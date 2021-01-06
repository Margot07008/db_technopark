package repository

import (
	"db_technopark/application/models"
	"db_technopark/application/thread"
	"github.com/jackc/pgx"
)

type pgThreadRepository struct {
	conn *pgx.ConnPool
}

func NewPgThreadRepository(conn *pgx.ConnPool) thread.Repository {
	return &pgThreadRepository{conn: conn}
}

func (p pgThreadRepository) CreateThread(forum models.Forum, user models.User, thread models.Thread) (models.Thread, *models.Error) {
	thread.Forum = forum.Slug
	thread.Author = user.Nickname

	tx, _ := p.conn.Begin()
	defer tx.Rollback()
	if thread.Slug == "" {
		err := tx.QueryRow(`insert into main.threads (author, created, forum, message, title) values ($1, $2, $3, $4, $5) returning id`,
			thread.Author, thread.Created, thread.Forum, thread.Message,
			thread.Title).Scan(&thread.Id)
		if err == pgx.ErrNoRows || err != nil {
			return models.Thread{}, models.NewError(409, models.ConflictError)
		}

	} else {
		err := tx.QueryRow(`insert into main.threads (author, created, forum, message, slug, title) values ($1, $2, $3, $4, $5, $6) returning id`,
			thread.Author, thread.Created, thread.Forum, thread.Message, thread.Slug,
			thread.Title).Scan(&thread.Id)
		if err == pgx.ErrNoRows || err != nil {
			return models.Thread{}, models.NewError(409, models.ConflictError)
		}

	}

	err := tx.Commit()
	if err != nil {
		return models.Thread{}, models.NewError(500, models.InternalError)
	}

	return thread, nil
}

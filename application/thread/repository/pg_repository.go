package repository

import (
	"db_technopark/application/models"
	"db_technopark/application/thread"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
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

func (p pgThreadRepository) GetThreadsBySlug(forum models.Forum, query models.PostsRequestQuery) (models.Threads, *models.Error) {
	baseSQL := "select * from main.threads"
	baseSQL += " where forum = '" + forum.Slug + "'"

	if query.Since != "" {
		if query.Desc {
			baseSQL += " and created <= '" + query.Since + "'"
		} else {
			baseSQL += " and created >= '" + query.Since + "'"
		}
	}

	if query.Desc {
		baseSQL += " order by created desc"
	} else {
		baseSQL += " order by created"
	}

	if query.Limit > 0 {
		baseSQL += " limit " + query.GetStringLimit()
	}
	res, err := p.conn.Query(baseSQL)
	if err != nil {
		return models.Threads{}, models.NewError(500, models.DBParsingError)
	}
	buffer := models.Thread{}
	foundedThreads := models.Threads{}
	for res.Next() {
		nullSlug := pgtype.Varchar{}
		err = res.Scan(&buffer.Id, &buffer.Title, &buffer.Author, &buffer.Forum, &buffer.Message, &buffer.Votes, &nullSlug, &buffer.Created)
		if err != nil {
			return models.Threads{}, models.NewError(500, models.InternalError)
		}
		buffer.Slug = nullSlug.String
		foundedThreads = append(foundedThreads, buffer)
	}

	return foundedThreads, nil
}

func (p pgThreadRepository) GetBySlug(slug string) (models.Thread, *models.Error) {
	if slug == "" {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}

	fThread := models.Thread{}
	res, err := p.conn.Query(`select * from main.threads where lower(slug) = lower($1)`, slug)

	if err != nil {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&fThread.Id, &fThread.Title, &fThread.Author, &fThread.Forum, &fThread.Message, &fThread.Votes, &fThread.Slug, &fThread.Created)
		if err != nil {
			return models.Thread{}, models.NewError(500, models.DBParsingError)
		}
		return fThread, nil
	}
	return models.Thread{}, models.NewError(404, models.NotFoundError)
}

func (p pgThreadRepository) GetByID(id int32) (models.Thread, *models.Error) {
	fThread := models.Thread{}

	if id == -1 {
		return models.Thread{}, models.NewError(400, models.BadRequestError)
	}

	res, err := p.conn.Query(`select * from main.threads where id = $1`, id)
	if err != nil {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	nullSlug := &pgtype.Varchar{}
	if res.Next() {
		err = res.Scan(&fThread.Id, &fThread.Title, &fThread.Author, &fThread.Forum, &fThread.Message, &fThread.Votes, nullSlug, &fThread.Created)
		if err != nil {
			return models.Thread{}, models.NewError(500, models.InternalError)
		}

		fThread.Slug = nullSlug.String
		return fThread, nil
	}

	return models.Thread{}, models.NewError(404, models.NotFoundError)
}

func (p pgThreadRepository) Update(thread models.Thread, threadUpdate models.ThreadUpdate) (models.Thread, *models.Error) {
	if threadUpdate.Message == "" && threadUpdate.Title == "" {
		return thread, nil
	}

	baseSQL := "update main.threads set"
	if threadUpdate.Message == "" {
		baseSQL += " message = message,"
	} else {
		thread.Message = threadUpdate.Message
		baseSQL += " message = '" + threadUpdate.Message + "',"
	}

	if threadUpdate.Title == "" {
		baseSQL += " title = title"
	} else {
		thread.Title = threadUpdate.Title
		baseSQL += " title = '" + threadUpdate.Title + "'"
	}

	baseSQL += " where slug = '" + thread.Slug + "'"

	res, err := p.conn.Exec(baseSQL)
	if err != nil {
		return models.Thread{}, models.NewError(500, models.UpdateError)
	}

	if res.RowsAffected() == 0 {
		return models.Thread{}, models.NewError(404, models.NotFoundError)
	}

	return thread, nil
}

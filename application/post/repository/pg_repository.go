package repository

import (
	"db_technopark/application/models"
	"db_technopark/application/post"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/lib/pq"
	"strings"
)

type pgPostRepository struct {
	conn *pgx.ConnPool
}

func NewPgPostRepository(conn *pgx.ConnPool) post.Repository {
	return &pgPostRepository{conn: conn}
}

func (p pgPostRepository) CreatePosts(posts models.Posts, thread models.Thread) (models.Posts, *models.Error) {
	if len(posts) == 0 {
		return models.Posts{}, nil
	}

	tx, _ := p.conn.Begin()
	defer tx.Rollback()

	mapParents := make(map[int64]models.Post)

	for _, item := range posts {
		if _, ok := mapParents[item.Parent]; !ok && item.Parent != 0 {
			parentPostQuery, err := p.GetById(item.Parent)
			if err != nil {
				err.StatusCode = 409
				return models.Posts{}, err
			}

			if parentPostQuery.Thread != thread.Id {
				return models.Posts{}, models.NewError(409, models.BadRequestError)
			}

			mapParents[item.Parent] = parentPostQuery
		}
	}

	postIdsRows, err := tx.Query(fmt.Sprintf(`select nextval(pg_get_serial_sequence('main.posts', 'id')) from generate_series(1, %d);`, len(posts)))
	if err != nil {
		return models.Posts{}, models.NewError(404, models.NotFoundError)
	}

	var postIds []int64
	for postIdsRows.Next() {
		var availableId int64
		_ = postIdsRows.Scan(&availableId)
		postIds = append(postIds, availableId)
	}
	postIdsRows.Close()

	if len(postIds) == 0 {
		return models.Posts{}, models.NewError(500, models.DBError)
	}

	posts[0].Path = append(mapParents[posts[0].Parent].Path, postIds[0])
	err = tx.QueryRow(`insert into main.posts (id, author, forum, message, parent, thread, path) values ($1, $2, $3, $4, $5, $6, $7) returning created`,
		postIds[0], posts[0].Author, thread.Forum, posts[0].Message, posts[0].Parent,
		thread.Id,
		"{"+strings.Trim(strings.Replace(fmt.Sprint(posts[0].Path), " ", ",", -1), "[]")+"}").
		Scan(&posts[0].Created)

	if err != nil {
		return models.Posts{}, models.NewError(404, models.CreateError)
	}

	now := posts[0].Created
	posts[0].Forum = thread.Forum
	posts[0].Thread = thread.Id
	posts[0].Created = now
	posts[0].Id = postIds[0]

	for i, item := range posts {
		if i == 0 {
			continue
		}

		item.Path = append(mapParents[item.Parent].Path, postIds[i])
		resInsert, err := tx.Exec(`insert into main.posts (id, author, created, forum, message, parent, thread, path) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
			postIds[i], item.Author, now, thread.Forum, item.Message, item.Parent, thread.Id,
			"{"+strings.Trim(strings.Replace(fmt.Sprint(item.Path), " ", ",", -1), "[]")+"}")

		if err != nil {
			return models.Posts{}, models.NewError(500, models.CreateError)
		}

		if resInsert.RowsAffected() == 0 {
			return models.Posts{}, models.NewError(500, models.CreateError)
		}

		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.Id
		posts[i].Created = now
		posts[i].Id = postIds[i]
	}

	_, err = tx.Exec(`update main.forums set posts = posts + $1 where slug = $2`, len(posts), thread.Forum)
	if err != nil {
		return models.Posts{}, models.NewError(500, models.InternalError)
	}

	err = tx.Commit()

	if err != nil {
		return models.Posts{}, models.NewError(500, models.DBError)
	}

	return posts, nil
}

func (p pgPostRepository) GetById(id int64) (models.Post, *models.Error) {
	res, err := p.conn.Query("select * from main.posts where id = $1", id)
	if err != nil {
		return models.Post{}, models.NewError(404, models.NotFoundError)
	}
	defer res.Close()

	fPost := models.Post{}

	for res.Next() {
		err = res.Scan(&fPost.Id, &fPost.Forum,
			&fPost.Parent, &fPost.Author,
			&fPost.Created, &fPost.IsEdited,
			&fPost.Message, &fPost.Thread, pq.Array(&fPost.Path))

		if err != nil {
			return models.Post{}, models.NewError(500, models.DBParsingError)
		}

		return fPost, nil
	}

	return models.Post{}, models.NewError(404, models.NotFoundError)
}

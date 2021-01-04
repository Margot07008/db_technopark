package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/forum"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
)

type PGRepository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) forum.Repository {
	return &PGRepository{db: db}
}

func (p *PGRepository) CreateForum(forum models.Forum) (*models.Forum, error) {
	query := fmt.Sprintf(`insert into main.forum 
					(slug, title, user) values ('%s', '%s', '%s')`, forum.Slug, forum.Title, forum.User)

	_, err := p.db.Query(&forum, query)
	if err != nil {
		return nil, err
	}

	return &forum, nil
}

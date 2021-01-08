package postsSQLGen

import (
	"db_technopark/application/models"
	"strconv"
)

type PostsSQLGen interface {
	FlatSort() string
	TreeSort() string
	ParentTreeSort() string
}

type postsSQLGen struct {
	thread models.Thread
	query  models.PostsRequestQuery
}

func NewPostsSQLGen(thread models.Thread, query models.PostsRequestQuery) PostsSQLGen {
	return &postsSQLGen{
		thread: thread,
		query:  query,
	}
}

func (p postsSQLGen) FlatSort() string {
	strID := strconv.FormatInt(int64(p.thread.Id), 10)
	baseSQL := ""

	baseSQL = "select author, created, forum, id, isedited, message, parent, thread FROM main.posts where thread = " + strID

	if p.query.Since != "" {
		if p.query.Desc {
			baseSQL += " and id < " + p.query.Since
		} else {
			baseSQL += " and id > " + p.query.Since
		}
	}

	if p.query.Desc {
		baseSQL += " order by id desc"
	} else {
		baseSQL += " order by id asc"
	}

	baseSQL += " limit " + strconv.Itoa(p.query.Limit)

	return baseSQL
}

func (p postsSQLGen) TreeSort() string {
	strID := strconv.FormatInt(int64(p.thread.Id), 10)
	baseSQL := ""

	baseSQL = "select author, created, forum, id, isedited, message, parent, thread from main.posts where thread = " + strID

	if p.query.Since != "" {
		if p.query.Desc {
			baseSQL += " and path < (select path from main.posts where id = " + p.query.Since + ")"
		} else {
			baseSQL += " and path > (select path from main.posts where id = " + p.query.Since + ")"
		}
	}

	if p.query.Desc {
		baseSQL += " order by path desc, id desc"
	} else {
		baseSQL += " order by path asc, id asc"
	}

	baseSQL += " limit " + strconv.Itoa(p.query.Limit)

	return baseSQL
}

func (p postsSQLGen) ParentTreeSort() string {
	baseSQL := ""

	baseSQL = "select author, created, forum, id, isedited, message, parent, thread from main.posts where path[1]" +
		" in (select id from main.posts where thread = " + strconv.FormatInt(int64(p.thread.Id), 10) +
		" and parent = 0"

	if p.query.Since != "" {
		if p.query.Desc {
			baseSQL += " and path[1] < (select path[1] from main.posts where id = " + p.query.Since + ")"
		} else {
			baseSQL += " and path[1] > (select path[1] from main.posts where id = " + p.query.Since + ")"
		}
	}

	if p.query.Desc {
		baseSQL += " order by id desc"
	} else {
		baseSQL += " order by id asc"
	}

	baseSQL += " limit " + strconv.Itoa(p.query.Limit) + ")"

	if p.query.Desc {
		baseSQL += " order by path[1] desc, path asc, id asc"
	} else {
		baseSQL += " order by path asc"
	}

	return baseSQL
}

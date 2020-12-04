package threadRepository

import (
	"database/sql"
	"fmt"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"net/url"
	"strconv"
)

type IThreadRepository interface {
	Create(thread *models.Thread) (*models.Thread, error)
	Get(slugOrId string) (*models.Thread, error)
	GetPosts(thread *models.Thread, params url.Values) (models.Posts, error)
	Update(slugOrId string, updateThread *models.ThreadUpdate) (*models.Thread, error)
	Vote(thread *models.Thread, vote *models.Vote) (*models.Thread, error)
}

type ThreadRepository struct {
	DB *sql.DB
}

func (t *ThreadRepository) Create(thread *models.Thread) (*models.Thread, error) {
	newThread := &models.Thread{}
	err := t.DB.QueryRow(
		"INSERT INTO threads (author, forum, message, slug, title, created) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, forum, author, created, message, title, votes, slug",
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Slug,
			thread.Title,
			thread.Created,
	).Scan(
			&newThread.Id,
			&newThread.Forum,
			&newThread.Author,
			&newThread.Created,
			&newThread.Message,
			&newThread.Title,
			&newThread.Votes,
			&newThread.Slug,
	)

	if err != nil {
		return nil, err
	}

	_, err = t.DB.Exec("UPDATE forums SET threads = threads + 1 WHERE slug = $1",
			thread.Forum,
		)

	if err != nil {
		return nil, err
	}

	return newThread, nil
}

func (t *ThreadRepository) Get(slugOrId string) (*models.Thread, error) {
	searchKey := ""
	if _, err := strconv.Atoi(slugOrId); err != nil {
		searchKey = "slug = $1 "
	} else {
		searchKey = "id = $1 "
	}

	thread := &models.Thread{}
	query := "SELECT id, forum, author, created, message, slug, title, votes " +
		"FROM threads WHERE " + searchKey

	err := t.DB.QueryRow(query,
			slugOrId,
		).Scan(
			&thread.Id,
			&thread.Forum,
			&thread.Author,
			&thread.Created,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (t *ThreadRepository) Update(slugOrId string, updateThread *models.ThreadUpdate) (*models.Thread, error) {
	searchKey := ""
	if _, err := strconv.Atoi(slugOrId); err != nil {
		searchKey = "slug = $3 "
	} else {
		searchKey = "id = $3 "
	}

	thread := &models.Thread{}
	err := t.DB.QueryRow(
		"UPDATE threads SET title = $1, message = $2 WHERE " +
			searchKey +
			"RETURNING id, slug, title, message, forum, author, created, votes",
			updateThread.Title,
			updateThread.Message,
			slugOrId,
		).Scan(
			&thread.Id,
			&thread.Slug,
			&thread.Title,
			&thread.Message,
			&thread.Forum,
			&thread.Author,
			&thread.Created,
			&thread.Votes,
		)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (t *ThreadRepository) Vote(thread *models.Thread, newVote *models.Vote) (*models.Thread, error) {
	err := t.DB.QueryRow("SELECT nickname FROM users WHERE nickname = $1", newVote.Nickname).
		Scan(&newVote.Nickname)  // check if exists
	if err != nil {
		return nil, err
	}

	vote := &models.Vote{}
	err = t.DB.QueryRow(  // check if exists
		"SELECT voice FROM votes WHERE thread = $1 AND nickname = $2",
			thread.Id,
			newVote.Nickname,
		).Scan(
			&vote.Voice,
		)
	if err != nil {  // создаем голос
		_, err = t.DB.Exec(
			"INSERT INTO votes(thread, voice, nickname)" +
			"VALUES ($1, $2, $3)",
			thread.Id,
			newVote.Voice,
			newVote.Nickname,
		)

		_, err = t.DB.Exec(  // TODO  узнать
			"UPDATE threads SET votes = votes + $1" +
				"WHERE id = $2",
			newVote.Voice,
			thread.Id,
		)

		if newVote.Voice < 0 {
			thread.Votes--
		} else {
			thread.Votes++
		}
	} else if vote.Voice != newVote.Voice {
		_, err = t.DB.Exec(
			"UPDATE votes SET voice = $1" +
				"WHERE thread = $2 AND nickname = $3",
			newVote.Voice,
			thread.Id,
			newVote.Nickname,
		)

		_, err = t.DB.Exec(
			"UPDATE threads SET votes = votes + CASE WHEN $1 < 0 THEN -2 ELSE 2 END " +
				" WHERE id = $2",
			newVote.Voice,
			thread.Id,
		)

		if newVote.Voice < 0 {
			thread.Votes -= 2
		} else {
			thread.Votes += 2
		}
	}

	return thread, err
}

func (t *ThreadRepository) GetPosts(thread *models.Thread, params url.Values) (models.Posts, error) {
	limit, order, sortType, sincePost, sortSign := extractParams(params)

	var sortedRows *sql.Rows
	var err error

	switch sortType {
	case "flat":
		sortedRows, err = t.getFlatSortedRows(thread.Id, order, sincePost, limit, sortSign)
	case "tree":
		sortedRows, err = t.getTreeSortedRows(thread.Id, order, sincePost, limit, sortSign)
	case "parent_tree":
		sortedRows, err = t.getParentTreeSortedRows(thread.Id, order, sincePost, limit, sortSign)
	}

	if err != nil {
		return nil, err
	}

	posts := make(models.Posts, 0)
	for sortedRows.Next() {
		p := models.Post{}
		err = sortedRows.Scan(
				&p.Id,
				&p.Parent,
				&p.Thread,
				&p.Forum,
				&p.Author,
				&p.Message,
				&p.Created,
				&p.IsEdited,
			)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err = sortedRows.Close(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (t *ThreadRepository) getFlatSortedRows(id int32, order, sincePost, limit, sortSign string) (sortedRows *sql.Rows, err error) {
	query := "SELECT id, parent, thread, forum, author, message, created, is_edited FROM posts WHERE thread = $1 "
	if sincePost != "" {
		query += fmt.Sprintf(" AND id %s %s", sortSign, sincePost)
	}
	query += fmt.Sprintf("ORDER BY created %s, id %s LIMIT %s", order, order, limit)

	sortedRows, err = t.DB.Query(query, id)
	return sortedRows, err
}

func (t *ThreadRepository) getTreeSortedRows(id int32, order, sincePost, limit, sortSign string) (sortedRows *sql.Rows, err error) {
	query := "SELECT id, parent, thread, forum, author, message, created, is_edited FROM posts WHERE thread = $1 "
	if sincePost != "" {
		query += fmt.Sprintf(" AND path %s (SELECT path FROM posts WHERE id = %s)", sortSign, sincePost)
	}
	query += fmt.Sprintf("ORDER BY path[1] %s, path %s LIMIT %s", order, order, limit)

	sortedRows, err = t.DB.Query(query, id)
	return sortedRows, err
}

func (t *ThreadRepository) getParentTreeSortedRows(id int32, order, sincePost, limit, sortSign string) (sortedRows *sql.Rows, err error) {
	query := "SELECT id, parent, thread, forum, author, message, created, is_edited FROM posts WHERE thread = $1 " +
		"AND path[1] IN (SELECT path[1] FROM posts WHERE thread = $1 AND array_length(path, 1) = 1 "
	if sincePost != "" {
		query += fmt.Sprintf(" AND path[1] %s (SELECT path FROM posts WHERE id = %s) ", sortSign, sincePost)
	}
	query += fmt.Sprintf(" ORDER BY path[1] %s, path %s LIMIT %s) ", order, order, limit)
	query += fmt.Sprintf(" ORDER BY path[1] %s, path ", order)

	sortedRows, err = t.DB.Query(query, id)
	return sortedRows, err
}

func extractParams(params url.Values) (limit string, order string, sort string, sincePost string, sortSign string) {
	limit = "500"
	if _, ok := params["limit"]; ok {
		limit = params["limit"][0]
	}

	sincePost = ""
	if _, ok := params["since"]; ok {
		sincePost = params["since"][0]
	}

	order = "ASC"
	if _, ok := params["desc"]; ok {
		order = "DESC"
	}

	if order == "ASC" {
		sortSign = ">"
	} else {
		sortSign = "<"
	}

	sort = "flat"
	if _, ok := params["sort"]; ok {
		sort = params["sort"][0]
	}

	return
}
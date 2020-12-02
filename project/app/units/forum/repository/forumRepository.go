package repository

import (
	"database/sql"
	"fmt"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"net/url"
)

type IForumRepository interface {
	Create(forum *models.Forum) error
	Get(slug string) (*models.Forum, error)
	GetUsers(forum *models.Forum, params url.Values) (models.Users, error)
	GetThreads(forum *models.Forum, params url.Values) (models.Threads, error)
}

type ForumRepository struct {
	DB *sql.DB
}

func (r ForumRepository) Create(forum *models.Forum) error {
	return r.DB.QueryRow(
		"INSERT INTO forums (slug, title, nickname) " +
			"VALUES ($1, $2, $3) RETURNING slug, title, nickname, posts, threads",
		forum.Slug,
		forum.Title,
		forum.User,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)
}

func (r ForumRepository) Get(slug string) (*models.Forum, error) {
	forum := &models.Forum{}

	err := r.DB.QueryRow("SELECT id, slug, title, nickname, posts, threads FROM forums WHERE slug = $1",
			slug,
		).Scan(
			&forum.Id,
			&forum.Slug,
			&forum.Title,
			&forum.User,
			&forum.Posts,
			&forum.Threads,
		)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (r ForumRepository) GetThreads(forum *models.Forum, params url.Values) (models.Threads, error) {
	limit, order, sort, since := extractParams(params)

	threads := models.Threads{}

	format := "SELECT id, forum, slug, author, title, message, votes, created FROM threads " +
		"WHERE forum = $1 AND created  %s '%s' ORDER BY created %s LIMIT %s"

	query := fmt.Sprintf(
		format,
		sort, since, order, limit)

	rows, err := r.DB.Query(query, forum.Slug)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		thread := models.Thread{}
		err := rows.Scan(
			&thread.Id,
			&thread.Forum,
			&thread.Author,
			&thread.Slug,
			&thread.Created,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
		)

		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (r ForumRepository) GetUsers(forum *models.Forum, params url.Values) (models.Users, error) {
	limit, order, sort, since := extractParams(params)

	users := models.Users{}

	format := "SELECT u.id, u.nickname, u.fullname, u.email, u.about FROM forum_users fu" +
		"INNER JOIN users u ON fu.nickname = u.nickname " +
		"WHERE fu.forum_id = $1 AND fu.nickname  %s '%s' " +
		"ORDER BY fu.nickname %s LIMIT %s"

	query := fmt.Sprintf(
		format,
		sort, since, order, limit)

	rows, err := r.DB.Query(query, forum.Id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.Id,
			&user.Nickname,
			&user.FullName,
			&user.Email,
			&user.About,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil

}

func extractParams(params url.Values) (limit string, order string, sort string, since string) {
	limit = "500"
	if _, ok := params["limit"]; ok {
		limit = params["limit"][0]
	}

	since = "1990-03-01 00:00:00-06"
	if _, ok := params["since"]; ok {
		since = params["since"][0]
	}

	order = "ASC"
	sort = ">="
	if _, ok := params["desc"]; ok {
		order = "DESC"
		sort = "<="
		since = "2050-05-01 00:00:00-06"
	}

	return
}
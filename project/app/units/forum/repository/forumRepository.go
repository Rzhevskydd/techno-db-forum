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
	CreateForumUser(forumId, userId int64) error
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
	limit, order, sort, since := extractThreadsQueryParams(params)

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
			&thread.Slug,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
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

// TODO forum_users обновлять при новом треде/посте
func (r ForumRepository) GetUsers(forum *models.Forum, params url.Values) (models.Users, error) {
	limit, order, signSort, since := extractUsersQueryParams(params)

	users := models.Users{}

	//format := "SELECT id, nickname, fullname, email, about FROM users  " +
	//	" WHERE id IN (SELECT user_id FROM forum_users WHERE forum_id = $1) "
	format := "SELECT u.id, u.nickname, u.fullname, u.email, u.about FROM forum_users fu " +
		" INNER JOIN users u ON fu.nickname = u.nickname WHERE forum = $1 "
		//fmt.Sprintf(" AND u.nickname %s $2 ", signSort) +
		//" ORDER BY nickname %s LIMIT $3 "
	if since != "" {
		format += fmt.Sprintf(" AND u.nickname  %s '%s' ", signSort, since)
		//format += fmt.Sprintf(" AND fu.user_id  %s '%s' ", signSort, since)
	}
	format += " ORDER BY nickname %s LIMIT $2 "
	//query += " ORDER BY fu.nickname %s LIMIT $3 "

	query := fmt.Sprintf(format, order)

	rows, err := r.DB.Query(query, forum.Slug, limit)

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

func extractThreadsQueryParams(params url.Values) (limit string, order string, sort string, since string) {
	limit = "500"
	if _, ok := params["limit"]; ok {
		limit = params["limit"][0]
	}

	since = "1990-03-01 00:00:00-06"

	order = "ASC"
	sort = ">="
	if _, ok := params["desc"]; ok && params["desc"][0] == "true"{
		order = "DESC"
		sort = "<="
		since = "2050-05-01 00:00:00-06"
	}

	if _, ok := params["since"]; ok && params["since"][0] != "" {
		since = params["since"][0]
	}

	return
}

func extractUsersQueryParams(params url.Values) (limit string, order string, signSort string, since string) {
	limit = "500"
	if _, ok := params["limit"]; ok {
		limit = params["limit"][0]
	}

	order = "ASC"
	signSort = ">"
	if _, ok := params["desc"]; ok && params["desc"][0] == "true" {
		order = "DESC"
		signSort = "<"
	}

	if _, ok := params["since"]; ok && params["since"][0] != ""{
		since = params["since"][0]
	}

	return
}

func (r ForumRepository) CreateForumUser(forumSlug, author string) error {
	_, err := r.DB.Exec("INSERT INTO forum_users(forum, nickname) VALUES ($1, $2)", forumSlug, author)
	return err
}

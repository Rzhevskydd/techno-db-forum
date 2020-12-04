package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"net/url"
	"time"
)

type IPostRepository interface {
	Create(thread *models.Thread, posts models.Posts) (models.Posts, error)
	Get(slug string) (*models.Forum, error)
	GetUsers(forum *models.Forum, params url.Values) (models.Users, error)
	GetThreads(forum *models.Forum, params url.Values) (models.Threads, error)
}

type PostRepository struct {
	DB *sql.DB
}

func (p *PostRepository) setParentId(thread *models.Thread, id int64) error {
	err := p.DB.QueryRow(
			"SELECT id FROM posts WHERE thread = $1 AND id = $2",
			thread.Id,
			id,
		).Scan(
			&id,
		)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostRepository) setAuthor(author string) error {
	return p.DB.QueryRow(
			"SELECT nickname FROM users WHERE nickname = $1",
			author,
		).Scan(
			&author,
		)
}

func (p *PostRepository) Create(thread *models.Thread, posts models.Posts) (models.Posts, error) {
	var err error
	for _, post := range posts {
		post.Thread = thread.Id
		post.Forum = thread.Forum

		if post.Parent != 0 {
			err = p.setParentId(thread, post.Parent)
		}

		if err != nil {
			return nil, errors.New("404")  // 404 | 500
		}

		err = p.setAuthor(post.Author)
		if err != nil {
			return nil, errors.New("404")  // 404
		}
	}

	now := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	for _, post := range posts {
		post.Created = now
		err = p.DB.QueryRow(
				"INSERT INTO posts (parent, thread, forum, author, created, message, path) " +
					"VALUES ($1, $2, $3, $4, $5, $6, " +
					"(SELECT path FROM posts WHERE id = $1) || " +
					"currval(pg_get_serial_sequence('posts', 'id'))::bigint) " +
					"RETURNING id ",
					post.Parent,
					post.Thread,
					post.Forum,
					post.Author,
					post.Created,
					post.Message,
			).Scan(
				&post.Id,
			)

		if err != nil {
			return nil, err // 500
		}

		// there's trigger on_new_thread_inserted
		//_, _ = p.DB.Exec("INSERT INTO forum_users(forum, nickname) VALUES($1, $2)",
		//	post.Forum,
		//	post.Author,
		//)
	}

	_, err = p.DB.Exec("UPDATE forums SET posts = posts + $1 WHERE slug = $2",
		len(posts),
		thread.Forum,
	)

	if err != nil {
		return nil, errors.New("500") // 500
	}

	return posts, nil
}

//func (r PostRepository) Get(slug string) (*models.Forum, error) {
//	forum := &models.Forum{}
//
//	err := r.DB.QueryRow("SELECT id, slug, title, nickname, posts, threads FROM forums WHERE slug = $1",
//			slug,
//		).Scan(
//			&forum.Id,
//			&forum.Slug,
//			&forum.Title,
//			&forum.User,
//			&forum.Posts,
//			&forum.Threads,
//		)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return forum, nil
//}
//
//func (r PostRepository) GetThreads(forum *models.Forum, params url.Values) (models.Threads, error) {
//	limit, order, sort, since := extractParams(params)
//
//	threads := models.Threads{}
//
//	format := "SELECT id, forum, slug, author, title, message, votes, created FROM threads " +
//		"WHERE forum = $1 AND created  %s '%s' ORDER BY created %s LIMIT %s"
//
//	query := fmt.Sprintf(
//		format,
//		sort, since, order, limit)
//
//	rows, err := r.DB.Query(query, forum.Slug)
//
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		thread := models.Thread{}
//		err := rows.Scan(
//			&thread.Id,
//			&thread.Forum,
//			&thread.Author,
//			&thread.Slug,
//			&thread.Created,
//			&thread.Title,
//			&thread.Message,
//			&thread.Votes,
//		)
//
//		if err != nil {
//			return nil, err
//		}
//
//		threads = append(threads, thread)
//	}
//
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//
//	return threads, nil
//}
//
//func (r PostRepository) GetUsers(forum *models.Forum, params url.Values) (models.Users, error) {
//	limit, order, sort, since := extractParams(params)
//
//	users := models.Users{}
//
//	format := "SELECT u.id, u.nickname, u.fullname, u.email, u.about FROM forum_users fu" +
//		"INNER JOIN users u ON fu.nickname = u.nickname " +
//		"WHERE fu.forum_id = $1 AND fu.nickname  %s '%s' " +
//		"ORDER BY fu.nickname %s LIMIT %s"
//
//	query := fmt.Sprintf(
//		format,
//		sort, since, order, limit)
//
//	rows, err := r.DB.Query(query, forum.Id)
//
//	if err != nil {
//		return nil, err
//	}
//
//	for rows.Next() {
//		user := models.User{}
//		err := rows.Scan(
//			&user.Id,
//			&user.Nickname,
//			&user.FullName,
//			&user.Email,
//			&user.About,
//		)
//
//		if err != nil {
//			return nil, err
//		}
//
//		users = append(users, user)
//	}
//
//	if err := rows.Close(); err != nil {
//		return nil, err
//	}
//
//	return users, nil
//
//}

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
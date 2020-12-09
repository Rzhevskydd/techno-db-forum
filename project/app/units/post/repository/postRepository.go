package repository

import (
	"database/sql"
	"errors"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"net/url"
	"strconv"
	"time"
)

type IPostRepository interface {
	Create(thread *models.Thread, posts models.Posts) (models.Posts, error)
	GetWithRelated(id string, related models.PostRelatedFields) (*models.PostFull, error)
	GetPost(idStr string) (*models.Post, error)
	Update(idStr, message string) (*models.Post, error)
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

	now := time.Now()
	for _, post := range posts {
		//post.Created = now
		err = p.DB.QueryRow(
				"INSERT INTO posts (parent, thread, forum, author, created, message, path) " +
					"VALUES ($1, $2, $3, $4, $5, $6, " +
					"(SELECT path FROM posts WHERE id = $1) || " +
					"currval(pg_get_serial_sequence('posts', 'id'))::bigint) " +
					"RETURNING id, created ",
					post.Parent,
					post.Thread,
					post.Forum,
					post.Author,
					now,
					post.Message,
			).Scan(
				&post.Id,
				&post.Created,
			)

		if err != nil {
			return nil, err // 500
		}

		// there's trigger on_new_thread_inserted
		_, _ = p.DB.Exec("INSERT INTO forum_users(forum, nickname) VALUES($1, $2)",
			post.Forum,
			post.Author,
		)
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

func (p *PostRepository) GetWithRelated(idStr string, related models.PostRelatedFields) (*models.PostFull, error) {
	postFull := &models.PostFull{}
	postFull.Post = &models.Post{}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	err = p.DB.QueryRow(
			"SELECT id, parent, thread, forum, author, created, is_edited, message FROM posts WHERE id = $1",
			id,
		).Scan(
			&postFull.Post.Id,
			&postFull.Post.Parent,
			&postFull.Post.Thread,
			&postFull.Post.Forum,
			&postFull.Post.Author,
			&postFull.Post.Created,
			&postFull.Post.IsEdited,
			&postFull.Post.Message,
		)

	if err != nil {
		return nil, err
	}

	if related.WithUser {
		postFull.Author = &models.User{}

		err = p.DB.QueryRow(
			"SELECT id, nickname, email, fullname, about FROM users WHERE nickname = $1",
			postFull.Post.Author,
		).Scan(
			&postFull.Author.Id,
			&postFull.Author.Nickname,
			&postFull.Author.Email,
			&postFull.Author.FullName,
			&postFull.Author.About,
		)

		if err != nil {
			return nil, err
		}
	}

	if related.WithForum {
		postFull.Forum = &models.Forum{}

		err = p.DB.QueryRow(
			"SELECT id, slug, title, nickname, posts, threads FROM forums WHERE slug = $1",
			postFull.Post.Forum,
		).Scan(
			&postFull.Forum.Id,
			&postFull.Forum.Slug,
			&postFull.Forum.Title,
			&postFull.Forum.User,
			&postFull.Forum.Posts,
			&postFull.Forum.Threads,
		)

		if err != nil {
			return nil, err
		}
	}

	if related.WithThread {
		postFull.Thread = &models.Thread{}

		err = p.DB.QueryRow(
			"SELECT id, forum, author, created, message, title, votes, slug FROM threads WHERE id = $1",
			postFull.Post.Thread,
		).Scan(
			&postFull.Thread.Id,
			&postFull.Thread.Forum,
			&postFull.Thread.Author,
			&postFull.Thread.Created,
			&postFull.Thread.Message,
			&postFull.Thread.Title,
			&postFull.Thread.Votes,
			&postFull.Thread.Slug,
		)

		if err != nil {
			return nil, err
		}
	}

	return postFull, nil
}

func (p *PostRepository) GetPost(idStr string) (*models.Post, error) {
	post := &models.Post{}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	err = p.DB.QueryRow(
		"SELECT id, parent, thread, forum, author, created, is_edited, message FROM posts WHERE id = $1",
		id,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Thread,
		&post.Forum,
		&post.Author,
		&post.Created,
		&post.IsEdited,
		&post.Message,
	)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p *PostRepository) Update(idStr, message string) (*models.Post, error) {
	post := &models.Post{}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	if err = p.DB.QueryRow(
		"UPDATE posts SET message = $2, is_edited = TRUE WHERE id = $1 " +
			"RETURNING id, parent, thread, forum, author, created, is_edited, message",
		id,
		message,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Thread,
		&post.Forum,
		&post.Author,
		&post.Created,
		&post.IsEdited,
		&post.Message,
	); err != nil {
		return nil, err
	}

	return post, nil
}

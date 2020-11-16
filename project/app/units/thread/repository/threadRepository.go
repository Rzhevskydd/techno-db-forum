package threadRepository

import (
	"database/sql"
	"errors"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
)

type IThreadRepository interface {
	Create(thread *models.Thread) (*models.Thread, error)
	GetPosts(thread *models.Thread) (*models.Posts, error)
	//Update(nickname string, email string) (models.Users, error)
	//Update(user *models.User) (*models.User, error)
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

	//f := &models.Forum{}
	_, err = t.DB.Exec("UPDATE forums SET threads = threads + 1 WHERE slug = $1",
			thread.Forum,
		)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

//func (r *ThreadRepository) Get(nickname string) (*models.User, error) {
//	user := &models.User{}
//	 err := r.DB.QueryRow("SELECT nickname, email, about, fullname " +
//		"FROM users WHERE LOWER(nickname) = LOWER($1)",
//		nickname,
//	 ).Scan(
//			&user.Nickname,
//			&user.Email,
//			&user.About,
//			&user.FullName,
//	)
//	 if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, nil
//		} else {
//			return nil, err
//		}
//	 }
//
//	return user, nil
//}
//
//func (t *ThreadRepository) GetAll(nickname string, email string) (models.Users, error) {
//	var users models.Users
//
//	rows, err := r.DB.Query(
//		"SELECT nickname, email, about, fullname FROM users " +
//				"WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2)",
//		nickname,
//		email)
//
//	if err != nil {
//		return users, err
//	}
//
//	for rows.Next() {
//		user := models.User{}
//		err = rows.Scan(&user.Nickname, &user.Email, &user.About, &user.FullName)
//		if err != nil {
//			return users, err
//		}
//		users = append(users, user)
//	}
//
//	if err = rows.Close(); err != nil {
//		return nil, err
//	}
//
//	return users, nil
//}
//
//func (t *ThreadRepository) Update(user *models.User) error {
//	_, err := r.DB.Exec(
//		"UPDATE users SET email = $1, about = $2, fullname = $3" +
//			" WHERE nickname = $4",
//		user.Email,
//		user.About,
//		user.FullName,
//		user.Nickname,
//	)
//	return err
//}

package userRepository

import (
	"database/sql"
	"errors"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
)

type IUserRepository interface {
	Create(forum *models.User) error
	Get(nickname string) (*models.User, error)
	GetAll(nickname string, email string) (models.Users, error)
	Update(user *models.User) (*models.User, error)
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
	 _, err := r.DB.Exec(
			"INSERT INTO users (nickname, email, about, fullname)" +
				"VALUES ($1, $2, $3, $4)",
				user.Nickname,
				user.Email,
				user.About,
				user.FullName,
		)
	return err
}

func (r *UserRepository) Get(nickname string) (*models.User, error) {
	user := &models.User{}
	 err := r.DB.QueryRow("SELECT id, nickname, email, about, fullname " +
		"FROM users WHERE LOWER(nickname) = LOWER($1)",
		nickname,
	 ).Scan(
	 		&user.Id,
			&user.Nickname,
			&user.Email,
			&user.About,
			&user.FullName,
	)
	 if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	 }

	return user, nil
}

func (r *UserRepository) GetAll(nickname string, email string) (models.Users, error) {
	var users models.Users

	rows, err := r.DB.Query(
		"SELECT nickname, email, about, fullname FROM users " +
				"WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2)",
		nickname,
		email)
	
	if err != nil {
		return users, err
	}
	
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Email, &user.About, &user.FullName)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}
	
	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.DB.Exec(
		"UPDATE users SET email = $1, about = $2, fullname = $3" +
			" WHERE LOWER(nickname) = LOWER($4)",
		user.Email,
		user.About,
		user.FullName,
		user.Nickname,
	)
	return err
}

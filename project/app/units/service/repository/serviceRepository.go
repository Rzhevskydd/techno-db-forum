package serviceRepository

import (
	"database/sql"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
)

type IServiceRepository interface {
	Clear() error
	GetStatus() (*models.Status, error)
}

type ServiceRepository struct {
	DB *sql.DB
}

func (s *ServiceRepository) Clear() error {
	if _, err := s.DB.Exec("TRUNCATE votes, users, posts, threads, forums RESTART IDENTITY CASCADE"); err != nil {
		return err
	}

	return nil
}

func (s *ServiceRepository) GetStatus() (*models.Status, error) {
	status := &models.Status{}
	var err error

	err = s.DB.QueryRow(
		`SELECT COUNT(*) 
		 FROM users`,
	).Scan(&status.User)
	if err != nil {
		return nil, err
	}

	err = s.DB.QueryRow(
		`SELECT COUNT(*) 
		 FROM forums`,
	).Scan(&status.Forum)
	if err != nil {
		return nil, err
	}

	err = s.DB.QueryRow(
		`SELECT COUNT(*) 
		 FROM threads`,
	).Scan(&status.Thread)
	if err != nil {
		return nil, err
	}

	err = s.DB.QueryRow(
		`SELECT COUNT(*) 
		 FROM posts`,
	).Scan(&status.Post)
	if err != nil {
		return nil, err
	}


	return status, nil
}
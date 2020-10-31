package forumRepository

import (
	"database/sql"
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


}
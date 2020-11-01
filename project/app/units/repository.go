package units

import (
	"database/sql"
	f "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/forumRepository"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/repository"
)

type Repositories struct {
	User  u.UserRepository
	Forum f.ForumRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:  u.UserRepository{DB: db},
		Forum: f.ForumRepository{DB: db},
	}
}

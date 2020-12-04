package units

import (
	"database/sql"
	f "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/repository"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/repository"
	t "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/repository"
	p "github.com/Rzhevskydd/techno-db-forum/project/app/units/post/repository"

)

type Repositories struct {
	User  u.UserRepository
	Forum f.ForumRepository
	Thread t.ThreadRepository
	Post p.PostRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:  u.UserRepository{DB: db},
		Forum: f.ForumRepository{DB: db},
		Thread: t.ThreadRepository{DB: db},
		Post: p.PostRepository{DB: db},
	}
}

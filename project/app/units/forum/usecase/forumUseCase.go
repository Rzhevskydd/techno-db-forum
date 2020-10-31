package forumUsecase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/app"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"net/url"
	"runtime"
)

type IForumUseCase interface {
	CreateForum(forum *models.Forum) error
	GetForum(slug string) (*models.Forum, error)
	GetForumUsers(forum *models.Forum, params url.Values) (models.Users, error)
	GetForumThreads(forum *models.Forum, params url.Values) (models.Threads, error)
}

type ForumUseCase struct {
	Repos app.Repositories
}

func (f *ForumUseCase) CreateForum(forum *models.Forum) error {

	//return runtime.RuntimeError()
}



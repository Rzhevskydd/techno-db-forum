package threadUseCase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	p "github.com/Rzhevskydd/techno-db-forum/project/app/units/post/repository"
	t "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/repository"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/repository"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/validator"
	"net/url"
)

type IThreadUseCase interface {
	CreateThreadPosts(slugOrId string, posts models.Posts) (models.Posts, int, error)
	GetThread(slugOrId string) (*models.Thread, error)
	UpdateThread(slugOrId string, updateThread *models.ThreadUpdate) (*models.Thread, error)
	GetThreadPosts(slugOrId string, params url.Values) (models.Posts, error)
	VoteThread(slugOrId string, vote *models.Vote) (*models.Thread, error)
}

type ThreadUseCase struct {
	UserRep u.UserRepository
	ThreadRep t.ThreadRepository
	PostRep p.PostRepository
}


// TODO forum_users обновлять при новом треде/посте
func (t *ThreadUseCase) CreateThreadPosts(slugOrId string, posts models.Posts) (models.Posts, int, error) {
	thread, err := t.ThreadRep.Get(slugOrId)
	if err != nil {
		return nil, 404, err
	}

	//if err != nil {
	//	return nil, 500
	//}

	newPosts, err := t.PostRep.Create(thread, posts)
	if err != nil {
		switch err.Error() {
		case "404":
			return nil, 404, err

		default:
			return nil, 409, err
		}


	}

	return newPosts, 201, nil
}

func (t *ThreadUseCase) GetThread(slugOrId string) (*models.Thread, error) {
	thread, err := t.ThreadRep.Get(slugOrId)
	if err != nil || thread == nil {
		return nil, err
	}

	return thread, nil
}

func (t *ThreadUseCase) UpdateThread(slugOrId string, updateThread *models.ThreadUpdate) (*models.Thread, error) {
	thread, err := t.ThreadRep.Get(slugOrId)
	if err != nil {
		return nil, err
	}

	if validator.IsEmpty(updateThread.Title + updateThread.Message) {
		return thread, nil
	}

	if validator.IsEmpty(updateThread.Title) {
		updateThread.Title = thread.Title
	}

	if validator.IsEmpty(updateThread.Message) {
		updateThread.Message = thread.Message
	}

	thread, err = t.ThreadRep.Update(slugOrId, updateThread)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (t *ThreadUseCase) VoteThread(slugOrId string, vote *models.Vote) (*models.Thread, error) {
	thread, err := t.ThreadRep.Get(slugOrId)
	if err != nil {
		return nil, err
	}

	thread, err = t.ThreadRep.Vote(thread, vote)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (t *ThreadUseCase) GetThreadPosts(slugOrId string, params url.Values) (models.Posts, error) {
	thread, err := t.ThreadRep.Get(slugOrId)
	if err != nil {
		return nil, err
	}

	posts, err := t.ThreadRep.GetPosts(thread, params)
	if err != nil {
		return nil, err
	}

	return posts, nil
}


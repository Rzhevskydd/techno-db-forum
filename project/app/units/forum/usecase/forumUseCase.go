package forumUsecase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	f "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/repository"
	t "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/repository"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/repository"
	"net/url"
)

type IForumUseCase interface {
	CreateForum(forum *models.Forum) (*models.Forum, int)
	GetForum(slug string) (*models.Forum, error)
	GetForumUsers(slug string, params url.Values) (models.Users, error)
	GetForumThreads(slug string, params url.Values) (models.Threads, error)
	CreateForumThread(thread *models.Thread, slug string) (*models.Thread, int, error)
}

type ForumUseCase struct {
	ForumRep f.ForumRepository
	UserRep u.UserRepository
	ThreadRep t.ThreadRepository

}

func (f *ForumUseCase) CreateForum(forum *models.Forum) (*models.Forum, int) {
	forumUser, err := f.UserRep.Get(forum.User)

	if err != nil {
		return nil, 500
	}

	if forumUser == nil {
		//forum.User = forumUser.Nickname
		return forum, 404
	}

	forum.User = forumUser.Nickname

	if err = f.ForumRep.Create(forum); err != nil {
		existingForum, err := f.ForumRep.Get(forum.Slug)
		if err != nil {
			return nil, 409
		}

		return existingForum, 409
	}

	return forum, 201
}

func (f *ForumUseCase) CreateForumThread(thread *models.Thread, slug string) (*models.Thread, int) {
	threadUser, err := f.UserRep.Get(thread.Author)
	if threadUser == nil || err != nil {
		//thread.Author = threadUser.Nickname
		return nil, 404
	}

	threadForum, err := f.GetForum(slug)
	if threadForum == nil || err != nil {
		return nil, 404
	}

	thread.Forum = threadForum.Slug

	existingThread, err := f.ThreadRep.Get(thread.Slug)
	if existingThread != nil {
		return existingThread, 409
	}


	newThread, err := f.ThreadRep.Create(thread)
	if err != nil {
		println(err)
		return nil, 409
	}

	_ = f.ForumRep.CreateForumUser(newThread.Forum, newThread.Author)
	//if err != nil {
	//	println(err.Error())
	//	return nil, 500
	//}

	return newThread, 201
}

func (f *ForumUseCase) GetForum(slug string) (*models.Forum, error) {
	forum, err := f.ForumRep.Get(slug)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (f *ForumUseCase) GetForumThreads(slug string, params url.Values) (models.Threads, error) {
	forum, err := f.ForumRep.Get(slug)
	if err != nil || forum == nil {
		return nil, err
	}

	threads, err := f.ForumRep.GetThreads(forum, params)
	if err != nil {
		return nil, err
	}

	return threads, nil
}


func (f *ForumUseCase) GetForumUsers(slug string, params url.Values) (models.Users, error) {
	forum, err := f.ForumRep.Get(slug)
	if err != nil || forum == nil {
		return nil, err
	}

	users, err := f.ForumRep.GetUsers(forum, params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

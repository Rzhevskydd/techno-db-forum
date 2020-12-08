package postUsecase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	p "github.com/Rzhevskydd/techno-db-forum/project/app/units/post/repository"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/validator"
	"net/url"
	"strings"
)

type IPostUseCase interface {
	CreatePosts(posts models.Posts) (models.Posts, int)
	GetPostDetails(id string, params url.Values) (*models.PostFull, error)
	UpdatePostDetails(id string, message string) (*models.Post, error)
	//GetForum(slug string) (*models.Forum, error)
	//GetForumUsers(slug string, params url.Values) (models.Users, error)
	//GetForumThreads(slug string, params url.Values) (models.Threads, error)
	//CreateForumThread(thread *models.Thread, slug string) (*models.Thread, int, error)
}

type PostUseCase struct {
	PostRep p.PostRepository
}

func getRelatedFields(params url.Values) models.PostRelatedFields {
	relatedFields := models.PostRelatedFields{
		WithUser: false,
		WithForum: false,
		WithThread: false,
	}

	if _, ok := params["related"]; ok {
		in := params["related"][0]
		//split := strings.Split(in, ",")

		if strings.Contains(in, "user") {
			relatedFields.WithUser = true
		}

		if strings.Contains(in, "forum") {
			relatedFields.WithForum = true
		}

		if strings.Contains(in, "thread") {
			relatedFields.WithThread = true
		}
	}
	return relatedFields
}

func (p *PostUseCase) GetPostDetails(id string, params url.Values) (*models.PostFull, error) {
	related := getRelatedFields(params)

	postFull, err := p.PostRep.GetWithRelated(id, related)
	if err != nil {
		return nil, err
	}

	return postFull, nil
}

func (p *PostUseCase) UpdatePostDetails(id string, message string) (*models.Post, error) {
	post, err := p.PostRep.GetPost(id)
	if err != nil {
		return nil, err
	}

	if validator.IsEmpty(message) || post.Message == message {
		return post, nil
	}

	post, err = p.PostRep.Update(id, message)
	if err != nil {
		return nil, err
	}

	return post, nil
}

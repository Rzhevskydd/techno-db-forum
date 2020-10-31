package userUseCase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/app"
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	"net/url"
	"runtime"
)

type IUserUseCase interface {
	CreateUser(user *models.User) (models.Users, error)
	GetUser(nickname string) (models.Users, error)
	UpdateUser(user *models.User) (*models.User, error)
}

type UserUseCase struct {
	Repos app.Repositories
}

func (u *UserUseCase) CreateUser(user *models.User) (models.Users, error) {
	users, err := u.Repos.User.GetAll(user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return users, nil
	}

	if err = u.Repos.User.Create(user); err != nil {
		return nil, err
	}

	users = append(users, *user)
	return users, nil
}

func (u *UserUseCase) GetUser(nickname string) (*models.User, error) {
	user, err := u.Repos.User.Get(nickname)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) UpdateUser(user *models.User) (*models.User, error) {


	//return runtime.RuntimeError()
}



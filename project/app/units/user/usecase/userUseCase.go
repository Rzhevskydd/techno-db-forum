package userUseCase

import (
	"github.com/Rzhevskydd/techno-db-forum/project/app/models"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/repository"
	"github.com/Rzhevskydd/techno-db-forum/project/app/utils/validator"
)

type IUserUseCase interface {
	CreateUser(user *models.User) (models.Users, error)
	GetUser(nickname string) (models.Users, error)
	UpdateUser(user *models.User) (*models.User, error)
}

type UserUseCase struct {
	UserRep u.UserRepository
}

func (u *UserUseCase) CreateUser(user *models.User) (models.Users, error) {
	users, err := u.UserRep.GetAll(user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return users, nil
	}

	if err = u.UserRep.Create(user); err != nil {
		return nil, err
	}

	users = append(users, *user)
	return nil, nil
}

func (u *UserUseCase) GetUser(nickname string) (*models.User, error) {
	return u.UserRep.Get(nickname)
}

func (u *UserUseCase) UpdateUser(user *models.User) (*models.User, error) {
	dbUser, err := u.UserRep.Get(user.Nickname)

	if err != nil {
		return nil, err
	}

	if dbUser == nil {
		return dbUser, err
	}

	// TODO валидаторы (через регулярки)

	if validator.IsEmpty(user.Email) {
		user.Email = dbUser.Email
	}

	if validator.IsEmpty(user.FullName) {
		user.FullName = dbUser.FullName
	}

	if validator.IsEmpty(user.About) {
		user.About = dbUser.About
	}

	if err = u.UserRep.Update(user); err != nil {
		return dbUser, err
	}

	return user, nil
}



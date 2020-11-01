package units

import (
	f "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/usecase"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/usecase"
)

type UseCase struct {
	User  u.UserUseCase
	Forum f.ForumUseCase
}

func NewUseCase(repos *Repositories) *UseCase {
	return &UseCase{
		User:  u.UserUseCase{UserRep: repos.User},
		Forum: f.ForumUseCase{UserRep: repos.User},
	}
}
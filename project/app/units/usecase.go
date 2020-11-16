package units

import (
	f "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/usecase"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/usecase"
	t "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/usecase"
)

type UseCase struct {
	User  u.UserUseCase
	Forum f.ForumUseCase
	Thread t.ThreadUseCase
}

func NewUseCase(repos *Repositories) *UseCase {
	return &UseCase{
		User:  u.UserUseCase{UserRep: repos.User},
		Forum: f.ForumUseCase{
			ForumRep: repos.Forum,
			UserRep: repos.User,
			ThreadRep: repos.Thread,
		},
	}
}
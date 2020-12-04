package units

import (
	f "github.com/Rzhevskydd/techno-db-forum/project/app/units/forum/usecase"
	t "github.com/Rzhevskydd/techno-db-forum/project/app/units/thread/usecase"
	u "github.com/Rzhevskydd/techno-db-forum/project/app/units/user/usecase"
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
		Thread: t.ThreadUseCase{
			UserRep:   repos.User,
			ThreadRep: repos.Thread,
			PostRep:   repos.Post,
		},
	}
}
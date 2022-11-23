package todo

import "testfiber/pkg/utility"

type Service struct {
	Repo Repository
	Sess Session
	ID   *utility.ID
}

func NewService(r Repository, s Session) Service {
	return Service{
		Repo: r,
		Sess: s,
		ID:   utility.NewID(),
	}
}

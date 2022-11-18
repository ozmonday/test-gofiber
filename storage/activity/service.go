package activity

type Service struct {
	Repo Repository
	Sess Session
}

func NewService(r Repository, s Session) Service {
	return Service{
		Repo: r,
		Sess: s,
	}
}

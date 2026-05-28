package users

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name, email string) User {
	return s.repo.Create(name, email)
}

func (s *Service) Update(id int64, name, email string) (User, error) {
	return s.repo.Update(id, name, email)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) List() []User {
	return s.repo.List()
}

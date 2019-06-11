package service

import (
	"gitlab.com/standard-go/project/internal/app/project"
)

func (s *Service) FetchIndexUser(pageRequest *project.PageRequest) ([]*project.User, int, error) {
	return s.Repo.FetchIndexUser(pageRequest)
}

func (s *Service) FetchShowUser(id string) (*project.User, error) {
	return s.Repo.FetchShowUser(id)
}

func (s *Service) FetchStoreUser(data *project.User) (*project.User, error) {
	return s.Repo.FetchStoreUser(data)
}

func (s *Service) FetchUpdateUser(id string, data *project.User) (*project.User, error) {
	return s.Repo.FetchUpdateUser(id, data)
}

func (s *Service) FetchDestroyUser(id string) error {
	return s.Repo.FetchDestroyUser(id)
}

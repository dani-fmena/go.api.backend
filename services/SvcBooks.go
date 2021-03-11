package services

import (
	"go.api.backend/data"
	"go.api.backend/repo/db"
)


type svcBook struct {
	*svcBase
}

// NewSvcBooks handles some of the CRUD operations for books schemes / entities.
// It depends on a movie repository for its actions.
// The code here decouple the data login from the higher level components.
// As a result, different repositories type can be used with the same logic without any additional changes.
func NewSvcBooks(repo *db.RepoDB) Service  {
	return &svcBook{NewService(repo)}								// We are using struct composition here. Hence the anonymous field (https://golangbot.com/inheritance/)
}

func (s *svcBook) GetAll() []data.Temporal {

	(*s.repo).DbGetAll()

	return make([]data.Temporal, 0)
}

func (s *svcBook) GetByID(id uint) (data.Temporal, bool) {
	panic("implement me")
}

func (s *svcBook) DeleteByID(id uint) bool {
	panic("implement me")
}

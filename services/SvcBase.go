package services

import (
	"go.api.backend/data"
	"go.api.backend/repo/db"
)


type svcBase struct {
	repo *db.RepoDB
}

// Service describe default service interface
// It's used everywhere, we may need to change or try an experimental different domain logic at the future.
// This interface allow that
type Service interface {
	GetAll() []data.Temporal
	GetByID(id uint) (data.Temporal, bool)
	DeleteByID(id uint) bool
}

// NewService create the default service struct
func NewService(db *db.RepoDB) *svcBase  {
	return &svcBase{repo: db}
}

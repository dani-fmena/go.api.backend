package db

import (
	"github.com/go-pg/pg/v10"
	"go.api.backend/data/models"
)

type repoDBBase struct {
	Pgdb *pg.DB `*pg.DB:"Database connection object"`
}

type RepoDB interface {
	DbGetAll()  []models.Book
}

// NewDbRepo creates a new Base Database Repository instance
func NewDbRepo(e *pg.DB) *repoDBBase {
	return &repoDBBase{Pgdb: e}
}

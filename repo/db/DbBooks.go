package db

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"go.api.backend/data/models"
)

type dbBooks struct {
	*repoDBBase
}

// NewRepoDbBook creates a new Temporal Database Repository instance
func NewRepoDbBook(dbCtx *pg.DB) RepoDB {
	return &dbBooks{NewDbRepo(dbCtx)} // We are using struct composition here. Hence the anonymous field (https://golangbot.com/inheritance/)
}

func (repo *dbBooks) DbGetAll() []models.Book {
	b := models.Book{}

	res, _ := repo.Pgdb.Query(&b, "SELECT * FROM books")

	fmt.Sprintf("%s", res)


	panic("asdsad")
}

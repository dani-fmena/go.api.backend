package db

import (
	"github.com/go-pg/pg/v10"
	"go.api.backend/data/models"
)

type RepoDbBook interface {
	GetAll(list *[]models.Book) error
	GetByID(book *models.Book) error
	DelByID(Id *uint) (uint, error)
}

type dbBooks struct {
	Pgdb *pg.DB `*pg.DB:"Database connection object"`
}

// NewRepoDbBook creates a new Temporal Database Repository instance
func NewRepoDbBook(dbCtx *pg.DB) RepoDbBook {
	return &dbBooks{dbCtx}
}

// GetAll get all record for a specific entity and set the result in the referenced (pointer) list (slice).
//
// - list [*[]models.Book] ~ A pointer to a slice for storing the query result
func (r *dbBooks) GetAll(list *[]models.Book) error {

	return r.Pgdb.Model(&list).Select()
	// _, err := r.Pgdb.Query(list, "SELECT * FROM list") hard coded query sample, allow placeholder see the docs (https://pg.uptrace.dev/placeholders/)
	// return err, this is another alternative 'cause the Pgdb method return the error
}

// GetByID get an entity by Id. If no entity found then err != nil.
//
// - entity [*models.Book] ~ A pointer to the holder entity struct to be found
func (r *dbBooks) GetByID(entity *models.Book) error {
	return r.Pgdb.Model(entity).WherePK().Select() 			// I'm not using & 'cause the param is already a pointer
}

// DelByID delete an entity by Id. If no entity found then err != nil.
// uint > 0 if any record was deleted, otherwise if 0 and no error then 404.
// - entity [*uint] ~ If of the entity to be deleted
func (r *dbBooks) DelByID(Id *uint) (uint, error)  {
	b := models.Book{Id: *Id}

	if res, err := r.Pgdb.Model(&b).WherePK().Delete(); res != nil {
		return uint(res.RowsAffected()), err
	} else {
		return 0, err
	}
}

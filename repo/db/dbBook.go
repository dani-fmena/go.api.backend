package db

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"go.api.backend/data"
	"go.api.backend/data/models"
)

type RepoDbBook interface {
	GetAll(list *[]models.Book) error
	GetByID(ent *models.Book) error
	DelByID(Id *uint) (uint, error)
	Add(ent *models.Book) error
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
// - ent [*models.Book] ~ A pointer to the holder entity struct to be found
func (r *dbBooks) GetByID(ent *models.Book) error {
	return r.Pgdb.Model(ent).WherePK().Select() 			// I'm not using & 'cause the param is already a pointer
}

// DelByID delete an entity by Id. If no entity found then err != nil.
// uint > 0 if any record was deleted, otherwise if 0 and no error then 404.
// - entity [*uint] ~ If of the entity to be deleted
func (r *dbBooks) DelByID(Id *uint) (uint, error) {
	b := models.Book{Id: *Id}

	if res, err := r.Pgdb.Model(&b).WherePK().Delete(); res != nil {
		return uint(res.RowsAffected()), err
	} else {
		return 0, err
	}
}


// Add a Book to the repository. If the book name already exist then err != nil.
// If something occurs during the ops also err != nil.
// - entity [*models.Book] ~ New book to be added to the repo
func (r *dbBooks) Add(entity *models.Book) error {

	isExist, e1 := r.Pgdb.Model(entity).Where("name = ?", entity.Name).Exists()
	if isExist && e1 == nil {
		return errors.New(data.ErrDuplicateKey)
	} else if e1 != nil {
		return e1								// Something happen
	} else {
		_, e2 := r.Pgdb.Model(entity).Insert()  // I'm not using & 'cause the param is already a pointer
		return e2
	}
}
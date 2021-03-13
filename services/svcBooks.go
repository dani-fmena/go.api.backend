package services

import (
	"go.api.backend/data/models"
	"go.api.backend/repo/db"
)

// ServiceBook is a sample for the services interface, defining its methods / functions
type ServiceBook interface {
	GetAll() ([]models.Book, error)
	GetByID(Id *uint) (models.Book, error)
	DelByID(Id *uint) (uint, error)
}

type svcBook struct {
	repo *db.RepoDbBook
}

// NewSvcBooks create the service Books that handles for the CRUD and other operations
// It depends on repository for accomplish his responsibility.
// The code here decouple the data login from the higher level components.
// As a result, different repositories type can be used with this same logic without any additional changes here.
//
// - repo [*db.RepoDbBook] ~ Repository instance
func NewSvcBooks(repo *db.RepoDbBook) ServiceBook {
	return &svcBook{repo}
}

// GetAll Get a list of all the books on the repository. If there is a error it's != from null
func (s *svcBook) GetAll() ([]models.Book, error) {
	list := make([]models.Book, 0)

	return list, (*s.repo).GetAll(&list)
}

// GetByID Get A book by its Id. If there is a error it's != from nil
//
// - id [*uint] ~ Book ID
func (s *svcBook) GetByID(Id *uint) (models.Book, error) {
	book := models.Book{Id: *Id}

	return book, (*s.repo).GetByID(&book)
}

// DelByID delete a book by its Id. If there is a error it's != from nil.
// Row affected > 0 if any record was deleted, otherwise if 0 and no error then 404.
//
// - id [*uint] ~ Book ID
func (s *svcBook) DelByID(Id *uint) (uint, error) {
	return (*s.repo).DelByID(Id)
}

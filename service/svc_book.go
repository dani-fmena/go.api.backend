package service

import (
	"go.api.backend/schema/models"
	"go.api.backend/repo/db"
)

// SvcBook is a sample for the service interface, defining its methods / functions
type SvcBook interface {
	GetAll() ([]models.Book, error)
	GetByID(Id *uint) (models.Book, error)
	DelByID(Id *uint) (uint, error)
	Create(book *models.Book) error
	UpdateBook(book *models.Book) (uint, error)
}

type svcBook struct {
	pRepo *db.RepoDbBook
}

// NewSvcBooks create the service Books that handles for the CRUD and other operations
// It depends on repository for accomplish his responsibility.
// The code here decouple the data login from the higher level components.
// As a result, different repositories type can be used with this same logic without any additional changes here.
//
// - pRepo [*db.RepoDbBook] ~ Repository instance pointer
func NewSvcBooks(pRepo *db.RepoDbBook) SvcBook {
	return &svcBook{pRepo}
}

// GetAll Get a list of all the books on the repository. If there is a error it's != from null
// Return a slice of books
func (s *svcBook) GetAll() ([]models.Book, error) {
	list := make([]models.Book, 0)

	return list, (*s.pRepo).GetAll(&list)
}

// GetByID Get A book by its Id. If there is a error it's != from nil
//
// - id [*uint] ~ Book ID pointer
func (s *svcBook) GetByID(pId *uint) (models.Book, error) {
	book := models.Book{Id: *pId}

	return book, (*s.pRepo).GetByID(&book)
}

// DelByID delete a book by its Id. If there is a error it's != from nil.
// Row affected (first return data) > 0 if any record was deleted, otherwise if 0 and no error then 404.
//
// - pId [*uint] ~ Book ID pointer
func (s *svcBook) DelByID(pId *uint) (uint, error) {
	return (*s.pRepo).DelByID(pId)
}

// Create creat a book. If there is a error it's != from nil.
// If the name key exist then a duplicated key error will be returned
//
// - pBook [*models.Book] ~ New book struct pointer to be created
func (s *svcBook) Create(pBook *models.Book) error {
	return (*s.pRepo).Add(pBook)
}

// UpdateBook update a book with the giving data
//
// - pBookDto [*models.Book] ~ Book data to be updated
func (s *svcBook) UpdateBook(pBook *models.Book) (uint, error) {
	return (*s.pRepo).Update(pBook)
}
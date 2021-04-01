package mapper

import (
	"go.api.backend/schema/dto"
	"go.api.backend/schema/models"
)

// TIP ref https://hellokoding.com/crud-restful-apis-with-go-modules-wire-gin-gorm-and-mysql/

// ToBookUpdateV map a dto.BookCreateIn to models.Book. This is the POST /create alternative
func ToBookCreateV(dto *dto.BookCreateIn) *models.Book {
	return &models.Book{Name: dto.Name, Items: dto.Items}
}

// ToBookUpdateV map a dto.BookUpdateIn to models.Book. This is the PUT / update alternative
func ToBookUpdateV(dto *dto.BookUpdateIn) *models.Book {
	return &models.Book{Id: dto.Id, Name: dto.Name, Items: dto.Items}
}

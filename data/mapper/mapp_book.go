package mapper

import (
	"go.api.backend/data/dto"
	"go.api.backend/data/models"
)

func ToBook(dto *dto.BookUpdateIn) *models.Book {
	return &models.Book{Id: dto.Id, Name: dto.Name, Items: dto.Items}
}
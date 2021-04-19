package mapper

import (
	"go.api.backend/schema/dto"
	"go.api.backend/schema/models"
	"strings"
)

// TIP ref https://hellokoding.com/crud-restful-apis-with-go-modules-wire-gin-gorm-and-mysql/
// Is we need it, this method can perform validation and return two values: the mapped struct and the error

// region ======== BOOKS =================================================================

// ToBookCreateV map a dto.BookCreateIn to models.Book with the necessary data to create a new one. This is the POST /create alternative
func ToBookCreateV(dto *dto.BookCreateIn) *models.Book {
	return &models.Book{Name: dto.Name, Items: dto.Items}
}

// ToBookUpdateV map a dto.BookUpdateIn to models.Book with the necessary data to make a update. This is the PUT / update alternative
func ToBookUpdateV(dto *dto.BookUpdateIn) *models.Book {
	return &models.Book{Id: dto.Id, Name: dto.Name, Items: dto.Items}
}
// endregion =============================================================================

// region ======== AUTHORIZATION =========================================================

// ToAccessTokenDataV map a dto dto.SisecAccessToken to dto.AccessTokenData
func ToAccessTokenDataV(obj *dto.SisecAccessToken) *dto.AccessTokenData {
	claims := dto.Claims{ Sub: obj.Client_Id, Rol: obj.Rol }

	return &dto.AccessTokenData{ Scope: strings.Fields(obj.Scope), Claims: claims }
}

// endregion =============================================================================

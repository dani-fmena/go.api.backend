package dto

// We can also declare Swagger / OpenAPI annotation for structs, useful ins custom struct params types.
// https://github.com/swaggo/swag#attribute | https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html#attribute
type BookUpdateIn struct {
	Id    uint   `example:"24" validate:"gte=0,numeric"`
	Name  string `example:"The Book of Eli" validate:"required,ascii,gte=3,lte=60"`
	Items uint   `example:"46" validate:"required,number,gte=0,lte=130"`
}

type BookCreateIn struct {
	Name  string `example:"The Book of Eli" validate:"required,ascii,gte=3,lte=60"`
	Items uint   `example:"46" validate:"required,number,gte=0,lte=130"`
}

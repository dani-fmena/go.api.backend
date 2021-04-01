package dto

type BookUpdateIn struct {
	Id    uint   `example:"24" validate:"gte=0,numeric"`
	Name  string `example:"The Book of Eli" validate:"required,ascii,gte=3,lte=60"`
	Items uint   `example:"46" validate:"required,number,gte=0,lte=130"`
}

type BookCreateIn struct {
	Name  string `example:"The Book of Eli" validate:"required,ascii,gte=3,lte=60"`
	Items uint   `example:"46" validate:"required,number,gte=0,lte=130"`
}

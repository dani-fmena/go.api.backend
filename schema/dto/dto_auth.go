package dto

// TODO check what happen if we pass a enum validation for the scope field ()
type UserCredIn struct {
	Username string `example:"mynickname" validate:"required,ascii,gte=3,lte=60"`
	Password string `example:"secret" validate:"required,ascii,gte=3,lte=20"`
	Scope    string `example:"web" validate:"required,ascii,gte=3,lte=10"`
}

//goland:noinspection GoSnakeCaseUsage
type SISECGrantIntentIn struct {
	Token_Type   string
	Access_Token sisecAccessToken
}

//goland:noinspection GoSnakeCaseUsage
type sisecAccessToken struct {
	Rol string
	Client_Id string
	Scope string
}

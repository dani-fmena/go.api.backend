package dto

// TODO check what happen if we pass a enum validation for the scope field ()
// This is a example declaring the validation for teh struct. It will be used when the
// struct is in the endpoint parameters
type UserCredIn struct {
	Username string `example:"mynickname" validate:"required,ascii,gte=3,lte=60"`
	Password string `example:"secret" validate:"required,ascii,gte=3,lte=20"`
	Domain   string `example:"web" validate:"required,ascii,gte=3,lte=10"`
}

//goland:noinspection GoSnakeCaseUsage
type SISECGrantIntentIn struct {
	Token_Type   string
	Access_Token SisecAccessToken			// SISEC response as access_token
}

//goland:noinspection GoSnakeCaseUsage
type SisecAccessToken struct {
	Rol string
	Client_Id string
	Scope string
}

// AccessTokenData using by this REST Api (HLF client node) to grant access to the resources
type AccessTokenData struct {
	Scope  []string
	Claims Claims
}

// Claims user claims
type Claims struct {
	Sub string
	Rol string
}

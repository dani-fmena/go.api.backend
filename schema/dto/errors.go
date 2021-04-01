package dto

// ApiError api documentation
type ApiError struct {
	Detail string `example:"Some error details"`
	Status uint   `example:"503"`
	Title  string `example:"err_code"`
}

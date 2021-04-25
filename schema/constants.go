package schema

// region ======== i18n ERROR KEYS =======================================================
const (
	ErrGeneric = "err.generic"
	ErrRepositoryOps = "err.repo_ops"
	ErrNotFound = "err.not_found"
	ErrHttpResError = "err.http_res_err"
	ErrDuplicateKey = "err.duplicate_key"

	ErrInvalidType = "err.wrong_type_assertion"
	ErrNetwork = "err.network"
	ErrJsonParse = "err.json_parse"
	ErrJwtGen = "err.jwt_generation"
	ErrWrongAuthProvider = "err.wrong_auth_provider"
	ErrUnauthorized = "err.unauthorized"
	ErrVal = "err.invalid_data"
)
// endregion =============================================================================


// region ======== ERROR DETAILS =========================================================
const (
	ErrDetNotFound        = "resource not found"
	ErrDetDuplicateKey 	  = "a unique resource field is duplicated"
	ErrDetHttpResError    = "there is an error on http request response"
	ErrDetInvalidType     = "invalid interface type (type assertion)"
	ErrDetInvalidCred     = "something was wrong with the provided user credentials"
	ErrDetInvalidProvider = "wrong or invalid provider"
)
// endregion =============================================================================


// region ======== SOME STRINGS ==========================================================
const (
	StrPgDuplicateKey = "23505" // Postgres error code for duplicate key
	StrDB404 = "no rows"
)
// endregion =============================================================================
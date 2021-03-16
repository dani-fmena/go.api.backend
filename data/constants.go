package data

// region ======== i18n ERROR KEYS =======================================================
const (
	ErrGeneric = "err.generic"
	ErrRepositoryOps = "err.repo_ops"
	ErrNotFound = "err.not_found"
	ErrDuplicateKey = "err.duplicate_key"
	ErrVal = "err.invalid_data"
)
// endregion =============================================================================


// region ======== ERROR DETAILS =========================================================
const (
	ErrDetNotFound = "Resource not found"
	ErrDetDuplicateKey = "A unique resource field is duplicated"
)
// endregion =============================================================================


// region ======== SOME STRINGS ==========================================================
const (
	StrPgDuplicateKey = "23505" // Postgres error code for duplicate key
	StrDB404 = "no rows"
)
// endregion =============================================================================
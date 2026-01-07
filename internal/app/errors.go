package app

type AppError struct {
	Code string
	Msg  string
}

func (e AppError) Error() string { return e.Msg }

func NewError(code, msg string) AppError { return AppError{Code: code, Msg: msg} }

const (
	ErrCodeInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	ErrCodeEmailExists        = "AUTH_EMAIL_EXISTS"
	ErrCodeValidation         = "VALIDATION_ERROR"
	ErrCodeInternal           = "INTERNAL_ERROR"
)

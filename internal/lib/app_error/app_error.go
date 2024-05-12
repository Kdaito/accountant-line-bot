package app_error

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (a *AppError) Error() string {
	return fmt.Sprintf("Error %d: %s - %v", a.Code, a.Message, a.Err)
}

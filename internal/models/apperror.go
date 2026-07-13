package models

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code, msg string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    msg,
		HTTPStatus: status,
		Err:        nil,
	}
}

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		HTTPStatus: e.HTTPStatus,
		Err:        err,
	}
}

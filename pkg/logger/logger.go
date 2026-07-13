package logger

type ErrorResponse struct {
}

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

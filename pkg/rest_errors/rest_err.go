package resterrors

import (
	"net/http"
)

const (
	badRequest          = "Bad Request"
	internalServerError = "Internal Server Error"
	notFound            = "Not Found"
	forbidden           = "Forbidden"
	unathorized         = "Unauthorized"
)

type RestErr struct {
	Message        string  `json:"message"`
	HttpErr        string  `json:"http_error"`
	HttpStatusCode int     `json:"status_code"`
	Errors         []error `json:"errors,omitempty"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewRestErr(message, err string, code int, errors []error) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        err,
		HttpStatusCode: code,
		Errors:         errors,
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        badRequest,
		HttpStatusCode: http.StatusBadRequest,
	}
}

func NewBadRequestValidationError(message string, errors []error) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        badRequest,
		HttpStatusCode: http.StatusBadRequest,
		Errors:         errors,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        internalServerError,
		HttpStatusCode: http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        notFound,
		HttpStatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        unathorized,
		HttpStatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message:        message,
		HttpErr:        forbidden,
		HttpStatusCode: http.StatusForbidden,
	}
}

package resterrors

import (
	"net/http"
)

const (
	badRequest          = "Bad Request"
	internalServerError = "Internal Server Error"
	notFound            = "Not Found"
	forbidden           = "Forbidden"
)

type RestErr struct {
	Message string  `json:"message"`
	Err     string  `json:"error"`
	Code    int     `json:"code"`
	Causes  []Cause `json:"causes,omitempty"`
}

type Cause struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewRestErr(message, err string, code int, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     err,
		Code:    code,
		Causes:  causes,
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     badRequest,
		Code:    http.StatusBadRequest,
	}
}

func NewBadRequestValidationError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     badRequest,
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     internalServerError,
		Code:    http.StatusBadRequest,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     notFound,
		Code:    http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     forbidden,
		Code:    http.StatusForbidden,
	}
}

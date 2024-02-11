package validation

import "fmt"

type ValidationError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("Key: %s, Error: %s", ve.Key, ve.Message)
}

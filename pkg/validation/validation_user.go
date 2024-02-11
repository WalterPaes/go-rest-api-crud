package validation

import (
	"encoding/json"
	"errors"

	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

func init() {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		unt := ut.New(en, en)
		transl, _ = unt.GetTranslator("en")
		en_translation.RegisterDefaultTranslations(val, transl)
	}
}

func ValidationUserError(validationErr error) *resterrors.RestErr {

	var jsonErr *json.UnmarshalTypeError
	var jsonValidationErr validator.ValidationErrors

	if errors.As(validationErr, &jsonErr) {
		return resterrors.NewBadRequestError("Invalid field type")
	}

	if errors.As(validationErr, &jsonValidationErr) {
		errorCauses := []resterrors.Cause{}

		for _, e := range validationErr.(validator.ValidationErrors) {
			cause := resterrors.Cause{
				Message: e.Translate(transl),
				Field:   e.Field(),
			}

			errorCauses = append(errorCauses, cause)
		}

		return resterrors.NewBadRequestValidationError("Some fields are invalid", errorCauses)
	}

	return resterrors.NewBadRequestError("Error trying to convert fields")
}

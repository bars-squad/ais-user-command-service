package validation

import (
	"fmt"

	"github.com/go-playground/validator"
)

type InvalidSchema struct {
	Field   string `json:"field"`
	Message any    `json:"message"`
}

func RequestBody(handler *validator.Validate, body interface{}) interface{} {
	err := handler.Struct(body)
	if err == nil {
		return nil
	}

	errorFields := err.(validator.ValidationErrors)
	var invalid []*InvalidSchema

	for _, errorField := range errorFields {
		invalid = append(invalid, &InvalidSchema{
			Field:   errorField.Field(),
			Message: fmt.Sprintf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value()),
		})
	}
	return invalid
}

package api

import (
	"github.com/ShubhKanodia/GoBank/util"
	"github.com/go-playground/validator/v10"
)

// custom validator for currency

func validateCurrency(field validator.FieldLevel) bool {
	// Check if the field's value is a string; if so, validate it as a supported currency
	if currency, ok := field.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

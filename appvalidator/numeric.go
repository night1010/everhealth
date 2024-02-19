package appvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var minNumeric validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}
	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.GreaterThanOrEqual(baseValue)
}

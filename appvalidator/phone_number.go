package appvalidator

import (
	"github.com/go-playground/validator/v10"
)

var phoneNumberPrefix validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if ok {
		if phone[0:1] != "0" {
			return false
		}
	}
	return true
}

var phoneNumberLength validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if ok {
		if len(phone) < 10 || len(phone) > 14 {
			return false
		}
	}
	return true
}

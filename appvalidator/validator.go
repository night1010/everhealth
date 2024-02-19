package appvalidator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidator() {
	v := binding.Validator.Engine().(*validator.Validate)
	v.RegisterValidation("phonenumberprefix", phoneNumberPrefix)
	v.RegisterValidation("phonenumberlength", phoneNumberLength)
	v.RegisterValidation("mind", minNumeric)

}

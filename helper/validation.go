package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorHandle(err error) string {
	var message string

	castedObject, ok := err.(validator.ValidationErrors)
	if ok {
		for _, v := range castedObject {
			switch v.Tag() {
			case "required":
				message = fmt.Sprintf("%s input value is required", v.Field())
			case "min":
				message = fmt.Sprintf("%s input value must be greater than %s character", v.Field(), v.Param())
			case "max":
				message = fmt.Sprintf("%s input value must be lower than %s character", v.Field(), v.Param())
			case "lte":
				message = fmt.Sprintf("%s input value must be below %s", v.Field(), v.Param())
			case "gte":
				message = fmt.Sprintf("%s input value must be above %s", v.Field(), v.Param())
			case "numeric":
				message = fmt.Sprintf("%s input value must be numeic", v.Field())
			case "url":
				message = fmt.Sprintf("%s input value must be am url", v.Field())
			case "email":
				message = fmt.Sprintf("%s input value must be an email", v.Field())
			case "password":
				message = fmt.Sprintf("%s input value must be filled", v.Field())
			}
		}
	}

	return message
}

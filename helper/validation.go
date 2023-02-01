package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorHandle(err error) string {
	messages := []string{}

	castedObject, ok := err.(validator.ValidationErrors)
	if ok {
		for _, v := range castedObject {
			switch v.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required", v.Field()))
			case "min":
				messages = append(messages, fmt.Sprintf("%s value must be greater than %s character", v.Field(), v.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s value must be lower than %s character", v.Field(), v.Param()))
			case "lte":
				messages = append(messages, fmt.Sprintf("%s value must be below %s", v.Field(), v.Param()))
			case "gte":
				messages = append(messages, fmt.Sprintf("%s value must be above %s", v.Field(), v.Param()))
			case "numeric":
				messages = append(messages, fmt.Sprintf("%s value must be numeic", v.Field()))
			case "url":
				messages = append(messages, fmt.Sprintf("%s value must be am url", v.Field()))
			case "email":
				messages = append(messages, fmt.Sprintf("%s value must be an email", v.Field()))
			case "password":
				messages = append(messages, fmt.Sprintf("%s value must be filled", v.Field()))
			}
		}
	}

	msg := strings.Join(messages, ", ")

	return msg
}

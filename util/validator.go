package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"error"`
}

func FormatValidatedField() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func ExtractValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Required"
	case "gte":
		return fmt.Sprintf("Should be greater than %s", fe.Param())
	case "gt":
		return fmt.Sprintf("Should be greater than %s", fe.Param())
	case "max":
		return fmt.Sprintf("Should be less than %s characters", fe.Param())
	case "email":
		return "Incorrect email format"
	default:
		return "Mismatch data type or malformed request"
	}
}

func ExtractUnmarshalError(je *json.UnmarshalTypeError) string {
	switch je.Field {
	case "author_id":
		return "Should be a number"
	case "title":
		return "Should be a string"
	case "description":
		return "Should be a string"
	case "quantity":
		return "Should be a number"
	case "cover":
		return "Should be a string"
	default:
		return "Mismatch data type or malformed request"
	}
}

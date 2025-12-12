package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = err.Error()
	return e
}

func NewValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		e.Errors[v.Field()] = fmt.Sprintf("%v", v.Tag())
	}
	return e
}

func AccessForbidden() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "access forbidden"
	return e
}

func NotFound(resource string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = fmt.Sprintf("%v not found", resource)
	return e
}

func InvalidField(field string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = fmt.Sprintf("Invalid field: %v", field)
	return e
}

func ErrorString(msg string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = msg
	return e
}

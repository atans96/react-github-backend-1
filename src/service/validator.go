package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"sync"
)

// structValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request. Gin provides a default implementation for this using
// https://github.com/go-playground/validator/tree/v8.18.2.
type structValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is a slice|array, the validation should be performed travel on every element.
	// If the received type is not a struct or slice|array, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) []*ErrorResponse

	// Engine returns the underlying validator engine which powers the
	// structValidator implementation.
	Engine() interface{}
}

// Validator is the default validator which implements the structValidator
// interface. It uses https://github.com/go-playground/validator/tree/v8.18.2
// under the hood.
var Validator structValidator = &defaultValidator{}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *defaultValidator) ValidateStruct(obj interface{}) []*ErrorResponse {
	if obj == nil {
		return nil
	}
	var errors []*ErrorResponse
	err := v.validateStruct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
func (v *defaultValidator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}
func (v *defaultValidator) validateStruct(obj interface{}) error {
	v.lazyInit()
	return v.validate.Struct(obj)
}
func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		fmt.Println("Validator")
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

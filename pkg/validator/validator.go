// Package validator provides input validation helpers.
//
// Wraps go-playground/validator behind a clean interface.
// Used by handlers to validate request bodies before
// passing data to use-cases.
//
// Example:
//
//	type Validator struct {
//	    validate *validator.Validate
//	}
//
//	func (v *Validator) ValidateStruct(s interface{}) error {
//	    return v.validate.Struct(s)
//	}
package validator

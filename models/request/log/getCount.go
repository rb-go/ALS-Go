package log

import (
	"gopkg.in/validator.v2"
)

type GetCount struct {
	Category string `validate:"nonzero"`
	SearchFilter map[string]interface{}
}

func (c GetCount) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}
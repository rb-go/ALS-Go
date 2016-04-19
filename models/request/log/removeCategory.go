package log

import "gopkg.in/validator.v2"

type RemoveCategory struct {
	Category string `validate:"nonzero"`
}



func (c RemoveCategory) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

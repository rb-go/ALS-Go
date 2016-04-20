package log

import "gopkg.in/validator.v2"

type RemoveLog struct {
	Category string `validate:"nonzero"`
	SearchFilter map[string]interface{}
}



func (c RemoveLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

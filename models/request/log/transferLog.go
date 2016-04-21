package log

import "gopkg.in/validator.v2"

type TransferLog struct {
	OldCategory string `validate:"nonzero"`
	NewCategory string `validate:"nonzero"`
	SearchFilter map[string]interface{}
}



func (c TransferLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

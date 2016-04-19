package log

import "gopkg.in/validator.v2"

type TransferLog struct {
	Old_category string `validate:"nonzero"`
	New_category string `validate:"nonzero"`
	Search_filter map[string]interface{}
}



func (c TransferLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

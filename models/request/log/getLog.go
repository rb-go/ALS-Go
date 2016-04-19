package log

import "gopkg.in/validator.v2"

type GetLog struct {
	Category string `validate:"nonzero"`
	Search_filter map[string]interface{}
	Limit int
	Offset int
	Sort_field string
	Sort_type string
}



func (c GetLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

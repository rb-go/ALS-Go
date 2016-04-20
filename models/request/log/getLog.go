package log

import "gopkg.in/validator.v2"

type GetLog struct {
	Category string `validate:"nonzero"`
	SearchFilter map[string]interface{}
	Limit int `validate:"max=1000, min=1"`
	Offset int
	Sort []string
}



func (c GetLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

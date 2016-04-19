package log

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

type AddCustom struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Category string `validate:"nonzero"`
	Level string `validate:"nonzero"`
	Message string `validate:"nonzero"`
	Timestamp int64 `validate:"nonzero"`
	ExpiresAt int64 `validate:"nonzero"`
	Tags []string
	AdditionalData interface{}
}

func (c AddCustom) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}
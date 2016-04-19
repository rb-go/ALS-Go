package log

import (
	"gopkg.in/validator.v2"
	"gopkg.in/mgo.v2/bson"
)

type Add struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Category string `validate:"nonzero"`
	Level string `validate:"nonzero"`
	Message string `validate:"nonzero"`
	Timestamp int64 `validate:"nonzero"`
	ExpiresAt int64 `validate:"nonzero"`
}

func (c Add) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}
package httpmodels

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
	"regexp"
	"fmt"
	"reflect"
	"errors"
)

type RequestLogAdd struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Category string `validate:"nonzero,CategoryNameValidators"`
	Level string `validate:"nonzero"`
	Message string `validate:"nonzero"`
	Timestamp int64 `validate:"nonzero"`
	ExpiresAt int64 `validate:"nonzero"`
}

func (c *RequestLogAdd) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}


type RequestLogAddCustom struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Category string `validate:"nonzero,CategoryNameValidators"`
	Level string `validate:"nonzero"`
	Message string `validate:"nonzero"`
	Timestamp int64 `validate:"nonzero"`
	ExpiresAt int64 `validate:"nonzero"`
	Tags []string
	AdditionalData interface{}
}

func (c RequestLogAddCustom) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}


type RequestLogGetCount struct {
	Category string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
}

func (c RequestLogGetCount) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}



type RequestLogGetLog struct {
	Category string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
	Limit int `validate:"max=1000, min=1"`
	Offset int
	Sort []string
}



func (c RequestLogGetLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}


type RequestLogRemoveCategory struct {
	Category string `validate:"nonzero,CategoryNameValidators"`
}



func (c RequestLogRemoveCategory) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}


type RequestLogRemoveLog struct {
	Category string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
}



func (c RequestLogRemoveLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

type RequestLogTransferLog struct {
	OldCategory string `validate:"nonzero,CategoryNameValidators"`
	NewCategory string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
}



func (c RequestLogTransferLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}


type RequestLogModifyTTL struct {
	Category string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{} `validate:"nonzero"`
	NewTTL int64 `validate:"nonzero"`
}

func (c RequestLogModifyTTL) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

func CategoryNameValidator(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}
	catNameValidator, err := regexp.Compile(`^[a-zA-Z0-9_]+$`)
	if err != nil {
		fmt.Println(err)
	}
	isMatch := catNameValidator.MatchString(st.String())
	if isMatch == false {
		return errors.New("field name not match [a-zA-Z0-9_]")
	}
	return nil
}
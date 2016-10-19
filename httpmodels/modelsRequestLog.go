package httpmodels

import (
	"errors"
	"reflect"
	"regexp"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

//RequestLogAdd Request Struct for LogAdd
type RequestLogAdd struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Category  string        `validate:"nonzero,CategoryNameValidators"`
	Level     string        `validate:"nonzero"`
	Message   string        `validate:"nonzero"`
	Timestamp int64         `validate:"nonzero"`
	ExpiresAt int64         `validate:"nonzero"`
}

//Validate Struct for LogAdd
func (c *RequestLogAdd) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogAddCustom Request Struct for LogAddCustom
type RequestLogAddCustom struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	Category       string        `validate:"nonzero,CategoryNameValidators"`
	Level          string        `validate:"nonzero"`
	Message        string        `validate:"nonzero"`
	Timestamp      int64         `validate:"nonzero"`
	ExpiresAt      int64         `validate:"nonzero"`
	Tags           []string
	AdditionalData interface{}
}

//Validate Struct for LogAddCustom
func (c RequestLogAddCustom) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogGetCount Request Struct for LogGetCount
type RequestLogGetCount struct {
	Category     string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
}

//Validate Struct for LogGetCount
func (c RequestLogGetCount) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogGetLog Request Struct for LogGetLog
type RequestLogGetLog struct {
	Category     string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
	Limit        int `validate:"max=1000, min=1"`
	Offset       int
	Sort         []string
}

//Validate Struct for LogGetLog
func (c RequestLogGetLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogRemoveCategory Request Struct for LogGetRemoveCategory
type RequestLogRemoveCategory struct {
	Category string `validate:"nonzero,CategoryNameValidators"`
}

//Validate Struct for LogGetRemoveCategory
func (c RequestLogRemoveCategory) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogRemoveLog Request Struct for LogGetRemoveLog
type RequestLogRemoveLog struct {
	Category     string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
}

//Validate Struct for LogGetRemoveLog
func (c RequestLogRemoveLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogTransferLog Request Struct for LogTransferLog
type RequestLogTransferLog struct {
	OldCategory  string `validate:"nonzero,CategoryNameValidators"`
	NewCategory  string `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{}
}

//Validate Struct for LogTransferLog
func (c RequestLogTransferLog) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//RequestLogModifyTTL Request Struct for LogModifyTTL
type RequestLogModifyTTL struct {
	Category     string                 `validate:"nonzero,CategoryNameValidators"`
	SearchFilter map[string]interface{} `validate:"nonzero"`
	NewTTL       int64                  `validate:"nonzero"`
}

//Validate Struct for LogModifyTTL
func (c RequestLogModifyTTL) Validate() error {
	if errs := validator.Validate(c); errs != nil {
		return errs
	}
	return nil
}

//CategoryNameValidator Method to validate category name. Not all names allowed. Allowed only: [a-zA-Z0-9_]
func CategoryNameValidator(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}
	catNameValidator, err := regexp.Compile(`^[a-zA-Z0-9_]+$`)
	if err != nil {
		return err
	}
	isMatch := catNameValidator.MatchString(st.String())
	if isMatch == false {
		return errors.New("field name not match [a-zA-Z0-9_]")
	}
	return nil
}

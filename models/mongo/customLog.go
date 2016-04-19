package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type CustomLog struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Category string
	Level string
	Message string
	Timestamp int64
	ExpiresAt time.Time `bson:"expiresAt"`
	Tags []string `json:",omitempty"`
	AdditionalData interface{} `bson:"additionalData" json:",omitempty"`
}

package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type CustomLog struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Category string `json:"category"`
	Level string `json:"level"`
	Message string `json:"message"`
	Timestamp int64 `json:"timestamp"`
	ExpiresAt time.Time `bson:"expiresAt"  json:"expiresAt"`
	Tags []string `json:"tags,omitempty"`
	AdditionalData interface{} `bson:"additionalData" json:"additionalData,omitempty"`
}

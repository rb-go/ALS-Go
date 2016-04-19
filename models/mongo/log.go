package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Log struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Category string
	Level string
	Message string
	Timestamp int64
	ExpiresAt time.Time `bson:"expiresAt"`
}

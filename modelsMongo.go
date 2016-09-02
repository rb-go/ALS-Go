package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoCustomLog struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Level string `json:"level"`
	Message string `json:"message"`
	Timestamp int64 `json:"timestamp"`
	ExpiresAt time.Time `bson:"expiresAt"  json:"-"`
	ExpiresAtShow int64 `bson:"expiresAtIntJustToShow"  json:"expiresAt"`
	Tags []string `json:"tags,omitempty"`
	AdditionalData interface{} `bson:"additionalData" json:"additionalData,omitempty"`
}


type MongoLog struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	Level string `json:"level"`
	Message string `json:"message"`
	Timestamp int64 `json:"timestamp"`
	ExpiresAt time.Time `bson:"expiresAt"  json:"-"`
	ExpiresAtShow int64 `bson:"expiresAtIntJustToShow"  json:"expiresAt"`
}
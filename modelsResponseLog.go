package main

import (
	"gopkg.in/mgo.v2/bson"
)

type ResponseLogAdd struct {
	LogId string `json:"logId"`
}


type ResponseLogGet struct {
	LogList []MongoCustomLog `json:"logList"`
}


type ResponseLogGetCategories struct {
	CategoriesList []string `json:"categoriesList"`
}


type ResponseLogRemoveCategory struct {
	Success int `json:"success"`
}


type ResponseLogRemoveLog struct {
	Matched int `json:"matched"`
	Removed int `json:"removed"`
}

type ResponseLogModifyTTL struct {
	Matched int `json:"matched"`
	Updated int `json:"updated"`
	UpsertedId interface{} `json:"upsertedId"`
}


type ResponseLogTransferLog struct {
	TransferedLogId []bson.ObjectId `json:"transferedLogId"`
}

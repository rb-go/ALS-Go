package httpmodels

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/Riftbit/ALS-Go/mongomodels"
)


//ResponseLogAdd Request Struct for LogAdd
type ResponseLogAdd struct {
	LogID string `json:"logId"`
}


//ResponseLogGet Request Struct for LogGet
type ResponseLogGet struct {
	LogList []mongomodels.MongoCustomLog `json:"logList"`
}


//ResponseLogGetCount Request Struct for LogGetCount
type ResponseLogGetCount struct {
	LogCount int `json:"logCount"`
}


//ResponseLogGetCategories Request Struct for LogGetCategories
type ResponseLogGetCategories struct {
	CategoriesList []string `json:"categoriesList"`
}


//ResponseLogRemoveCategory Request Struct for LogRemoveCategory
type ResponseLogRemoveCategory struct {
	Success int `json:"success"`
}


//ResponseLogRemoveLog Request Struct for LogRemoveLog
type ResponseLogRemoveLog struct {
	Matched int `json:"matched"`
	Removed int `json:"removed"`
}

//ResponseLogModifyTTL Request Struct for LogModifyTTL
type ResponseLogModifyTTL struct {
	Matched int `json:"matched"`
	Updated int `json:"updated"`
	UpsertedID interface{} `json:"upsertedId"`
}


//ResponseLogTransferLog Request Struct for LogTransferLog
type ResponseLogTransferLog struct {
	TransferedLogID []bson.ObjectId `json:"transferedLogId"`
}

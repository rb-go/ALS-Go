package main

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/Riftbit/ALS-Go/httpmodels"
	"github.com/stretchr/testify/assert"
)

var logApi *Log
var reqWithCorrectAuth *http.Request
var reqWithNotCorrectAuth *http.Request

func getReadyRequestFortests() {
	reqWithCorrectAuth, _ := http.NewRequest("POST", "http://api.local/", nil)
	reqWithCorrectAuth.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))

	reqWithNotCorrectAuth, _ := http.NewRequest("POST", "http://api.local/", nil)
	reqWithNotCorrectAuth.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
}

func init() {
	applicationExitFunction = func(c int) { okForTest = false }
	getReadyRequestFortests()
}

func TestApiLogAdd(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogAdd
	var reply *httpmodels.ResponseLogAdd

	args.Level = "error"
	args.Category = "api"
	args.Message = "This is test message to TestApiLogAdd"
	args.Timestamp = 1420074061
	args.ExpiresAt = 1490569965
	result := logApi.Add(reqWithCorrectAuth, args, reply)
	ass.Nil(result)

	result = logApi.Add(reqWithNotCorrectAuth, args, reply)
	ass.NotNil(result)
}

func TestApiLogAddCustom(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogAddCustom
	var reply *httpmodels.ResponseLogAdd

	type additionalDataStruct struct {
		Customer string
		State    int
	}

	args.Level = "error"
	args.Category = "api"
	args.Message = "This is test message to TestApiLogAddCustom"
	args.Timestamp = 1420074061
	args.ExpiresAt = 1490569965
	args.Tags = []string{"tags", "test", "go"}
	args.AdditionalData = additionalDataStruct{Customer: "apitester", State: 1}
	result := logApi.Add(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogGet(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogGetLog
	var reply *httpmodels.ResponseLogGet

	args.Category = "api"
	args.SearchFilter = struct{}{}
	args.Sort = []string{"+timestamp"}
	args.Limit = 1
	args.Offset = 0

	result := logApi.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogGetCount(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogGetCount
	var reply *httpmodels.ResponseLogGetCount

	args.Category = "api"
	args.SearchFilter = struct{}{}

	result := logApi.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogGetCategories(t *testing.T) {
	ass := assert.New(t)

	var args *struct{}
	var reply *httpmodels.ResponseLogGetCategories

	result := logApi.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogRemove(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogRemoveLog
	var reply *httpmodels.ResponseLogRemoveLog

	args.Category = "api"
	args.SearchFilter = struct{}{}

	result := logApi.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogRemoveCategory(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogRemoveCategory
	var reply *httpmodels.ResponseLogRemoveCategory

	args.Category = "api"

	result := logApi.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogTransfer(t *testing.T) {

	ass := assert.New(t)

	var args *httpmodels.RequestLogAdd
	var reply *httpmodels.ResponseLogAdd

	args.Level = "error"
	args.Category = "api"
	args.Message = "This is test message to TestApiLogTransfer"
	args.Timestamp = 1420074061
	args.ExpiresAt = 1490569965
	result := logApi.Add(reqWithCorrectAuth, args, reply)
	ass.Nil(result)

	var argss *httpmodels.RequestLogTransferLog
	var replyy *httpmodels.ResponseLogTransferLog

	argss.NewCategory = "api"
	argss.OldCategory = "api_new"
	argss.SearchFilter = struct{}{}

	result = logApi.Get(reqWithCorrectAuth, argss, replyy)
	ass.Nil(result)
}

func TestApiLogModifyTTL(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogModifyTTL
	var reply *httpmodels.ResponseLogModifyTTL

	args.Category = "api_new"
	args.SearchFilter = struct{}{}
	args.NewTTL = 1590569965

	result := logApi.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

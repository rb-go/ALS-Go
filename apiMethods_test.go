package main

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/Riftbit/ALS-Go/httpmodels"
	"github.com/stretchr/testify/assert"
)

var logAPI *Log
var reqWithCorrectAuth *http.Request
var reqWithNotCorrectAuth *http.Request
var emptySearchFilter map[string]interface{}

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
	result := logAPI.Add(reqWithCorrectAuth, args, reply)
	ass.Nil(result)

	result = logAPI.Add(reqWithNotCorrectAuth, args, reply)
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
	result := logAPI.AddCustom(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogGet(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogGetLog
	var reply *httpmodels.ResponseLogGet

	args.Category = "api"
	args.SearchFilter = emptySearchFilter
	args.Sort = []string{"+timestamp"}
	args.Limit = 1
	args.Offset = 0

	result := logAPI.Get(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogGetCount(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogGetCount
	var reply *httpmodels.ResponseLogGetCount

	args.Category = "api"
	args.SearchFilter = emptySearchFilter

	result := logAPI.GetCount(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogGetCategories(t *testing.T) {
	ass := assert.New(t)

	var args *struct{}
	var reply *httpmodels.ResponseLogGetCategories

	result := logAPI.GetCategories(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogRemove(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogRemoveLog
	var reply *httpmodels.ResponseLogRemoveLog

	args.Category = "api"
	args.SearchFilter = emptySearchFilter

	result := logAPI.Remove(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

func TestApiLogRemoveCategory(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogRemoveCategory
	var reply *httpmodels.ResponseLogRemoveCategory

	args.Category = "api"

	result := logAPI.RemoveCategory(reqWithCorrectAuth, args, reply)
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
	result := logAPI.Add(reqWithCorrectAuth, args, reply)
	ass.Nil(result)

	var argss *httpmodels.RequestLogTransferLog
	var replyy *httpmodels.ResponseLogTransferLog

	argss.NewCategory = "api"
	argss.OldCategory = "api_new"
	argss.SearchFilter = emptySearchFilter

	result = logAPI.Transfer(reqWithCorrectAuth, argss, replyy)
	ass.Nil(result)
}

func TestApiLogModifyTTL(t *testing.T) {
	ass := assert.New(t)

	var args *httpmodels.RequestLogModifyTTL
	var reply *httpmodels.ResponseLogModifyTTL

	args.Category = "api_new"
	args.SearchFilter = emptySearchFilter
	args.NewTTL = 1590569965

	result := logAPI.ModifyTTL(reqWithCorrectAuth, args, reply)
	ass.Nil(result)
}

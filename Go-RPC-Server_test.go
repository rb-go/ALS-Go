package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"fmt"

	"github.com/Riftbit/ALS-Go/httpmodels"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"gopkg.in/validator.v2"
)

var rawRequestBody string
var rawDataBody []byte
var okForTest bool
var testAdminMethodsList []string
var testBasicMethodsList []string

var logAPI *Log
var systemAPI *System
var emptySearchFilter map[string]interface{}

func init() {
	rawRequestBody = "{\"id\": \"55196eba27a55\", \"jsonrpc\": \"2.0\", \"method\": \"Log.GetCategories\", \"params\": {}}"
	applicationExitFunction = func(c int) { okForTest = false }

	validator.SetValidationFunc("CategoryNameValidators", httpmodels.CategoryNameValidator)
}

func TestLogPrintln(t *testing.T) {
	logPrintln("Starting tests!")
	os.Setenv("TESTING", "YES")
}

/*
====================================================
	CONFIG TESTS
====================================================
*/

//TestFailedInitConfigs - negative test
func TestFailedInitConfigsWhenFileNotExist(t *testing.T) {
	configPath = "./config.not.exists"
	initConfigs()
	if okForTest == true {
		t.Error("Wrong processing initConfigs when file not exists")
	}
	okForTest = true
}

//TestFailedInitConfigs - negative test
func TestFailedInitConfigs(t *testing.T) {
	configPath = "./config.wrong.yml"
	initConfigs()
	if okForTest == true {
		t.Error("Wrong processing initConfigs when config file not correct")
	}
	okForTest = true
}

func TestCommandLineFlags(t *testing.T) {
	parseCommandLineParams()
	configPath = "./config.smpl.yml"
}

func TestInitConfigs(t *testing.T) {
	initConfigs()
}

//TestFailInitLoggerWithWrongTimestampFormat - negative test
func TestFailInitLoggerWithWrongTimestampFormat(t *testing.T) {
	Configs.Log.TimestampFormat = "wrong"
	initLogger()
	if okForTest == true {
		t.Error("Wrong processing initConfigs when wrong Log TimestampFormat")
	}
	okForTest = true
	Configs.Log.TimestampFormat = "2006-01-02T15:04:05.999999999Z07:00"
}

//TestFailInitLoggerWithWrongFormatter - negative test
func TestFailInitLoggerWithWrongFormatter(t *testing.T) {
	Configs.Log.Formatter = "wrong"
	initLogger()
	if okForTest == true {
		t.Error("Wrong processing initConfigs when wrong Log Formatter")
	}
	okForTest = true
	Configs.Log.Formatter = "text"
}

//TestFailInitLoggerWithWrongFormatter - negative test
func TestFailInitLoggerWithWrongLogLevel(t *testing.T) {
	Configs.Log.LogLevel = "wrong"
	initLogger()
	if okForTest == true {
		t.Error("Wrong processing initConfigs when wrong LogLevel")
	}
	okForTest = true
	Configs.Log.LogLevel = "panic"
}

func TestInitLoggerWithJsonFormatter(t *testing.T) {
	Configs.Log.Formatter = "json"
	initLogger()
}

func TestInitLogger(t *testing.T) {
	initLogger()
}

/*
====================================================
	Go-RPC-Server TESTS
====================================================
*/

func TestInitRuntime(t *testing.T) {
	initRuntime()
}

func TestRpcPrepare(t *testing.T) {
	rpcPrepare()
}

func TestGetDataBody(t *testing.T) {
	req, err := http.NewRequest("POST", "http://api.local/", bytes.NewBufferString(rawRequestBody))
	if err != nil {
		t.Error("getDataBody Not correct http.NewRequest")
	}
	rawDataBody = getDataBody(req)
	if len(rawDataBody) < 1 {
		t.Error("getDataBody Not returned correct data", len(rawDataBody))
	}
}

func TestGetRequestJSON(t *testing.T) {
	ass := assert.New(t)
	jsonData, err := getRequestJSON(rawDataBody)
	if err != nil {
		t.Error(err)
	}
	ass.Equal("55196eba27a55", jsonData["id"], "Request ID should be equal")
}

func TestRegisterApi(t *testing.T) {
	testAdminMethodsList, testBasicMethodsList = registerAPI(rpcV2)
	ass := assert.New(t)
	ass.NotEmpty(testAdminMethodsList)
	ass.NotEmpty(testBasicMethodsList)
}

func TestFailInitDataBase(t *testing.T) {
	Configs.Db.DbConnectionString = Configs.Db.DbConnectionString + "&timeout=10ms"
	initDataBase()
	if okForTest == true {
		t.Error("Wrong processing initDataBase when wrong connection string")
	}
	okForTest = true
}

func TestInitDataBase(t *testing.T) {
	Configs.Db.DbType = "sqlite3"
	Configs.Db.DbConnectionString = "test_" + strconv.Itoa(int(time.Now().UTC().Unix())) + ".db"
	initDataBase()
}

func TestInitDatabaseStructure(t *testing.T) {
	initDatabaseStructure()
}

func TestInitDatabaseData(t *testing.T) {
	initDatabaseData(testAdminMethodsList, testBasicMethodsList)
}

func TestAuthentificator(t *testing.T) {
	ass := assert.New(t)

	ts := httptest.NewServer(authentificator(rpcV2))
	defer ts.Close()
	var u bytes.Buffer
	u.WriteString(string(ts.URL))

	EmptyRequestBody := ""
	correctRequestBody := "{\"id\": \"55196eba27a55\", \"jsonrpc\": \"2.0\", \"method\": \"Log.GetCategories\", \"params\": {}}"
	NotCorrectRequestBody := "{[\"id\": \"55196eba27.0\",] \"method\": \"Log.GetCategories\", \"params\": {}}]"

	//without auth
	client := &http.Client{}
	req, _ := http.NewRequest("POST", u.String(), bytes.NewBufferString(EmptyRequestBody))
	_, err := client.Do(req)
	ass.NoError(err)

	client = &http.Client{}
	req, _ = http.NewRequest("POST", u.String(), bytes.NewBufferString(EmptyRequestBody))
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
	_, err = client.Do(req)
	ass.NoError(err)

	client = &http.Client{}
	req, _ = http.NewRequest("POST", u.String(), bytes.NewBufferString(correctRequestBody))
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
	_, err = client.Do(req)
	ass.NoError(err)

	client = &http.Client{}
	req, _ = http.NewRequest("POST", u.String(), bytes.NewBufferString(NotCorrectRequestBody))
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
	_, err = client.Do(req)
	ass.NoError(err)
}

/*
====================================================
	MODELSDB TESTS
====================================================
*/
func TestFailCheckUserAuth(t *testing.T) {
	result := checkUserAuth("asdasd", "321")
	ass := assert.New(t)
	ass.Equal(false, result, "checkUserAuth should be wrong in this test")
	//cached
	result = checkUserAuth("asdasd", "321")
	ass.Equal(false, result, "checkUserAuth should be wrong in this test")
}

func TestCheckUserAuth(t *testing.T) {
	result := checkUserAuth(Configs.Admin.RootUser, Configs.Admin.RootPassword)
	ass := assert.New(t)
	ass.Equal(true, result, "checkUserAuth should be correct in this test")
	//cached
	result = checkUserAuth(Configs.Admin.RootUser, Configs.Admin.RootPassword)
	ass.Equal(true, result, "checkUserAuth should be correct in this test")
}

func TestFailCheckUserAccessToMethod(t *testing.T) {
	ass := assert.New(t)
	result := checkUserAccessToMethod("System.FlushCache", "dsa")
	ass.Equal(false, result, "checkUserAccessToMethod should not be correct in this test when wrong user")
	result = checkUserAccessToMethod("System.FlushCacheZ", "ergoz")
	ass.Equal(false, result, "checkUserAccessToMethod should not be correct in this test when wrong method")
	//cached
	result = checkUserAccessToMethod("System.FlushCache", "dsa")
	ass.Equal(false, result, "checkUserAccessToMethod should not be correct in this test when wrong user")
	result = checkUserAccessToMethod("System.FlushCacheZ", "ergoz")
	ass.Equal(false, result, "checkUserAccessToMethod should not be correct in this test when wrong method")
}

func TestCheckUserAccessToMethod(t *testing.T) {
	result := checkUserAccessToMethod("System.FlushCache", "ergoz")
	ass := assert.New(t)
	ass.Equal(true, result, "checkUserAuth should be correct in this test")
	//cached
	result = checkUserAccessToMethod("System.FlushCache", "ergoz")
	ass.Equal(true, result, "checkUserAuth should be correct in this test")
}

/*
====================================================
	AUTH TESTS
====================================================
*/

func TestParseRequestAuthData(t *testing.T) {
	req, err := http.NewRequest("POST", "http://api.local/", bytes.NewBufferString(rawRequestBody))
	if err != nil {
		t.Error("parseRequestAuthData Not correct http.NewRequest")
	}
	ass := assert.New(t)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
	user, password := parseRequestAuthData(req)
	ass.Equal(Configs.Admin.RootUser, user, "parseRequestAuthData should be correct in this test")
	ass.Equal(Configs.Admin.RootPassword, password, "parseRequestAuthData should be correct in this test")

	req, err = http.NewRequest("POST", "http://api.local/", bytes.NewBufferString(rawRequestBody))
	if err != nil {
		t.Error("parseRequestAuthData Not correct http.NewRequest")
	}
	user, password = parseRequestAuthData(req)
	ass.Equal("", user, "parseRequestAuthData should be correct in this test")
	ass.Equal("", password, "parseRequestAuthData should be correct in this test")
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("POST", "http://api.local/", bytes.NewBufferString(rawRequestBody))
	if err != nil {
		t.Error("getUser Not correct http.NewRequest")
	}
	ass := assert.New(t)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
	user := getUser(req)
	ass.Equal(Configs.Admin.RootUser, user, "getUser should be correct in this test")
}

func TestCheckAPIMethodAccess(t *testing.T) {
	ass := assert.New(t)
	isAllow := checkAPIMethodAccess(Configs.Admin.RootUser, "Log.Get")
	ass.Equal(true, isAllow, "checkAPIMethodAccess should be correct in this test")
	//cached
	isAllow = checkAPIMethodAccess(Configs.Admin.RootUser, "Log.Get")
	ass.Equal(true, isAllow, "checkAPIMethodAccess should be correct in this test")

	isAllow = checkAPIMethodAccess(Configs.Admin.RootUser, "ASD.DSA")
	ass.Equal(false, isAllow, "checkAPIMethodAccess should be correct in this test")
	//cached
	isAllow = checkAPIMethodAccess(Configs.Admin.RootUser, "ASD.DSA")
	ass.Equal(false, isAllow, "checkAPIMethodAccess should be correct in this test")
}

/*
====================================================
	HELPERS TESTS
====================================================
*/

func TestPrintObject(t *testing.T) {
	var result string
	m := make(map[string]int)
	m["route"] = 66
	result = printObject(m)
	ass := assert.New(t)
	ass.NotEmpty(result)
	ass.Equal("{\"route\":66}", result, "printObject data be equal")
}

func TestGetFuncName(t *testing.T) {
	result := getFuncName(1)
	ass := assert.New(t)
	ass.NotEmpty(result)
	ass.Contains(result, "TestGetFuncName", "getFuncName data be equal")
}

func TestGetLineCall(t *testing.T) {
	result := getLineCall(1)
	ass := assert.New(t)
	ass.NotZero(result)
	ass.Equal(getLineCall(1)-3, result, "getLineCall data be equal")
}

func TestGetFileCall(t *testing.T) {
	result := getFileCall(1)
	ass := assert.New(t)
	ass.NotEmpty(result)
	ass.Contains(result, "Go-RPC-Server_test.go", "getFileCall data be equal")
}

/*
====================================================
	API METHODS TESTS
====================================================
*/

func getReadyRequestForTests(correct bool) *http.Request {
	req, _ := http.NewRequest("POST", "http://api.local/", nil)
	if correct == true {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":"+Configs.Admin.RootPassword)))
	} else {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(Configs.Admin.RootUser+":wronguserpassword")))
	}
	return req
}

func TestApiLogAdd(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogAdd{}
	reply := httpmodels.ResponseLogAdd{}

	args.Level = "error"
	args.Category = "api"
	args.Message = "This is test message to TestApiLogAdd"
	args.Timestamp = 1420074061
	args.ExpiresAt = 1490569965

	result := logAPI.Add(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogAddCustom(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogAddCustom{}
	reply := httpmodels.ResponseLogAdd{}

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

	result := logAPI.AddCustom(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogGet(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogGetLog{}
	reply := httpmodels.ResponseLogGet{}

	args.Category = "api"
	args.SearchFilter = emptySearchFilter
	args.Sort = []string{"+timestamp"}
	args.Limit = 1
	args.Offset = 0

	result := logAPI.Get(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogGetCount(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogGetCount{}
	reply := httpmodels.ResponseLogGetCount{}

	args.Category = "api"
	args.SearchFilter = emptySearchFilter

	result := logAPI.GetCount(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogGetCategories(t *testing.T) {
	ass := assert.New(t)

	args := struct{}{}
	reply := httpmodels.ResponseLogGetCategories{}

	result := logAPI.GetCategories(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogRemove(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogRemoveLog{}
	reply := httpmodels.ResponseLogRemoveLog{}

	args.Category = "api"
	args.SearchFilter = emptySearchFilter

	result := logAPI.Remove(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogRemoveCategory(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogRemoveCategory{}
	reply := httpmodels.ResponseLogRemoveCategory{}

	args.Category = "api"

	result := logAPI.RemoveCategory(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestApiLogTransfer(t *testing.T) {

	ass := assert.New(t)

	args := httpmodels.RequestLogAdd{}
	reply := httpmodels.ResponseLogAdd{}

	args.Level = "error"
	args.Category = "api"
	args.Message = "This is test message to TestApiLogTransfer"
	args.Timestamp = 1420074061
	args.ExpiresAt = 1490569965

	result := logAPI.Add(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)

	argss := httpmodels.RequestLogTransferLog{}
	replyy := httpmodels.ResponseLogTransferLog{}

	argss.NewCategory = "api"
	argss.OldCategory = "api_new"
	argss.SearchFilter = emptySearchFilter

	result = logAPI.Transfer(getReadyRequestForTests(true), &argss, &replyy)
	ass.Nil(result)
}

func TestApiLogModifyTTL(t *testing.T) {
	ass := assert.New(t)

	args := httpmodels.RequestLogModifyTTL{}
	reply := httpmodels.ResponseLogModifyTTL{}

	type searchFilterWithExp map[string]interface{}
	paramsSearch := searchFilterWithExp{"ExpiresAt": 1490569965}

	args.Category = "api_new"
	args.SearchFilter = paramsSearch
	args.NewTTL = 1590569965

	result := logAPI.ModifyTTL(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
}

func TestDeleteDataBaseAfterMethodTests(t *testing.T) {
	DBConn.Close()
	err := os.Remove(Configs.Db.DbConnectionString)
	ass := assert.New(t)
	ass.Nil(err)
}

/*
====================================================
	SYSTEM API METHODS TESTS
====================================================
*/

func TestApiSystemGetCacheAll(t *testing.T) {
	ass := assert.New(t)

	args := struct{}{}
	reply := struct {
		Count int
		Items map[string]cache.Item
	}{}

	result := systemAPI.GetCacheAll(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
	ass.NotEqual(0, reply.Count)
}

func TestApiSystemGetCache(t *testing.T) {
	ass := assert.New(t)

	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", Configs.Admin.RootUser, "api")

	args := struct{ Key string }{Key: cahceKey}
	reply := struct {
		Key  string
		Data interface{}
	}{}

	result := systemAPI.GetCache(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
	ass.NotNil(reply.Data)
	ass.Equal(1, reply.Data)
}

func TestApiSystemDeleteCache(t *testing.T) {
	ass := assert.New(t)

	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", Configs.Admin.RootUser, "api")

	args := struct{ Key string }{Key: cahceKey}
	reply := struct{ Status int }{}

	result := systemAPI.DeleteCache(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
	ass.Equal(1, reply.Status)
}

func TestApiSystemGetCacheAfterDelete(t *testing.T) {
	ass := assert.New(t)

	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", Configs.Admin.RootUser, "api")

	args := struct{ Key string }{Key: cahceKey}
	reply := struct {
		Key  string
		Data interface{}
	}{}

	result := systemAPI.GetCache(getReadyRequestForTests(true), &args, &reply)
	ass.NotNil(result)
	ass.Equal(json2.ErrNullResult, result)
}

func TestApiSystemSetCache(t *testing.T) {
	ass := assert.New(t)

	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", Configs.Admin.RootUser, "api")

	args := struct {
		Key  string
		Data interface{}
		TTL  time.Duration
	}{}
	reply := struct{ Status int }{}

	args.Key = cahceKey
	args.Data = 1
	args.TTL = cache.NoExpiration

	result := systemAPI.SetCache(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
	ass.Equal(1, reply.Status)
}

func TestApiSystemGetCacheAfterSet(t *testing.T) {
	ass := assert.New(t)

	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", Configs.Admin.RootUser, "api")

	args := struct{ Key string }{Key: cahceKey}
	reply := struct {
		Key  string
		Data interface{}
	}{}

	result := systemAPI.GetCache(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
	ass.NotNil(reply.Data)
	ass.Equal(1, reply.Data)
}

func TestApiSystemFlushCache(t *testing.T) {
	ass := assert.New(t)

	args := struct{}{}
	reply := struct{ Status int }{}

	result := systemAPI.FlushCache(getReadyRequestForTests(true), &args, &reply)
	ass.Nil(result)
	ass.Equal(1, reply.Status)
}

func TestApiSystemGetCacheAfterFlush(t *testing.T) {
	ass := assert.New(t)

	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", Configs.Admin.RootUser, "api")

	args := struct{ Key string }{Key: cahceKey}
	reply := struct {
		Key  string
		Data interface{}
	}{}

	result := systemAPI.GetCache(getReadyRequestForTests(true), &args, &reply)
	ass.NotNil(result)
	ass.Equal(json2.ErrNullResult, result)
}

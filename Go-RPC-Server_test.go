package main

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawRequestBody string
var rawDataBody []byte
var okForTest bool

func init() {
	rawRequestBody = "{\"id\": \"55196eba27a55\", \"jsonrpc\": \"2.0\", \"method\": \"Log.GetCategories\", \"params\": {}}"
	applicationExitFunction = func(c int) { okForTest = false }
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
	if testing.Short() {
		configPath = "./config.smpl.yml"
		t.Skip("skipping test; this test not for race or run in more than 1 thread")
	}
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
	if testing.Short() {
		t.Skip("skipping test; this test not for race or run in more than 1 thread")
	}
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
	adminMethodsList, basicMethodsList := registerAPI(rpcV2)
	ass := assert.New(t)
	ass.NotEmpty(adminMethodsList)
	ass.NotEmpty(basicMethodsList)
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

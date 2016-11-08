package main

import (
	"bytes"
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

var rawRequestBody string = "{\"id\": \"55196eba27a55\", \"jsonrpc\": \"2.0\", \"method\": \"Log.GetCategories\", \"params\": {}}"
var rawDataBody []byte
var okForTest bool = false

func init() {
	applicationExitFunction = func(c int) { okForTest = true }
}

//TestFailedInitConfigs - negative test
func TestFailedInitConfigsWhenFileNotExist(t *testing.T) {
	configPath = "./config.not.exists"
	initConfigs()
	if okForTest == false {
		t.Error("Wrong processing initConfigs when file not exists")
	}
	okForTest = false
}

//TestFailedInitConfigs - negative test
func TestFailedInitConfigs(t *testing.T) {
	configPath = "./config.wrong.yml"
	initConfigs()
	if okForTest == false {
		t.Error("Wrong processing initConfigs when config file not correct")
	}
	okForTest = false
}

func TestInitConfigs(t *testing.T) {
	configPath = "./config.smpl.yml"
	initConfigs()
}

func TestInitRuntime(t *testing.T) {
	initRuntime()
}

//TestFailInitLoggerWithWrongTimestampFormat - negative test
func TestFailInitLoggerWithWrongTimestampFormat(t *testing.T) {
	Configs.Log.TimestampFormat = "wrong"
	initLogger()
	if okForTest == false {
		t.Error("Wrong processing initConfigs when wrong Log TimestampFormat")
	}
	okForTest = false
	Configs.Log.TimestampFormat = "2006-01-02T15:04:05.999999999Z07:00"
}

//TestFailInitLoggerWithWrongFormatter - negative test
func TestFailInitLoggerWithWrongFormatter(t *testing.T) {
	Configs.Log.Formatter = "wrong"
	initLogger()
	if okForTest == false {
		t.Error("Wrong processing initConfigs when wrong Log Formatter")
	}
	okForTest = false
	Configs.Log.Formatter = "text"
}

//TestFailInitLoggerWithWrongFormatter - negative test
func TestFailInitLoggerWithWrongLogLevel(t *testing.T) {
	Configs.Log.LogLevel = "wrong"
	initLogger()
	if okForTest == false {
		t.Error("Wrong processing initConfigs when wrong LogLevel")
	}
	okForTest = false
	Configs.Log.Formatter = "panic"
}

func TestInitLoggerWithJsonFormatter(t *testing.T) {
	Configs.Log.Formatter = "json"
	initLogger()
}

func TestInitLogger(t *testing.T) {
	initLogger()
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

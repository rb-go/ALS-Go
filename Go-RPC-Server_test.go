package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawRequestBody string = "{\"id\": \"55196eba27a55\", \"jsonrpc\": \"2.0\", \"method\": \"Log.GetCategories\", \"params\": {}}"
var rawDataBody []byte
var ok bool = false

func init() {
	application_exit_function = func(c int) { ok = true }
}

//TestFailedInitConfigs - negative test
func TestFailedInitConfigs(t *testing.T) {
	configPath = "./config.not.exists"
	initConfigs()
	if ok == false {
		t.Error("Wrong processing initConfigs when file not exists")
	}
	ok = false
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
	if ok == false {
		t.Error("Wrong processing initConfigs when file not exists")
	}
	ok = false
	Configs.Log.TimestampFormat = "2006-01-02T15:04:05.999999999Z07:00"
}

//TestFailInitLoggerWithWrongFormatter - negative test
func TestFailInitLoggerWithWrongFormatter(t *testing.T) {
	Configs.Log.Formatter = "wrong"
	initLogger()
	if ok == false {
		t.Error("Wrong processing initConfigs when file not exists")
	}
	ok = false
	Configs.Log.Formatter = "text"
}

func TestInitLogger(t *testing.T) {
	initLogger()
}

func TestGetDataBody(t *testing.T) {
	req := httptest.NewRequest("POST", "http://api.local/", bytes.NewBufferString(rawRequestBody))
	rawDataBody = getDataBody(req)
	if len(rawDataBody) < 10 {
		t.Error("getDataBody Not returned correct data")
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

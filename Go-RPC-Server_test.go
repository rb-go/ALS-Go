package main

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawRequestBody string
var rawDataBody []byte

func init() {
	rawRequestBody = "{\"id\": \"55196eba27a55\", \"jsonrpc\": \"2.0\", \"method\": \"Log.GetCategories\", \"params\": {}}"
}

func TestInitConfigs(t *testing.T) {
	configPath = "./config.smpl.yml"
	initConfigs()
}

func TestInitRuntime(t *testing.T) {
	initRuntime()
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

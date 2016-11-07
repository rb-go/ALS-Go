package main

import (
	"encoding/json"
	"runtime"
)

func printObject(v interface{}) string {
	res2B, _ := json.Marshal(v)
	return string(res2B)
}

func getFuncName(level int) string {
	pc, _, _, _ := runtime.Caller(level)
	return runtime.FuncForPC(pc).Name()
}

func getLineCall(level int) int {
	_, _, line, _ := runtime.Caller(level)
	return line
}

func getFileCall(level int) string {
	_, file, _, _ := runtime.Caller(level)
	return file
}

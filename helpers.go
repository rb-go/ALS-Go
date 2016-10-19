package main

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func printObject(v interface{}) {
	res2B, _ := json.Marshal(v)
	fmt.Println(string(res2B))
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

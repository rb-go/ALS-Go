package main

import (
	"fmt"
	"reflect"

	"github.com/gorilla/rpc/v2"
)

func register(rpcV2 *rpc.Server) {
	Logger.Info("Registering exported methods")
	rpcV2.RegisterService(new(Log), "")
	rpcV2.RegisterService(new(System), "")

	var adminMethodsList []string
	var basicMethodsList []string

	// prints a concise summary of the exported API calls
	listMethods := func(m interface{}) {
		typ := reflect.TypeOf(m)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		fooType := reflect.TypeOf(m)
		for i := 0; i < fooType.NumMethod(); i++ {
			method := fooType.Method(i)
			args := reflect.New(method.Type.In(2).Elem()).Elem().Interface()
			resp := reflect.New(method.Type.In(3).Elem()).Elem().Interface()
			Logger.Debugf("request = api.call('%s.%s', %+v) # response: %+v", typ.Name(), method.Name, args, resp)
			if typ.Name() == "System" {
				adminMethodsList = append(adminMethodsList, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			} else {
				basicMethodsList = append(basicMethodsList, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			}
		}
	}
	Logger.Debug("Start exported methods names")
	listMethods(new(System))
	listMethods(new(Log))
	Logger.Debug("End exporten methods names")

	initDatabaseStructure()
	initDatabaseData(adminMethodsList, basicMethodsList)
}

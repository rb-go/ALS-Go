package main

import (
	"github.com/gorilla/rpc/v2"
	"reflect"
	"fmt"
)


func Register(rpc_v2 *rpc.Server) {
	Logger.Info("Registering exported methods")
	rpc_v2.RegisterService(new(Log), "")
	rpc_v2.RegisterService(new(System), "")

	var admin_methods_list []string
	var basic_methods_list []string

	// prints a concise summary of the exported API calls
	list_methods := func(m interface{}) {
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
				admin_methods_list = append(admin_methods_list, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			} else {
				basic_methods_list = append(basic_methods_list, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			}
		}
	}
	Logger.Debug("Start exported methods names")
	list_methods(new(System))
	list_methods(new(Log))
	Logger.Debug("End exporten methods names")

	InitDatabaseStructure()
	InitDatabaseData(admin_methods_list, basic_methods_list)
}


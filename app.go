package main

import (
	"github.com/gorilla/rpc/v2"
	"reflect"
	"fmt"
	"log"
)


func Register(rpc_v2 *rpc.Server) {
	log.Println("... REGISTERING METHODS ...")
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
			log.Printf("request = api.call('%s.%s', %+v) # response: %+v", typ.Name(), method.Name, args, resp)
			if typ.Name() == "System" {
				admin_methods_list = append(admin_methods_list, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			} else {
				basic_methods_list = append(basic_methods_list, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			}
		}
	}
	log.Println("START EXPORTED METHOD NAMES")
	list_methods(new(System))
	list_methods(new(Log))
	log.Println("END EXPORTED METHOD NAMES")

	InitDatabaseStructure()
	InitDatabaseData(admin_methods_list, basic_methods_list)
}


package main

import (
	"net/http"
	"fmt"
	"github.com/patrickmn/go-cache"
)

func CheckAuth(r *http.Request) bool {
	username, password, _ := r.BasicAuth()
	if username == "" && password == "" {
		return false
	}
	return CheckUserAuth(username, password)
}

func GetUser(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

func CheckAPIMethodAccess(r *http.Request, json_data map[string]interface{}) bool {
	username := GetUser(r)
	method_name := json_data["method"].(string)

	access_right, found := Cache.Get(fmt.Sprintf("Access:%s:%s", username, method_name))
	if found == false {
		if !CheckUserAccessToMethod(method_name, username) {
			Logger.Warnf("No permissions for user '%s' to method '%s'", username, method_name)
			Cache.Set(fmt.Sprintf("Access:%s:%s", username, method_name), false, cache.NoExpiration)
			return false
		} else {
			Cache.Set(fmt.Sprintf("Access:%s:%s", username, method_name), true, cache.NoExpiration)
			return true
		}
	} else {
		return access_right.(bool)
	}
}
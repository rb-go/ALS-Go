package main

import (
	"net/http"
	"fmt"
	"github.com/patrickmn/go-cache"
)

func checkAuth(r *http.Request) bool {
	username, password, _ := r.BasicAuth()
	if username == "" && password == "" {
		return false
	}
	return checkUserAuth(username, password)
}

func getUser(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

func checkAPIMethodAccess(r *http.Request, jsonData map[string]interface{}) bool {
	username := getUser(r)
	methodName := jsonData["method"].(string)

	accessRight, found := Cache.Get(fmt.Sprintf("Access:%s:%s", username, methodName))
	if found == false {
		if !checkUserAccessToMethod(methodName, username) {
			Logger.Warnf("No permissions for user '%s' to method '%s'", username, methodName)
			Cache.Set(fmt.Sprintf("Access:%s:%s", username, methodName), false, cache.NoExpiration)
			return false
		}
		Cache.Set(fmt.Sprintf("Access:%s:%s", username, methodName), true, cache.NoExpiration)
		return true
	}
	return accessRight.(bool)
}
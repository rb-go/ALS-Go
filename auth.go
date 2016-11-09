package main

import (
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
)

func parseRequestAuthData(r *http.Request) (string, string) {
	username, password, _ := r.BasicAuth()
	return username, password
}

func getUser(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

func checkAPIMethodAccess(userName string, methodName string) bool {
	accessRight, found := Cache.Get(fmt.Sprintf("Access:%s:%s", userName, methodName))
	if found == false {
		if !checkUserAccessToMethod(methodName, userName) {
			Logger.Warnf("No permissions for user '%s' to method '%s'", userName, methodName)
			Cache.Set(fmt.Sprintf("Access:%s:%s", userName, methodName), false, cache.NoExpiration)
			return false
		}
		Cache.Set(fmt.Sprintf("Access:%s:%s", userName, methodName), true, cache.NoExpiration)
		return true
	}
	return accessRight.(bool)
}

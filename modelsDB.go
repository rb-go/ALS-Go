package main

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

//Method struct for db
type Method struct {
	ID   int    `gorm:"primary_key"`
	Name string `sql:"size:255;not null;unique"` // Default size for string is 255, you could reset it with this tag
	//Users        []User  `gorm:"many2many:user2method;"` // Many-To-Many relationship, 'user_languages' is join table
}

//TableName Exporting table name
func (c Method) TableName() string {
	return "methods"
}

//User struct for db
type User struct {
	ID       int      `gorm:"primary_key"`
	Login    string   `sql:"type:varchar(40);not null;unique"`
	Password string   `sql:"type:varchar(40)"`
	Email    string   `sql:"size:255"`
	Status   int      `sql:"type:int(2);not null;DEFAULT:0"`
	Methods  []Method `gorm:"many2many:user2method;"`
}

//TableName Exporting table name
func (c User) TableName() string {
	return "users"
}

func checkUserAccessToMethod(method, user string) bool {
	accessRight, found := Cache.Get(fmt.Sprintf("UserMethodAccess:%s:%s", user, method))
	if found == false {
		var u User
		db := DBConn.Preload("Methods", Method{Name: method}).First(&u, User{Login: user})
		if db.Error != nil {
			Logger.Error(db.Error)
			Cache.Set(fmt.Sprintf("UserMethodAccess:%s:%s", user, method), false, cache.NoExpiration)
			return false
		}
		Logger.Debug("[checkUserAccessToMethod] ", printObject(db.Value))
		if u.Methods == nil {
			Cache.Set(fmt.Sprintf("UserMethodAccess:%s:%s", user, method), false, cache.NoExpiration)
			return false
		}
		Cache.Set(fmt.Sprintf("UserMethodAccess:%s:%s", user, method), true, cache.NoExpiration)
		return true
	}
	return accessRight.(bool)
}

func checkUserAuth(user, password string) bool {
	accessRight, found := Cache.Get(fmt.Sprintf("UserAuth:%s:%s", user, password))
	if found == false {
		var u User
		db := DBConn.First(&u, User{Login: user, Password: password})
		Logger.Debug("[checkUserAuth] ", db.Value)
		if db.Error != nil {
			Cache.Set(fmt.Sprintf("UserAuth:%s:%s", user, password), false, cache.NoExpiration)
			Logger.Errorf("DB_ERROR: %s | Params: %s,%s", db.Error, user, password)
			return false
		}
		Cache.Set(fmt.Sprintf("UserAuth:%s:%s", user, password), true, cache.NoExpiration)
		return true
	}
	return accessRight.(bool)
}

func isDBConnected() bool {
	err := DBConn.DB().Ping()
	if err != nil {
		return false
	}
	return true
}

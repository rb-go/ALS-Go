package main


import (
	"fmt"
	"github.com/patrickmn/go-cache"
)

type method struct {
	ID           int     `gorm:"primary_key"`
	Name         string  `sql:"size:255;not null;unique"` // Default size for string is 255, you could reset it with this tag
							      //Users        []User  `gorm:"many2many:user2method;"` // Many-To-Many relationship, 'user_languages' is join table
}

//TableName Exporting table name
func (c method) TableName() string {
	return "methods"
}

type user struct {
	ID           int      `gorm:"primary_key"`
	Login        string   `sql:"type:varchar(20);not null;unique"`
	Password     string   `sql:"type:varchar(40)"`
	Email        string   `sql:"size:255"`
	Status       int      `sql:"type:int(2);not null;DEFAULT:0"`
	Methods      []method `gorm:"many2many:user2method;"`
}

//TableName Exporting table name
func (c user) TableName() string {
	return "users"
}

func checkUserAccessToMethod(method, user string) bool {
	var u user
	db := DBConn.Preload("Methods", method{Name:method}).First(&u, user{Login: user})
	if db.Error != nil {
		Logger.Error(db.Error)
	}
	if u.Methods == nil {
		return false
	}
	return true
}

func checkUserAuth(user, password string) bool {
	accessRight, found := Cache.Get(fmt.Sprintf("UserAuth:%s:%s", user,password))
	if found == false {
		var u user
		db := DBConn.First(&u, user{Login: user, Password: password})
		if db.Error != nil {
			Cache.Set(fmt.Sprintf("UserAuth:%s:%s", user, password), false, cache.NoExpiration)
			Logger.Errorf("DB_ERROR: %s",db.Error)
			return false
		}
		Cache.Set(fmt.Sprintf("UserAuth:%s:%s", user, password), true, cache.NoExpiration)
		return true
	}
	return accessRight.(bool)
}

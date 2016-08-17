package db

import (
	"gitlab.com/ergoz/ALS-Go/configs"
	"log"
	"fmt"
	"github.com/patrickmn/go-cache"
)

type User struct {
	Id           int      `gorm:"primary_key"`
	Login        string   `sql:"type:varchar(20);not null;unique"`
	Password     string   `sql:"type:varchar(40)"`
	Email        string   `sql:"size:255"`
	Status       int      `sql:"type:int(2);not null;DEFAULT:0"`
	Methods      []Method `gorm:"many2many:user2method;"`
}

func (c User) TableName() string {
	return "users"
}

func CheckUserAccessToMethod(method, user string) bool {
	var u User
	db := configs.DBConn.Preload("Methods", Method{Name:method}).First(&u, User{Login: user})
	if db.Error != nil {
		log.Println(db.Error)
	}
	if u.Methods == nil {
		return false
	} else {
		return true
	}
}

func CheckUserAuth(user, password string) bool {

	access_right, found := configs.Cache.Get(fmt.Sprintf("UserAuth:%s:%s", user,password))
	if found == false {
		var u User
		db := configs.DBConn.First(&u, User{Login: user, Password: password})
		if db.Error != nil {
			configs.Cache.Set(fmt.Sprintf("UserAuth:%s:%s", user, password), false, cache.NoExpiration)
			log.Printf("DB_ERROR: %s",db.Error)
			return false
		} else {
			configs.Cache.Set(fmt.Sprintf("UserAuth:%s:%s", user, password), true, cache.NoExpiration)
			return true
		}
	}
	return access_right.(bool)

}

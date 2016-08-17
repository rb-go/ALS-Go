package models

import "gitlab.com/ergoz/ALS-Go/configs"
import "gitlab.com/ergoz/ALS-Go/models/db"

func InitDatabaseStructure() {
	configs.DBConn.AutoMigrate(&db.Method{},&db.User{})
}

func InitDatabaseData(admin_methods_list, basic_methods_list []string) {
	allMethods := []db.Method{}
	adminMethods := []db.Method{}
	otherMethods := []db.Method{}

	for _,method_name := range admin_methods_list {
		var method = db.Method{}
		configs.DBConn.FirstOrCreate(&method, db.Method{Name: method_name})
		adminMethods = append(adminMethods, method)
		allMethods = append(allMethods, method)
	}

	for _,method_name := range basic_methods_list {
		var method = db.Method{}
		configs.DBConn.FirstOrCreate(&method, db.Method{Name: method_name})
		otherMethods = append(otherMethods, method)
		allMethods = append(allMethods, method)
	}

	var user = db.User{}
	user_to_create := db.User{Password: configs.Configs.Admin.RootPassword, Email: configs.Configs.Admin.RootPassword, Status: 1, Methods: allMethods}
	configs.DBConn.Attrs(user_to_create).FirstOrCreate(&user, db.User{Login: configs.Configs.Admin.RootUser})

	//configs.DBConn.Model(&user_to_create).Association("Methods").Append(db.Method{Name: "System.Test"})

	//var methods []db.Method
	//configs.DBConn.Model(&user_to_create).Association("Methods").Find(&methods)

	//userd := db.User{Id:1}
	//configs.DBConn.Model(&userd).Association("Methods").Append(db.Method{Id:1})
	//log.Printf("User method count: %d", configs.DBConn.Model(&userd).Association("Methods").Count())
	//userd = db.User{Id:1}
	//configs.DBConn.Model(&user).Association("Methods").Delete(db.Method{Id:1})
	//log.Printf("User method count: %d", configs.DBConn.Model(&userd).Association("Methods").Count())
	//userd = db.User{Id:1}
	//configs.DBConn.Model(&userd).Association("Methods").Append(db.Method{Id:1})
	//configs.DBConn.Model(&userd).Association("Methods").Append(db.Method{Id:2})
	//configs.DBConn.Model(&userd).Association("Methods").Append(db.Method{Id:3})
	//log.Printf("User method count: %d", configs.DBConn.Model(&userd).Association("Methods").Count())
	//userd = db.User{Id:1}
	//configs.DBConn.Model(&user).Association("Methods").Clear()
	//log.Printf("User method count: %d", configs.DBConn.Model(&userd).Association("Methods").Count())

	//var user_new = db.User{}
	//configs.DBConn.Preload("Methods").First(&user_new, db.User{Login:"ergoz"})
	//user_new.DeleteAccessMethod(db.Method{Name:"System.Test"})
	//configs.DBConn.Save(&user_new)

	//log.Println(user_new)
	//log.Println(allMethods)
	//log.Println(adminMethods)
	//log.Println(otherMethods)

	//db.CheckUserAccessToMethod()
}
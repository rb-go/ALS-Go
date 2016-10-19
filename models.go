package main


func initDatabaseStructure() {
	DBConn.AutoMigrate(&Method{},&User{})
}

func initDatabaseData(adminMethodsList, basicMethodsList []string) {
	allMethods := []Method{}
	adminMethods := []Method{}
	otherMethods := []Method{}

	for _,methodName := range adminMethodsList {
		var method = Method{}
		DBConn.FirstOrCreate(&method, Method{Name: methodName})
		adminMethods = append(adminMethods, method)
		allMethods = append(allMethods, method)
	}

	for _,methodName := range basicMethodsList {
		var method = Method{}
		DBConn.FirstOrCreate(&method, Method{Name: methodName})
		otherMethods = append(otherMethods, method)
		allMethods = append(allMethods, method)
	}

	var user = User{}
	userToCreate := User{Password: Configs.Admin.RootPassword, Email: Configs.Admin.RootPassword, Status: 1, Methods: allMethods}
	DBConn.Attrs(userToCreate).FirstOrCreate(&user, User{Login: Configs.Admin.RootUser})

	//DBConn.Model(&user_to_create).Association("Methods").Append(Method{Name: "System.Test"})

	//var methods []Method
	//DBConn.Model(&user_to_create).Association("Methods").Find(&methods)

	//userd := User{Id:1}
	//DBConn.Model(&userd).Association("Methods").Append(Method{Id:1})
	//log.Printf("User method count: %d", DBConn.Model(&userd).Association("Methods").Count())
	//userd = User{Id:1}
	//DBConn.Model(&user).Association("Methods").Delete(Method{Id:1})
	//log.Printf("User method count: %d", DBConn.Model(&userd).Association("Methods").Count())
	//userd = User{Id:1}
	//DBConn.Model(&userd).Association("Methods").Append(Method{Id:1})
	//DBConn.Model(&userd).Association("Methods").Append(Method{Id:2})
	//DBConn.Model(&userd).Association("Methods").Append(Method{Id:3})
	//log.Printf("User method count: %d", DBConn.Model(&userd).Association("Methods").Count())
	//userd = User{Id:1}
	//DBConn.Model(&user).Association("Methods").Clear()
	//log.Printf("User method count: %d", DBConn.Model(&userd).Association("Methods").Count())

	//var user_new = User{}
	//DBConn.Preload("Methods").First(&user_new, User{Login:"ergoz"})
	//user_new.DeleteAccessMethod(Method{Name:"System.Test"})
	//DBConn.Save(&user_new)

	//log.Println(user_new)
	//log.Println(allMethods)
	//log.Println(adminMethods)
	//log.Println(otherMethods)

	//CheckUserAccessToMethod()
}
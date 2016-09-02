package main


func InitDatabaseStructure() {
	DBConn.AutoMigrate(&Method{},&User{})
}

func InitDatabaseData(admin_methods_list, basic_methods_list []string) {
	allMethods := []Method{}
	adminMethods := []Method{}
	otherMethods := []Method{}

	for _,method_name := range admin_methods_list {
		var method = Method{}
		DBConn.FirstOrCreate(&method, Method{Name: method_name})
		adminMethods = append(adminMethods, method)
		allMethods = append(allMethods, method)
	}

	for _,method_name := range basic_methods_list {
		var method = Method{}
		DBConn.FirstOrCreate(&method, Method{Name: method_name})
		otherMethods = append(otherMethods, method)
		allMethods = append(allMethods, method)
	}

	var user = User{}
	user_to_create := User{Password: Configs.Admin.RootPassword, Email: Configs.Admin.RootPassword, Status: 1, Methods: allMethods}
	DBConn.Attrs(user_to_create).FirstOrCreate(&user, User{Login: Configs.Admin.RootUser})

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
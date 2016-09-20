package main

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/patrickmn/go-cache"
	"github.com/riftbit/ALS-Go/httpmodels"
	"github.com/riftbit/ALS-Go/mongomodels"
)


type Log struct{}

func (h *Log) Add(r *http.Request, args *httpmodels.RequestLogAdd, reply *httpmodels.ResponseLogAdd) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("[" + GetFuncName(1) + "] CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	args.ID = bson.NewObjectId()

	logData := mongomodels.MongoLog{}
	logData.ID = args.ID
	logData.Level = args.Level
	logData.Message = args.Message
	logData.Timestamp = args.Timestamp
	logData.ExpiresAt = time.Unix(args.ExpiresAt, 0)
	logData.ExpiresAtShow = args.ExpiresAt

	if InsertMGO(collection, logData) != nil {
		log.Println("InsertMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert Problems"}
	}

	reply.LogId=args.ID.Hex()
	return nil
}



func (h *Log) AddCustom(r *http.Request, args *httpmodels.RequestLogAddCustom, reply *httpmodels.ResponseLogAdd) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	args.ID = bson.NewObjectId()

	logData := mongomodels.MongoCustomLog{}
	logData.ID = args.ID
	logData.Level = args.Level
	logData.Message = args.Message
	logData.Timestamp = args.Timestamp
	logData.ExpiresAt = time.Unix(args.ExpiresAt, 0)
	logData.ExpiresAtShow = args.ExpiresAt
	logData.Tags = args.Tags
	logData.AdditionalData = args.AdditionalData

	if InsertMGO(collection, logData) != nil {
		log.Println("InsertMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert Problems"}
	}

	reply.LogId=args.ID.Hex()
	return nil
}


func (h *Log) Get(r *http.Request, args *httpmodels.RequestLogGetLog, reply *httpmodels.ResponseLogGet) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()


	reply.LogList = GetFromMGO(collection, args.SearchFilter, args.Limit, args.Offset, args.Sort)

	return nil
}


func (h *Log) GetCount(r *http.Request, args *httpmodels.RequestLogGetCount, reply *httpmodels.ResponseLogGetCount) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	count, err := GetCountMGO(collection, args.SearchFilter)
	if err != nil {
		log.Println("GetCountMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Data Problems"}
	}

	reply.LogCount = count
	return nil
}



func (h *Log) GetCategories(r *http.Request, args *struct{}, reply *httpmodels.ResponseLogGetCategories) error {

	serverList := GetServersList()

	for _, srv := range serverList {
		collections, err := GetCollectionsList(srv, GetUser(r))
		if err != nil {
			log.Println("GetCollectionsList: "+err.Error())
			return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Categories Getting Problems"}
		}
		for _,col := range collections {
			reply.CategoriesList = append(reply.CategoriesList, col)
		}
	}

	return nil
}


func (h *Log) Remove(r *http.Request, args *httpmodels.RequestLogRemoveLog, reply *httpmodels.ResponseLogRemoveLog) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()


	info, err := RemoveAllFromMGO(collection, args.SearchFilter)
	if err != nil {
		log.Println("RemoveAllFromMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}

	reply.Matched = info.Matched
	reply.Removed = info.Removed
	return nil
}


func (h *Log) RemoveCategory(r *http.Request, args *httpmodels.RequestLogRemoveCategory, reply *httpmodels.ResponseLogRemoveCategory) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()


	err = collection.DropCollection()
	if err != nil {
		log.Println("DropCollection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}
	reply.Success = 1
	return nil
}




func (h *Log) Transfer(r *http.Request, args *httpmodels.RequestLogTransferLog, reply *httpmodels.ResponseLogTransferLog) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	//CONNECT TO "FROM" DB AND GET DATA
	connectionStringFrom := GetConnectionStringByCategory(args.OldCategory)

	sessionFrom, err := CreateMGOConnection(connectionStringFrom)
	defer sessionFrom.Close()
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collectionFrom, err := UseMGOCol(UseMGODB(sessionFrom, GetUser(r)), args.OldCategory)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}

	found := GetFromMGO(collectionFrom, args.SearchFilter, 0, 0, nil)


	//CONNECT TO "TO" DB
	connectionStringTo := GetConnectionStringByCategory(args.NewCategory)
	sessionTo, err := CreateMGOConnection(connectionStringTo)
	defer sessionTo.Close()
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collectionTo, err := UseMGOCol(UseMGODB(sessionTo, GetUser(r)), args.NewCategory)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}


	//TRANSFERING DATA "FROM" -> "TO"
	for _, element := range found {
		err = collectionFrom.RemoveId(element.ID)
		if err != nil {
			return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Removing from Collection Problems: " + err.Error()}
		}
		err = collectionTo.Insert(element)
		if err != nil {
			return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert to Collection Problems: " + err.Error()}
		}
		reply.TransferedLogId = append(reply.TransferedLogId, element.ID)
	}

	return nil
}





func (h *Log) ModifyTTL(r *http.Request, args *httpmodels.RequestLogModifyTTL, reply *httpmodels.ResponseLogModifyTTL) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	connectionString := GetConnectionStringByCategory(args.Category)
	session, err := CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := UseMGOCol(UseMGODB(session, GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	updateData := bson.M{"$set": bson.M{"expiresAt": time.Unix(args.NewTTL, 0), "expiresAtIntJustToShow": args.NewTTL}}

	info, err := UpdateAllMGO(collection, args.SearchFilter, updateData)
	if err != nil {
		log.Println("UpdateAllMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}

	reply.Matched = info.Matched
	reply.Updated = info.Updated
	reply.UpsertedId = info.UpsertedId

	return nil
}




/*
-------------------------------------------------
INTERNAL CALL: Some general system tests, checking/showing how errors and other conditions are handled
-------------------------------------------------
*/

type System struct{}

func (h *System) Test(r *http.Request, args *struct{ Test string }, reply *struct{ Answer string }) error {
	log.Println(GetUser(r))
	if args.Test == "" {
		return &json2.Error{Code: json2.E_INVALID_REQ, Message: "Missing required parameter: Test"}
	}
	if args.Test == "fatal" {
		// fatal, programming error
		x := 0
		y := 0
		x = x / y
	}
	reply.Answer = "Hello, "+args.Test
	return nil
}



func (h *System) GetCacheAll(r *http.Request, args *struct{}, reply *struct{ Count int
									     Items map[string]cache.Item }) error {
	reply.Count = Cache.ItemCount()
	reply.Items = Cache.Items()
	return nil
}

func (h *System) GetCache(r *http.Request, args *struct{Key string}, reply *struct{ Key string
										    Data interface{} }) error {
	var found bool
	reply.Key = args.Key
	reply.Data, found = Cache.Get(args.Key)
	if found == false {
		return json2.ErrNullResult
	}
	return nil
}

func (h *System) DeleteCache(r *http.Request, args *struct{Key string}, reply *struct{ Status int }) error {
	Cache.Delete(args.Key)
	reply.Status = 1
	return nil
}

func (h *System) SetCache(r *http.Request, args *struct{Key string
							Data interface{}
							Ttl time.Duration}, reply *struct{ Status int }) error {
	Cache.Set(args.Key, args.Data, args.Ttl)
	reply.Status = 1
	return nil
}

func (h *System) FlushCache(r *http.Request, args *struct{}, reply *struct{ Status int }) error {
	Cache.Flush()
	reply.Status = 1
	return nil
}
package api

import (
	"net/http"
	"log"
	"time"
	"fmt"

	logReq "../models/request/log"
	logResp "../models/response/log"
	mongoModel "../models/mongo"

	"../helpers/auth"
	"../helpers/mgolibs"
	"../helpers"

	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/rpc/v2/json2"
)


type Log struct{}

type GeneralArgs struct{}
type GeneralReply struct{Warnings []string `json:",omitempty"`}



func (h *Log) Add(r *http.Request, args *logReq.Add, reply *logResp.Add) error {

	//res2B, _ := json.Marshal(args)
	//log.Println(string(res2B))
	errs := args.Validate()
	if errs != nil {
		return errs
	}

	connectionString := mgolibs.GetConnectionStringByCategory(args.Category)
	session, err := mgolibs.CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("[" + helpers.GetFuncName(1) + "] CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(session, auth.GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	args.ID = bson.NewObjectId()

	logData := mongoModel.Log{}
	logData.ID = args.ID
	logData.Category = args.Category
	logData.Level = args.Level
	logData.Message = args.Message
	logData.Timestamp = args.Timestamp
	logData.ExpiresAt = time.Unix(args.ExpiresAt, 0)

	if mgolibs.InsertMGO(collection, logData) != nil {
		log.Println("InsertMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert Problems"}
	}

	reply.LogId=args.ID.Hex()
	return nil
}



func (h *Log) AddCustom(r *http.Request, args *logReq.AddCustom, reply *logResp.Add) error {

	//res2B, _ := json.Marshal(args)
	//log.Println(string(res2B))
	errs := args.Validate()
	if errs != nil {
		return errs
	}

	connectionString := mgolibs.GetConnectionStringByCategory(args.Category)
	session, err := mgolibs.CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(session, auth.GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	args.ID = bson.NewObjectId()

	logData := mongoModel.CustomLog{}
	logData.ID = args.ID
	logData.Category = args.Category
	logData.Level = args.Level
	logData.Message = args.Message
	logData.Timestamp = args.Timestamp
	logData.ExpiresAt = time.Unix(args.ExpiresAt, 0)
	logData.Tags = args.Tags
	logData.AdditionalData = args.AdditionalData

	if mgolibs.InsertMGO(collection, logData) != nil {
		log.Println("InsertMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert Problems"}
	}

	reply.LogId=args.ID.Hex()
	return nil
}


func (h *Log) Get(r *http.Request, args *logReq.GetLog, reply *logResp.Get) error {

	//res2B, _ := json.Marshal(args)
	//log.Println(string(res2B))
	errs := args.Validate()
	if errs != nil {
		return errs
	}

	connectionString := mgolibs.GetConnectionStringByCategory(args.Category)
	session, err := mgolibs.CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(session, auth.GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()


	reply.LogList = mgolibs.GetFromMGO(collection, args.Search_filter, 10)

	return nil
}


func (h *Log) GetCount(r *http.Request, args *logReq.GetCount, reply *struct{LogCount int}) error {

	//res2B, _ := json.Marshal(args)
	//log.Println(string(res2B))
	errs := args.Validate()
	if errs != nil {
		return errs
	}


	connectionString := mgolibs.GetConnectionStringByCategory(args.Category)
	session, err := mgolibs.CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(session, auth.GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()

	count, err := mgolibs.GetCountMGO(collection, args.Search_filter)
	if err != nil {
		log.Println("GetCountMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Data Problems"}
	}

	reply.LogCount = count
	return nil
}



func (h *Log) GetCategories(r *http.Request, args *struct{}, reply *logResp.GetCategories) error {

	//res2B, _ := json.Marshal(args)
	//log.Println(string(res2B))

	serverList := mgolibs.GetServersList()

	for _, srv := range serverList {
		collections, err := mgolibs.GetCollectionsList(srv, auth.GetUser(r))
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


func (h *Log) Remove(r *http.Request, args *logReq.RemoveLog, reply *logResp.RemoveLog) error {

	//res2B, _ := json.Marshal(args)
	//log.Println(string(res2B))
	errs := args.Validate()
	if errs != nil {
		return errs
	}

	connectionString := mgolibs.GetConnectionStringByCategory(args.Category)
	session, err := mgolibs.CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(session, auth.GetUser(r)), args.Category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	defer session.Close()


	info, err := mgolibs.RemoveAllFromMGO(collection, args.Search_filter)
	if err != nil {
		log.Println("RemoveAllFromMGO: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}

	reply.Matched = info.Matched
	reply.Removed = info.Removed
	return nil
}


func (h *Log) RemoveCategory(r *http.Request, args *logReq.RemoveCategory, reply *logResp.RemoveCategory) error {

	errs := args.Validate()
	if errs != nil {
		return errs
	}

	connectionString := mgolibs.GetConnectionStringByCategory(args.Category)
	session, err := mgolibs.CreateMGOConnection(connectionString)
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collection, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(session, auth.GetUser(r)), args.Category)
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




func (h *Log) Transfer(r *http.Request, args *logReq.TransferLog, reply *logResp.TransferLog) error {

	errs := args.Validate()
	if errs != nil {
		return errs
	}


	//CONNECT TO "FROM" DB AND GET DATA
	connectionStringFrom := mgolibs.GetConnectionStringByCategory(args.Old_category)

	sessionFrom, err := mgolibs.CreateMGOConnection(connectionStringFrom)
	defer sessionFrom.Close()
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collectionFrom, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(sessionFrom, auth.GetUser(r)), args.Old_category)
	if err != nil {
		log.Println("UseMGOCol: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}

	found := mgolibs.GetFromMGO(collectionFrom, args.Search_filter, 0)


	//CONNECT TO "TO" DB
	connectionStringTo := mgolibs.GetConnectionStringByCategory(args.New_category)
	sessionTo, err := mgolibs.CreateMGOConnection(connectionStringTo)
	defer sessionTo.Close()
	if err != nil {
		log.Println("CreateMGOConnection: "+err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}
	collectionTo, err := mgolibs.UseMGOCol(mgolibs.UseMGODB(sessionTo, auth.GetUser(r)), args.New_category)
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




/*
func (h *Log) Modify(r *http.Request, args *struct{}, reply *CurrencyAnswer) error {
	return nil
}
*/
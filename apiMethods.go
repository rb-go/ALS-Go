package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Riftbit/ALS-Go/httpmodels"
	"github.com/Riftbit/ALS-Go/mongomodels"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func makeDBConnection(db string, category string) (*mgo.Session, *mgo.Collection, *json2.Error) {
	connectionString := getConnectionStringByCategory(category)

	session, err := createMGOConnection(connectionString)
	if err != nil {
		Logger.Error("[" + getFuncName(1) + "] createMGOConnection: " + err.Error())
		return session, nil, &json2.Error{Code: json2.E_INTERNAL, Message: "Log Connection Problems"}
	}

	collection, err := useMGOCol(useMGODB(session, db), category)
	if err != nil {
		Logger.Error("useMGOCol: " + err.Error())
		return session, collection, &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Collection Problems"}
	}
	return session, collection, nil
}

//Log area
type Log struct{}

//Add Method to add Log
func (h *Log) Add(r *http.Request, args *httpmodels.RequestLogAdd, reply *httpmodels.ResponseLogAdd) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
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

	if err := insertMGO(collection, logData); err != nil {
		Logger.Error("InsertMGO: " + err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert Problems"}
	}

	reply.LogID = args.ID.Hex()
	return nil
}

//AddCustom Method to add Log with additional params
func (h *Log) AddCustom(r *http.Request, args *httpmodels.RequestLogAddCustom, reply *httpmodels.ResponseLogAdd) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
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

	if err := insertMGO(collection, logData); err != nil {
		Logger.Error("InsertMGO: " + err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert Problems"}
	}

	reply.LogID = args.ID.Hex()
	return nil
}

//Get Method to get Log/Logs
func (h *Log) Get(r *http.Request, args *httpmodels.RequestLogGetLog, reply *httpmodels.ResponseLogGet) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
	}
	defer session.Close()

	reply.LogList = getFromMGO(collection, args.SearchFilter, args.Limit, args.Offset, args.Sort)

	return nil
}

//GetCount Method to get Log counts
func (h *Log) GetCount(r *http.Request, args *httpmodels.RequestLogGetCount, reply *httpmodels.ResponseLogGetCount) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
	}
	defer session.Close()

	count, err := getCountMGO(collection, args.SearchFilter)
	if err != nil {
		Logger.Error("GetCountMGO: " + err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Select Data Problems"}
	}

	reply.LogCount = count
	return nil
}

//GetCategories Method to get Log categories
func (h *Log) GetCategories(r *http.Request, args *struct{}, reply *httpmodels.ResponseLogGetCategories) error {

	serverList := getServersList()

	for _, srv := range serverList {
		collections, err := getCollectionsList(srv, getUser(r))
		if err != nil {
			Logger.Error("GetCollectionsList: " + err.Error())
			return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Categories Getting Problems"}
		}
		for _, col := range collections {
			reply.CategoriesList = append(reply.CategoriesList, col)
		}
	}

	return nil
}

//Remove Method to remove Log/Logs
func (h *Log) Remove(r *http.Request, args *httpmodels.RequestLogRemoveLog, reply *httpmodels.ResponseLogRemoveLog) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
	}
	defer session.Close()

	info, err := removeAllFromMGO(collection, args.SearchFilter)
	if err != nil {
		Logger.Error("RemoveAllFromMGO: " + err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}

	reply.Matched = info.Matched
	reply.Removed = info.Removed
	return nil
}

//RemoveCategory Method to remove Log category
func (h *Log) RemoveCategory(r *http.Request, args *httpmodels.RequestLogRemoveCategory, reply *httpmodels.ResponseLogRemoveCategory) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
	}
	defer session.Close()

	if err := collection.DropCollection(); err != nil {
		Logger.Error("DropCollection: " + err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}
	reply.Success = 1
	return nil
}

//Transfer Method to transfer Log/Logs to another category
func (h *Log) Transfer(r *http.Request, args *httpmodels.RequestLogTransferLog, reply *httpmodels.ResponseLogTransferLog) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	//CONNECT TO "FROM" DB AND GET DATA
	sessionFrom, collectionFrom, errjs := makeDBConnection(getUser(r), args.OldCategory)
	if errjs != nil {
		return errjs
	}
	defer sessionFrom.Close()

	found := getFromMGO(collectionFrom, args.SearchFilter, -1, 0, nil)

	fmt.Println(printObject(found))

	//CONNECT TO "TO" DB
	sessionTo, collectionTo, errjs := makeDBConnection(getUser(r), args.NewCategory)
	if errjs != nil {
		return errjs
	}
	defer sessionTo.Close()

	//TRANSFERING DATA "FROM" -> "TO"
	for _, element := range found {

		if err := collectionFrom.RemoveId(element.ID); err != nil {
			return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Removing from Collection Problems: " + err.Error()}
		}
		if err := collectionTo.Insert(element); err != nil {
			return &json2.Error{Code: json2.E_INTERNAL, Message: "Log Insert to Collection Problems: " + err.Error()}
		}
		reply.TransferedLogID = append(reply.TransferedLogID, element.ID)
	}

	return nil
}

// ModifyTTL Method to modify Log/Logs TTL value
func (h *Log) ModifyTTL(r *http.Request, args *httpmodels.RequestLogModifyTTL, reply *httpmodels.ResponseLogModifyTTL) error {

	errs := args.Validate()
	if errs != nil {
		return &json2.Error{Code: json2.E_BAD_PARAMS, Message: errs.Error()}
	}

	session, collection, errjs := makeDBConnection(getUser(r), args.Category)
	if errjs != nil {
		return errjs
	}
	defer session.Close()

	updateData := bson.M{"$set": bson.M{"expiresAt": time.Unix(args.NewTTL, 0), "expiresAtIntJustToShow": args.NewTTL}}

	info, err := updateAllMGO(collection, args.SearchFilter, updateData)
	if err != nil {
		Logger.Error("UpdateAllMGO: " + err.Error())
		return &json2.Error{Code: json2.E_INTERNAL, Message: fmt.Sprintf("%s", err)}
	}

	reply.Matched = info.Matched
	reply.Updated = info.Updated
	reply.UpsertedID = info.UpsertedId

	return nil
}

/*
-------------------------------------------------
INTERNAL CALL: Some general system tests, checking/showing how errors and other conditions are handled
-------------------------------------------------
*/

//System ...
type System struct{}

// GetCacheAll cache records
func (h *System) GetCacheAll(r *http.Request, args *struct{}, reply *struct {
	Count int
	Items map[string]cache.Item
}) error {
	reply.Count = Cache.ItemCount()
	reply.Items = Cache.Items()
	return nil
}

// GetCache record by key
func (h *System) GetCache(r *http.Request, args *struct{ Key string }, reply *struct {
	Key  string
	Data interface{}
}) error {
	var found bool
	reply.Key = args.Key
	reply.Data, found = Cache.Get(args.Key)
	if found == false {
		return json2.ErrNullResult
	}
	return nil
}

// DeleteCache record by key
func (h *System) DeleteCache(r *http.Request, args *struct{ Key string }, reply *struct{ Status int }) error {
	Cache.Delete(args.Key)
	reply.Status = 1
	return nil
}

// SetCache for cache element
func (h *System) SetCache(r *http.Request, args *struct {
	Key  string
	Data interface{}
	TTL  time.Duration
}, reply *struct{ Status int }) error {
	Cache.Set(args.Key, args.Data, args.TTL)
	reply.Status = 1
	return nil
}

// FlushCache clean all cached data
func (h *System) FlushCache(r *http.Request, args *struct{}, reply *struct{ Status int }) error {
	Cache.Flush()
	reply.Status = 1
	return nil
}

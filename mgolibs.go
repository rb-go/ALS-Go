package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/riftbit/ALS-Go/mongomodels"
	"github.com/tmc/mgohacks"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func createMGOConnection(connectionString string) (*mgo.Session, error) {
	session, err := mgo.DialWithTimeout(connectionString, Configs.Mongo.ConnectionTimeout*time.Millisecond)
	if err != nil {
		return nil, err
	}
	Logger.Debug("[createMGOConnection] Selected connection string: ", connectionString)
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func useMGODB(session *mgo.Session, dbName string) *mgo.Database {
	return session.DB(dbName)

}

func useMGOCol(mgDB *mgo.Database, category string) (*mgo.Collection, error) {
	c := mgDB.C(category)
	//Set EnsureIndex if we need it
	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", mgDB.Name, category)
	indexValue, existIndexCache := Cache.Get(cahceKey)
	Logger.Debug("[useMGOCol] Check cache for ", cahceKey)
	if existIndexCache == false || indexValue != "1" {
		index := mgo.Index{
			Key:         []string{"expiresAt"},
			Background:  true,
			ExpireAfter: 0,
		}
		err := mgohacks.EnsureTTLIndex(c, index)
		if err != nil {
			return nil, err
		}
		Cache.Set(cahceKey, "1", cache.NoExpiration)
		Logger.Debug("[useMGOCol] Set cache for ", cahceKey)
	}

	//Set Index on Timestamp field if we need it
	cahceKey = fmt.Sprintf("MongoDB:TimestampIndex:%s:%s", mgDB.Name, category)
	indexValue, existIndexCache = Cache.Get(cahceKey)
	Logger.Debug("[useMGOCol] Check cache for ", cahceKey)
	if existIndexCache == false || indexValue != "1" {
		index := mgo.Index{
			Key:        []string{"timestamp"},
			Background: true,
		}
		err := c.EnsureIndex(index)
		if err != nil {
			return nil, err
		}
		Cache.Set(cahceKey, "1", cache.NoExpiration)
		Logger.Debug("[useMGOCol] Set cache for ", cahceKey)
	}

	return c, nil
}

func insertMGO(c *mgo.Collection, args interface{}) error {
	Logger.Debug("[insertMGO] args: ", printObject(args))
	err := c.Insert(&args)
	return err
}

func getFromMGO(c *mgo.Collection, searchFilter map[string]interface{}, limit int, offset int, sortBy []string) []mongomodels.MongoCustomLog {
	var results []mongomodels.MongoCustomLog
	searchFilter = prepareSearchFilter(searchFilter)
	Logger.Debug("[getFromMGO] searchFilter: ", printObject(searchFilter))
	if limit < 0 {
		c.Find(&searchFilter).Sort(sortBy...).Skip(offset).All(&results)
	} else {
		c.Find(&searchFilter).Sort(sortBy...).Skip(offset).Limit(limit).All(&results)
	}
	Logger.Debug("[getFromMGO] results: ", printObject(results))
	n, _ := c.Count()
	Logger.Debug("[getFromMGO] count: ", n)

	for i := range results {
		if results[i].ExpiresAtShow == 0 {
			results[i].ExpiresAtShow = results[i].ExpiresAt.Unix()
		}
	}
	return results
}

func removeAllFromMGO(c *mgo.Collection, searchFilter map[string]interface{}) (*mgo.ChangeInfo, error) {
	searchFilter = prepareSearchFilter(searchFilter)
	Logger.Debug("[removeAllFromMGO] searchFilter: ", printObject(searchFilter))
	return c.RemoveAll(&searchFilter)
}

func updateAllMGO(c *mgo.Collection, searchFilter map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	searchFilter = prepareSearchFilter(searchFilter)
	Logger.Debug("[updateAllMGO] searchFilter: ", printObject(searchFilter))
	return c.UpdateAll(searchFilter, update)
}

func getCountMGO(c *mgo.Collection, args map[string]interface{}) (int, error) {
	count, err := c.Find(args).Count()
	Logger.Debug("[getCountMGO] args: ", printObject(args))
	return count, err
}

func getConnectionStringByCategory(category string) string {
	connData := mGOadditionalCollectionsConn[category]
	if connData == "" {
		connData = Configs.Mongo.CommonDB.ConnectionString
	}
	return connData
}

func getServersList() []string {
	servers := make(map[string]int)
	var result []string
	for _, coll := range mGOadditionalCollectionsConn {
		servers[coll] = 1
	}
	servers[Configs.Mongo.CommonDB.ConnectionString] = 1
	for srv := range servers {
		result = append(result, srv)
	}
	Logger.Debug("[getServersList] ", result)
	return result
}

func getCollectionsList(server string, dbName string) ([]string, error) {
	connection, err := createMGOConnection(server)
	if err != nil {
		return nil, err
	}
	db := useMGODB(connection, dbName)
	collectionNames, err := db.CollectionNames()
	if err != nil {
		return nil, err
	}
	Logger.Debug("[getCollectionsList] ", collectionNames)
	return collectionNames, nil
}

func prepareSearchFilter(searchFilter map[string]interface{}) map[string]interface{} {
	for key := range searchFilter {
		if key == "_id" {
			findIDKeyValuesAndFixThem(&searchFilter, key)
		}
		if md, ok := searchFilter[key].(map[string]interface{}); ok {
			prepareSearchFilter(md)
		}
		if md, ok := searchFilter[key].([]interface{}); ok {
			for _, v := range md {
				if reflect.TypeOf(v).Kind() == reflect.Map || reflect.TypeOf(v).Kind() == reflect.Slice {
					prepareSearchFilter(v.(map[string]interface{}))
				}
			}
		}
	}
	return searchFilter
}

func findIDKeyValuesAndFixThem(searchFilter *map[string]interface{}, baseKey string) {
	data := *searchFilter
	currentType := reflect.TypeOf(data[baseKey]).Kind()
	if currentType == reflect.String {
		data[baseKey] = bson.ObjectIdHex(data[baseKey].(string))
	}

	if currentType == reflect.Map {
		for k := range data[baseKey].(map[string]interface{}) {
			newlevel := data[baseKey].(map[string]interface{})
			findIDKeyValuesAndFixThem(&newlevel, k)
		}
	}
	if currentType == reflect.Slice {
		for k, v := range data[baseKey].([]interface{}) {
			data[baseKey].([]interface{})[k] = bson.ObjectIdHex(v.(string))
		}
	}
}

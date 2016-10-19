package main

import (
	"github.com/tmc/mgohacks"
	"gopkg.in/mgo.v2"
	"time"
	"reflect"
	"gopkg.in/mgo.v2/bson"
	"github.com/patrickmn/go-cache"
	"fmt"
	"github.com/Riftbit/ALS-Go/mongomodels"
)

func createMGOConnection(connectionString string) (*mgo.Session, error ){
	session, err := mgo.DialWithTimeout(connectionString, Configs.Mongo.ConnectionTimeout*time.Millisecond)
	if err != nil {
		return nil, err
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func useMGODB(session *mgo.Session, dbName string) *mgo.Database {
	return session.DB(dbName)

}

func useMGOCol(mgDB *mgo.Database, category string) (*mgo.Collection, error) {
	c := mgDB.C(category)
	cahceKey := fmt.Sprintf("MongoDB:EnsureIndex:%s:%s", mgDB.Name, category)
	indexValue, existIndexCache := Cache.Get(cahceKey)
	if existIndexCache == false || indexValue != "1" {
		index := mgo.Index{
			Key:        []string{"expiresAt"},
			Background: true,
			ExpireAfter: 0,
		}
		err := mgohacks.EnsureTTLIndex(c, index)
		if err != nil {
			return nil, err
		}
		Cache.Set(cahceKey, "1", cache.NoExpiration)
	}
	return c, nil
}

func insertMGO(c *mgo.Collection, args interface{}) error{
	err := c.Insert(&args)
	return err
}


func getFromMGO(c *mgo.Collection, searchFilter map[string]interface{}, limit int, offset int, sortBy []string) []mongomodels.MongoCustomLog {
	var results []mongomodels.MongoCustomLog
	searchFilter = prepareSearchFilter(searchFilter)
	if limit < 0 {
		c.Find(&searchFilter).Sort(sortBy...).Skip(offset).All(&results)
	} else {
		c.Find(&searchFilter).Sort(sortBy...).Skip(offset).Limit(limit).All(&results)
	}
	for i := range results {
		if(results[i].ExpiresAtShow == 0) {
			results[i].ExpiresAtShow = results[i].ExpiresAt.Unix()
		}
	}
	return results
}

func removeAllFromMGO(c *mgo.Collection, searchFilter map[string]interface{}) (*mgo.ChangeInfo, error){
	searchFilter = prepareSearchFilter(searchFilter)
	return c.RemoveAll(&searchFilter)
}


func updateAllMGO(c *mgo.Collection, searchFilter map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	searchFilter = prepareSearchFilter(searchFilter)
	return c.UpdateAll(searchFilter, update)
}

func getCountMGO(c *mgo.Collection, args map[string]interface{}) (int, error) {
	count, err := c.Find(args).Count()
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
	return collectionNames, nil
}


func prepareSearchFilter(searchFilter map[string]interface{}) map[string]interface{} {
	for key := range searchFilter {
		if key == "_id" {
			findIDKeyValuesAndFixThem(&searchFilter, key)
		}
		if md, ok := searchFilter[key].(map[string]interface{}) ; ok {
			prepareSearchFilter(md)
		}
		if md, ok := searchFilter[key].([]interface{}) ; ok {
			for _,v := range md {
				if reflect.TypeOf(v).Kind() == reflect.Map || reflect.TypeOf(v).Kind() == reflect.Slice {
					prepareSearchFilter(v.(map[string]interface{}))
				}
			}
		}
	}
	return searchFilter
}

func findIDKeyValuesAndFixThem(searchFilter *map[string]interface{}, baseKey string) {
	if baseKey == "" {
		baseKey = "_id"
	}
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
		for k,v := range data[baseKey].([]interface{}) {
			data[baseKey].([]interface {})[k] = bson.ObjectIdHex(v.(string))
		}
	}
}
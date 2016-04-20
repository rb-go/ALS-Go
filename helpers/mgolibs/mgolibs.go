package mgolibs

import (
	"github.com/tmc/mgohacks"
	"gopkg.in/mgo.v2"
	"../../configs"
	"../../models/mongo"
	"time"
)

func CreateMGOConnection(connectionString string) (*mgo.Session, error ){
	session, err := mgo.DialWithTimeout(connectionString, configs.Configs.Mongo.ConnectionTimeout*time.Millisecond)
	if err != nil {
		return nil, err
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func UseMGODB(session *mgo.Session, dbName string) *mgo.Database {
	return session.DB(dbName)

}

func UseMGOCol(mgDB *mgo.Database, category string) (*mgo.Collection, error) {
	c := mgDB.C(category)
	index := mgo.Index{
		Key:        []string{"expiresAt"},
		Background: true,
		ExpireAfter: 0,
	}

	err := mgohacks.EnsureTTLIndex(c, index)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func InsertMGO(c *mgo.Collection, args interface{}) error{
	err := c.Insert(&args)
	return err
}


func GetFromMGO(c *mgo.Collection, args map[string]interface{}, limit int, offset int, sortBy []string) []mongo.CustomLog {
	var results []mongo.CustomLog
	query := c.Find(&args)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Skip(limit)
	}
	if sortBy != nil {
		query = query.Sort(sortBy...)
	}

	query.All(&results)
	return results
}

func RemoveAllFromMGO(c *mgo.Collection, args map[string]interface{}) (*mgo.ChangeInfo, error){
	info, err := c.RemoveAll(&args)
	return info, err
}

func RemoveFromMGO(c *mgo.Collection, args map[string]interface{}) error {
	err := c.Remove(&args)
	return err
}

func GetCountMGO(c *mgo.Collection, args map[string]interface{}) (int, error) {
	count, err := c.Find(args).Count()
	return count, err
}

func GetConnectionStringByCategory(category string) string {
	connData := configs.MGOadditionalCollectionsConn[category]
	if connData == "" {
		connData = configs.Configs.Mongo.CommonDB.ConnectionString
	}
	return connData
}


func GetServersList() []string {
	servers := make(map[string]int)
	var result []string
	for _, coll := range configs.MGOadditionalCollectionsConn {
		servers[coll] = 1
	}
	servers[configs.Configs.Mongo.CommonDB.ConnectionString] = 1
	for srv, _ := range servers {
		result = append(result, srv)
	}

	return result
}

func GetCollectionsList(server string, dbName string) ([]string, error) {
	connection, err := CreateMGOConnection(server)
	if err != nil {
		return nil, err
	}
	db := UseMGODB(connection, dbName)
	collectionNames, err := db.CollectionNames()
	if err != nil {
		return nil, err
	}
	return collectionNames, nil
}
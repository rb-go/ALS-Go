package main

import (
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache.Cache

type MongoCommonServerConf struct {
	ConnectionString string  `yaml:"connection"`
}

type MongoAdditionalServerConf struct {
	ConnectionString string  `yaml:"connection"`
	Collections  []string `yaml:"collections"`
}

type Conf struct {

	System struct {
	       		MaxThreads int  `yaml:"maxThreads"`
	       		ListenOn  string `yaml:"listenOn"`
		}

	Admin struct {
		      RootUser string `yaml:"rootUser"`
		      RootPassword string `yaml:"rootPassword"`
		      RootEmail string `yaml:"rootEmail"`
	      }

	Db struct {
		   DbConnectionString string `yaml:"dbConnectionString"`
	   }

	Mongo struct {
		      ConnectionTimeout time.Duration  `yaml:"connectionTimeout"`
		      CommonDB MongoCommonServerConf `yaml:"commonDB"`
		      AdditionalDB []MongoAdditionalServerConf `yaml:"additionalDB"`
	      }
}

var Configs Conf
var ConfigPath string

var DBConn *gorm.DB
var MGOadditionalCollectionsConn map[string]string

func ProcessMGOAdditionalConf() {
	MGOadditionalCollectionsConn = make(map[string]string)
	if len(Configs.Mongo.AdditionalDB) > 0 {
		for _, additDB := range Configs.Mongo.AdditionalDB {
			if len(additDB.Collections) > 0 {
				for _, coll := range additDB.Collections {
					MGOadditionalCollectionsConn[coll] = additDB.ConnectionString
				}
			}
		}
	}
}


func IsDBConnected() bool {
	err := DBConn.DB().Ping()
	if err != nil {
		return false
	} else {
		return true
	}
}
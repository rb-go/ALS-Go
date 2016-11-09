package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

//Cache ...
var Cache *cache.Cache

//Logger ...
var Logger *logrus.Logger

type mongoCommonServerConf struct {
	ConnectionString string `yaml:"connection"`
}

type mongoAdditionalServerConf struct {
	ConnectionString string   `yaml:"connection"`
	Collections      []string `yaml:"collections"`
}

type conf struct {
	System struct {
		MaxThreads int    `yaml:"maxThreads"`
		ListenOn   string `yaml:"listenOn"`
	}
	Admin struct {
		RootUser     string `yaml:"rootUser"`
		RootPassword string `yaml:"rootPassword"`
		RootEmail    string `yaml:"rootEmail"`
	}
	Db struct {
		DbType             string `yaml:"dbType"`
		DbConnectionString string `yaml:"dbConnectionString"`
	}
	Log struct {
		Formatter       string `yaml:"formatter"` //text, json
		LogLevel        string `yaml:"logLevel"`  // panic, fatal, error, warn, warning, info, debug
		DisableColors   bool   `yaml:"disableColors"`
		TimestampFormat string `yaml:"timestampFormat"`
	}
	Mongo struct {
		ConnectionTimeout time.Duration               `yaml:"connectionTimeout"`
		CommonDB          mongoCommonServerConf       `yaml:"commonDB"`
		AdditionalDB      []mongoAdditionalServerConf `yaml:"additionalDB"`
	}
}

//Configs ...
var Configs conf

var configPath string

//DBConn ...
var DBConn *gorm.DB

var mGOadditionalCollectionsConn map[string]string

func processMGOAdditionalConf() {
	mGOadditionalCollectionsConn = make(map[string]string)
	if len(Configs.Mongo.AdditionalDB) > 0 {
		for _, additDB := range Configs.Mongo.AdditionalDB {
			if len(additDB.Collections) > 0 {
				for _, coll := range additDB.Collections {
					mGOadditionalCollectionsConn[coll] = additDB.ConnectionString
				}
			}
		}
	}
}

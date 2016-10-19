package main

import (
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/Sirupsen/logrus"
	"strings"
	"os"
	"fmt"
)

//Cache ...
var Cache *cache.Cache

//Logger ...
var Logger *logrus.Logger

type mongoCommonServerConf struct {
	ConnectionString string  `yaml:"connection"`
}

type mongoAdditionalServerConf struct {
	ConnectionString string  `yaml:"connection"`
	Collections      []string `yaml:"collections"`
}

type conf struct {
	System struct {
		       MaxThreads int  `yaml:"maxThreads"`
		       ListenOn   string `yaml:"listenOn"`
	       }
	Admin  struct {
		       RootUser     string `yaml:"rootUser"`
		       RootPassword string `yaml:"rootPassword"`
		       RootEmail    string `yaml:"rootEmail"`
	       }
	Db     struct {
		       DbConnectionString string `yaml:"dbConnectionString"`
	       }
	Log    struct {
		       Formatter       string `yaml:"formatter"` //text, json
		       LogLevel        string `yaml:"logLevel"`  // panic, fatal, error, warn, warning, info, debug
		       DisableColors   bool `yaml:"disableColors"`
		       TimestampFormat string `yaml:"timestampFormat"`
	       }
	Mongo  struct {
		       ConnectionTimeout time.Duration  `yaml:"connectionTimeout"`
		       CommonDB          mongoCommonServerConf `yaml:"commonDB"`
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

func isDBConnected() bool {
	err := DBConn.DB().Ping()
	if err != nil {
		return false
	}
	return true
}

func initLogger() {

	allowedTimestampsFormat := map[string]int{
		"Mon Jan _2 15:04:05 2006": 1,
		"Mon Jan _2 15:04:05 MST 2006": 1,
		"Mon Jan 02 15:04:05 -0700 2006": 1,
		"02 Jan 06 15:04 MST": 1,
		"02 Jan 06 15:04 -0700": 1,
		"Monday, 02-Jan-06 15:04:05 MST": 1,
		"Mon, 02 Jan 2006 15:04:05 MST": 1,
		"Mon, 02 Jan 2006 15:04:05 -0700": 1,
		"2006-01-02T15:04:05Z07:00": 1,
		"2006-01-02T15:04:05.999999999Z07:00": 1,
		"3:04PM": 1,
		"Jan _2 15:04:05": 1,
		"Jan _2 15:04:05.000": 1,
		"Jan _2 15:04:05.000000": 1,
		"Jan _2 15:04:05.000000000": 1,
	}

	_, ok := allowedTimestampsFormat[Configs.Log.TimestampFormat]
	if ok == false {
		fmt.Println("Wrong Timestamp Format value in config!")
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}


	var formatter logrus.Formatter

	switch strings.ToLower(Configs.Log.Formatter) {
	case "text":
		formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: Configs.Log.DisableColors, TimestampFormat: Configs.Log.TimestampFormat}
		break
	case "json":
		formatter = &logrus.JSONFormatter{TimestampFormat: Configs.Log.TimestampFormat}
		break
	default:
		fmt.Println("Error Log config formatter")
		time.Sleep(1 * time.Second)
		os.Exit(1)
		break
	}

	level, err := logrus.ParseLevel(Configs.Log.LogLevel)
	if err != nil {
		fmt.Println(err)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	Logger = &logrus.Logger{
		Out: os.Stdout,
		Formatter: formatter,
		Level: level,
	}
}

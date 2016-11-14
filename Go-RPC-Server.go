package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/Riftbit/ALS-Go/httpmodels"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	cache "github.com/patrickmn/go-cache"
	validator "gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var applicationExitFunction = func(code int) { os.Exit(code) }
var rpcV2 *rpc.Server

func abstractExitFunction(exit int) {
	applicationExitFunction(exit)
}

func initConfigs() {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		logPrintln(err)
		time.Sleep(10 * time.Millisecond)
		abstractExitFunction(1)
	}

	err = yaml.Unmarshal(data, &Configs)
	if err != nil {
		logPrintln("error reading config", err)
		time.Sleep(10 * time.Millisecond)
		abstractExitFunction(1)
	}
	Cache = cache.New(10*time.Minute, 30*time.Second)
	processMGOAdditionalConf()
}

func initDataBase() {
	var err error
	DBConn, err = gorm.Open(Configs.Db.DbType, Configs.Db.DbConnectionString)
	if err != nil {
		Logger.Errorf("ORM NOT WORKS! - %s", err)
		time.Sleep(10 * time.Millisecond)
		abstractExitFunction(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	if !isDBConnected() {
		Logger.Errorf("DB Connection NOT WORKS! - %s", err.Error())
		time.Sleep(10 * time.Millisecond)
		abstractExitFunction(1)
	} else {
		Logger.Info("DB Connection WORKS!")
		Logger.Info("DB Data and structs initialized!")
	}
}

func initRuntime() {
	numCPU := runtime.NumCPU()
	Logger.Infof("Init runtime to use %d CPUs and %d threads", numCPU, Configs.System.MaxThreads)
	runtime.GOMAXPROCS(numCPU)
	debug.SetMaxThreads(Configs.System.MaxThreads)
}

func parseCommandLineParams() {
	flag.StringVar(&configPath, "config", "./config.yml", "Path to config.yml")
	flag.Parse()
}

func rpcPrepare() {
	rpcV2 = rpc.NewServer()
	rpcV2.RegisterCodec(json2.NewCodec(), "text/plain")
	rpcV2.RegisterCodec(json2.NewCodec(), "application/json")
	rpcV2.RegisterCodec(json2.NewCodec(), "text/plain; charset=utf-8")       // For firefox 11 and other browsers which append the charset=UTF-8
	rpcV2.RegisterCodec(json2.NewCodec(), "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	http.Handle("/", authentificator(rpcV2))
}

func prepareServerWithConfigs() {
	initLogger()
	initRuntime()
	initDataBase()
	rpcPrepare()

	adminMethodsList, basicMethodsList := registerAPI(rpcV2)

	initDatabaseStructure()
	initDatabaseData(adminMethodsList, basicMethodsList)

	validator.SetValidationFunc("CategoryNameValidators", httpmodels.CategoryNameValidator)

	Logger.Infof("Starting server on <%s>", Configs.System.ListenOn)
}

func main() {
	parseCommandLineParams()
	initConfigs()
	prepareServerWithConfigs()
	Logger.Fatal(http.ListenAndServe(Configs.System.ListenOn, nil))
}

func registerAPI(rpcV2 *rpc.Server) ([]string, []string) {
	Logger.Info("Registering exported methods")
	rpcV2.RegisterService(new(Log), "")
	rpcV2.RegisterService(new(System), "")

	var adminMethodsList []string
	var basicMethodsList []string

	// prints a concise summary of the exported API calls
	listMethods := func(m interface{}) {
		typ := reflect.TypeOf(m)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		fooType := reflect.TypeOf(m)
		for i := 0; i < fooType.NumMethod(); i++ {
			method := fooType.Method(i)
			args := reflect.New(method.Type.In(2).Elem()).Elem().Interface()
			resp := reflect.New(method.Type.In(3).Elem()).Elem().Interface()
			Logger.Debugf("request = api.call('%s.%s', %+v) # response: %+v", typ.Name(), method.Name, args, resp)
			if typ.Name() == "System" {
				adminMethodsList = append(adminMethodsList, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			} else {
				basicMethodsList = append(basicMethodsList, fmt.Sprintf("%s.%s", typ.Name(), method.Name))
			}
		}
	}
	Logger.Debug("Start exported methods names")
	listMethods(new(System))
	listMethods(new(Log))
	Logger.Debug("End exporten methods names")
	return adminMethodsList, basicMethodsList
}

func authentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userlogin, password, methodName string
		userlogin, password = parseRequestAuthData(r)
		if checkUserAuth(userlogin, password) {
			rawDataBody := getDataBody(r)
			Logger.Debug("[authentificator] Received Request: ", string(rawDataBody))
			if rawDataBody == nil {
				answerWriter(w, 405, []byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`), true)
			} else {
				jsonData, err := getRequestJSON(rawDataBody)
				if err != nil {
					answerWriter(w, 405, []byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`), true)
				} else {
					methodName = jsonData["method"].(string)
					if checkAPIMethodAccess(userlogin, methodName) == false {
						answerWriter(w, 403, []byte(`{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Permission denied"}, "id": null}`), true)
					} else {
						next.ServeHTTP(w, r)
					}
				}
			}
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="API"`)
			answerWriter(w, 401, []byte(`401 Unauthorized\n`), true)
		}
	})
}

func answerWriter(w http.ResponseWriter, code int, body []byte, isJSON bool) {
	if isJSON {
		w.Header().Set("Content-Type", `application/json; charset=utf-8`)
	}
	w.WriteHeader(code)
	w.Write(body)
}

func getDataBody(r *http.Request) []byte {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	return data
}

func getRequestJSON(data []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	return jsonData, err
}

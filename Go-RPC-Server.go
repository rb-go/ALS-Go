package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"reflect"

	"github.com/Riftbit/ALS-Go/httpmodels"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

var application_exit_function func(code int) = os.Exit

func AbstractExitFunction(exit int) {
	application_exit_function(exit)
}

func initConfigs() {

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(1 * time.Second)
		AbstractExitFunction(1)
	}

	err = yaml.Unmarshal(data, &Configs)
	if err != nil {
		fmt.Println("error reading config", err)
		time.Sleep(1 * time.Second)
		AbstractExitFunction(1)
	}

	initLogger()

	Cache = cache.New(10*time.Minute, 30*time.Second)

	processMGOAdditionalConf()
}

func initDataBase() {
	var err error
	DBConn, err = gorm.Open("mysql", Configs.Db.DbConnectionString)
	if err != nil {
		Logger.Fatalf("ORM NOT WORKS! - %s", err)
		time.Sleep(1 * time.Second)
		AbstractExitFunction(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	if !isDBConnected() {
		Logger.Fatalf("DB Connection NOT WORKS! - %s", err.Error())
		time.Sleep(1 * time.Second)
		AbstractExitFunction(1)
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
	initValidators()
}

func initValidators() {
	validator.SetValidationFunc("CategoryNameValidators", httpmodels.CategoryNameValidator)
}

func parseCommandLineParams() {
	flag.StringVar(&configPath, "c", "./config.yml", "Path to config.yml")
	flag.Parse()
}

func main() {
	parseCommandLineParams()
	initConfigs()
	initRuntime()
	initDataBase()

	rpcV2 := rpc.NewServer()
	rpcV2.RegisterCodec(json2.NewCodec(), "text/plain")
	rpcV2.RegisterCodec(json2.NewCodec(), "application/json")
	rpcV2.RegisterCodec(json2.NewCodec(), "text/plain; charset=utf-8")       // For firefox 11 and other browsers which append the charset=UTF-8
	rpcV2.RegisterCodec(json2.NewCodec(), "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	register(rpcV2)

	http.Handle("/", authentificator(rpcV2))

	Logger.Infof("Starting server on <%s>", Configs.System.ListenOn)
	Logger.Fatal(http.ListenAndServe(Configs.System.ListenOn, nil))
}

func register(rpcV2 *rpc.Server) {
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

	initDatabaseStructure()
	initDatabaseData(adminMethodsList, basicMethodsList)
}

func authentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(r) {
			rawDataBody := getDataBody(r)
			Logger.Debug("[authentificator] Received Request: ", string(rawDataBody))
			if rawDataBody == nil {
				w.Header().Set("Content-Type", `application/json; charset=utf-8`)
				w.WriteHeader(405)
				w.Write([]byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`))
			} else {
				jsonData, err := getRequestJSON(rawDataBody)
				if err != nil {
					w.Header().Set("Content-Type", `application/json; charset=utf-8`)
					w.WriteHeader(405)
					w.Write([]byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`))
				} else {
					if checkAPIMethodAccess(r, jsonData) == false {
						w.Header().Set("Content-Type", `application/json; charset=utf-8`)
						w.WriteHeader(403)
						w.Write([]byte(`{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Permission denied"}, "id": null}`))
					} else {
						next.ServeHTTP(w, r)
					}
				}
			}
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="API"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
		}
	})
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

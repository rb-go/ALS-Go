package main

/*
    test with curl
        curl -X POST -H "Content-Type: application/json" -d '{"method":"HelloService.Say","params":[{"Who":"Test"}], "id":"1"}' http://localhost:10000/api
*/

import (
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v2"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
	"runtime"
	"runtime/debug"
	"flag"
	"io/ioutil"
	"os"
	"encoding/json"
	"bytes"
	"gopkg.in/validator.v2"
	"fmt"
	"github.com/Riftbit/ALS-Go/httpmodels"
)

func initConfigs() {

	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		Logger.Fatal(err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &Configs)
	if err != nil {
		panic(fmt.Sprintf("error reading config: %v", err))
	}

	initLogger()

	DBConn, err = gorm.Open("mysql", Configs.Db.DbConnectionString)
	if err != nil {
		Logger.Fatalf("ORM NOT WORKS! - %s", err)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	if !IsDBConnected() {
		Logger.Fatalf("DB Connection NOT WORKS! - %s", err.Error())
		time.Sleep(1 * time.Second)
		os.Exit(1)
	} else {
		Logger.Info("DB Connection WORKS!")
		Logger.Info("DB Data and structs initialized!")
	}

	Cache = cache.New(5*time.Minute, 30*time.Second)

	ProcessMGOAdditionalConf()
}

func initRuntime() {
	numCpu := runtime.NumCPU()
	Logger.Infof("Init runtime to use %d CPUs and %d threads", numCpu, Configs.System.MaxThreads)
	runtime.GOMAXPROCS(numCpu)
	debug.SetMaxThreads(Configs.System.MaxThreads)
	initValidators()
}

func initValidators() {
	validator.SetValidationFunc("CategoryNameValidators", httpmodels.CategoryNameValidator)
}

func init() {
}

func main() {
	flag.StringVar(&ConfigPath, "c", "./config.yml", "Path to config.yml")
	time.Sleep(1 * time.Second)
	flag.Parse()

	time.Sleep(1 * time.Second)
	initConfigs()
	time.Sleep(1 * time.Second)
	initRuntime()
	time.Sleep(1 * time.Second)

	rpc_v2 := rpc.NewServer()
	rpc_v2.RegisterCodec(json2.NewCodec(), "text/plain")
	rpc_v2.RegisterCodec(json2.NewCodec(), "application/json")
	rpc_v2.RegisterCodec(json2.NewCodec(), "text/plain; charset=utf-8") // For firefox 11 and other browsers which append the charset=UTF-8
	rpc_v2.RegisterCodec(json2.NewCodec(), "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	Register(rpc_v2)

	http.Handle("/", Authentificator(rpc_v2))

	Logger.Infof("Starting server on <%s>", Configs.System.ListenOn)

	Logger.Fatal(http.ListenAndServe(Configs.System.ListenOn, nil))
}

func Authentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if CheckAuth(r) {
			rawDataBody := getDataBody(r)
			if rawDataBody == nil {
				w.Header().Set("Content-Type", `application/json; charset=utf-8`)
				w.WriteHeader(405)
				w.Write([]byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`))
			} else {
				json_data, err := getRequestJson(rawDataBody)
				if err != nil {
					w.Header().Set("Content-Type", `application/json; charset=utf-8`)
					w.WriteHeader(405)
					w.Write([]byte(`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`))
				} else {
					if CheckAPIMethodAccess(r, json_data) == false {
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


type myReader struct {
	*bytes.Buffer
}
func (m myReader) Close() error { return nil }

func getDataBody(r *http.Request) []byte {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	rdr1 := myReader{bytes.NewBuffer(buf)}
	rdr2 := myReader{bytes.NewBuffer(buf)}
	r.Body = rdr2
	data, err := ioutil.ReadAll(rdr1)
	rdr1.Close()
	rdr2.Close()
	if err != nil {
		Logger.Error(err)
		return nil
	}
	return data
}

func getRequestJson(data []byte) (map[string]interface{}, error) {
	var json_data map[string]interface{}
	err := json.Unmarshal(data, &json_data)
	return json_data, err
}
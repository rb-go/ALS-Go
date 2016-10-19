package main

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

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err.Error())
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &Configs)
	if err != nil {
		fmt.Println("error reading config", err)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}

	initLogger()

	DBConn, err = gorm.Open("mysql", Configs.Db.DbConnectionString)
	if err != nil {
		Logger.Fatalf("ORM NOT WORKS! - %s", err)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	if !isDBConnected() {
		Logger.Fatalf("DB Connection NOT WORKS! - %s", err.Error())
		time.Sleep(1 * time.Second)
		os.Exit(1)
	} else {
		Logger.Info("DB Connection WORKS!")
		Logger.Info("DB Data and structs initialized!")
	}

	//it is to fix:
	//[mysql] packets.go:33: unexpected EOF
	//[mysql] packets.go:124: write tcp 127.0.0.1:59804->127.0.0.1:3306: write: broken pipe
	//DBConn.DB().SetMaxIdleConns(0)
	//DBConn.DB().SetMaxOpenConns(100)
	//but now we will try to set in my.cnf wait_timeout=2147483

	Cache = cache.New(10*time.Minute, 30*time.Second)

	processMGOAdditionalConf()
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

func init() {
}

func main() {
	flag.StringVar(&configPath, "c", "./config.yml", "Path to config.yml")
	flag.Parse()
	initConfigs()
	initRuntime()

	rpcV2 := rpc.NewServer()
	rpcV2.RegisterCodec(json2.NewCodec(), "text/plain")
	rpcV2.RegisterCodec(json2.NewCodec(), "application/json")
	rpcV2.RegisterCodec(json2.NewCodec(), "text/plain; charset=utf-8") // For firefox 11 and other browsers which append the charset=UTF-8
	rpcV2.RegisterCodec(json2.NewCodec(), "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	register(rpcV2)

	http.Handle("/", authentificator(rpcV2))

	Logger.Infof("Starting server on <%s>", Configs.System.ListenOn)

	Logger.Fatal(http.ListenAndServe(Configs.System.ListenOn, nil))
}


func authentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(r) {
			rawDataBody := getDataBody(r)
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

func getRequestJSON(data []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	return jsonData, err
}
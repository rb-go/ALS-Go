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
	"log"
	"io/ioutil"
	"os"
	"encoding/json"
	"bytes"
	"gitlab.com/ergoz/ALS-Go/helpers/auth"
	"gitlab.com/ergoz/ALS-Go/app"
	"gitlab.com/ergoz/ALS-Go/configs"
)

func initConfigs() {

	data, err := ioutil.ReadFile(configs.ConfigPath)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &configs.Configs)
	if err != nil {
		log.Printf("error reading config: %v", err)
	}

	configs.DBConn, err = gorm.Open("mysql", configs.Configs.Db.DbConnectionString)
	if err != nil {
		log.Printf("ORM NOT WORKS! - %s", err)
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
	// Open doesn't open a connection. Validate DSN data:
	if !configs.IsDBConnected() {
		log.Printf("DB Connection NOT WORKS! - %s", err.Error())
		time.Sleep(1 * time.Second)
		os.Exit(1)
	} else {
		log.Println("DB Connection WORKS!")
		log.Println("DB Data and structs initialized!")
	}

	configs.Cache = cache.New(5*time.Minute, 30*time.Second)

	configs.ProcessMGOAdditionalConf()
}

func initRuntime() {
	numCpu := runtime.NumCPU()
	log.Printf("Initializing runtime to use %d CPUs and %d threads", numCpu, configs.Configs.System.MaxThreads)
	runtime.GOMAXPROCS(numCpu)
	debug.SetMaxThreads(configs.Configs.System.MaxThreads)
}

func init() {

	log.SetFlags(log.LstdFlags + log.Lshortfile)

	flag.StringVar(&configs.ConfigPath, "-c", "./config.yml", "Path to config.yml without tralling slash at the end, like /etc/als-go")
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
	app.Register(rpc_v2)

	http.Handle("/", Authentificator(rpc_v2))

	log.Printf("Starting server on <%s>", configs.Configs.System.ListenOn)

	log.Fatal(http.ListenAndServe(configs.Configs.System.ListenOn, nil))
}

func main() {
}

func Authentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.CheckAuth(r) {
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
					if auth.CheckAPIMethodAccess(r, json_data) == false {
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
		log.Printf("%s", err)
		return nil
	}
	rdr1 := myReader{bytes.NewBuffer(buf)}
	rdr2 := myReader{bytes.NewBuffer(buf)}
	r.Body = rdr2
	data, err := ioutil.ReadAll(rdr1)
	rdr1.Close()
	rdr2.Close()
	if err != nil {
		log.Printf("%s", err)
		return nil
	}
	return data
}

func getRequestJson(data []byte) (map[string]interface{}, error) {
	var json_data map[string]interface{}
	err := json.Unmarshal(data, &json_data)
	return json_data, err
}
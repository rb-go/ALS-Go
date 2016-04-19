package api
import (
	"../helpers/auth"
	"../configs"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/patrickmn/go-cache"
	"time"
	"log"
	"net/http"
)

/*
-------------------------------------------------
INTERNAL CALL: Some general system tests, checking/showing how errors and other conditions are handled
-------------------------------------------------
*/

type System struct{}

func (h *System) Test(r *http.Request, args *struct{ Test string }, reply *struct{ Answer string }) error {
	log.Println(auth.GetUser(r))
	if args.Test == "" {
		return &json2.Error{Code: json2.E_INVALID_REQ, Message: "Missing required parameter: Test"}
	}
	if args.Test == "fatal" {
		// fatal, programming error
		x := 0
		y := 0
		x = x / y
	}
	reply.Answer = "Hello, "+args.Test
	return nil
}



func (h *System) GetCacheAll(r *http.Request, args *struct{}, reply *struct{ Count int
									     Items map[string]cache.Item }) error {
	reply.Count = configs.Cache.ItemCount()
	reply.Items = configs.Cache.Items()
	return nil
}

func (h *System) GetCache(r *http.Request, args *struct{Key string}, reply *struct{ Key string
										    Data interface{} }) error {
	var found bool
	reply.Key = args.Key
	reply.Data, found = configs.Cache.Get(args.Key)
	if found == false {
		return json2.ErrNullResult
	}
	return nil
}

func (h *System) DeleteCache(r *http.Request, args *struct{Key string}, reply *struct{ Status int }) error {
	configs.Cache.Delete(args.Key)
	reply.Status = 1
	return nil
}

func (h *System) SetCache(r *http.Request, args *struct{Key string
							Data interface{}
							Ttl time.Duration}, reply *struct{ Status int }) error {
	configs.Cache.Set(args.Key, args.Data, args.Ttl)
	reply.Status = 1
	return nil
}

func (h *System) FlushCache(r *http.Request, args *struct{}, reply *struct{ Status int }) error {
	configs.Cache.Flush()
	reply.Status = 1
	return nil
}
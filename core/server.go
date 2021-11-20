package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// global database variable
var tkvdb = make(map[string]interface{})

// private key for server
//var global_private_key []byte

// auth key is key for making database/server connection safer
// it is creating new uuid key
// everytime user is saving
//var Http_security_key uuid.UUID

//var saves = 0

// config
var (
	// used in 119 line in sync_with_servers() function
	// It is optional you can leave it blank
	server_name string
	server_url  = fmt.Sprintf("http://localhost:%s", port)
	// declare servers and child servers names
	// !!!
	// always declare current server name first
	servers_json = map[string]string{
		"node": fmt.Sprintf("http://localhost:%s", port),
		// example of second server
		//"child1": fmt.Sprintf("http://localhost:8894",),
	}
	server_json_path = "./servers.json"

	// port of server
	// you can set it to whatever port you want
	port = "80"

	cache_path = "./cache.json"
	// SAVE_IN_JSON as said its save your database in ./cache.json
	// if SAVE_IN_JSON is enabled all your data will not be lost
	// and restored when server will be started
	//
	// if you have disable it all your data when server will stop will be gone
	save_cache = true

	// FALSE
	// whenever you will store new key in database
	// if this key exist it will not be changed
	// ---
	// TRUE
	// whenever you will store new key the old key will be repalced with the new one
	replace_key = true
)

type TrashKvMuxConfig struct {
	Port       string
	SaveCache  bool
	CachePath  string
	ReplaceKey bool
}

func (config *TrashKvMuxConfig) Serve() {
	server_name = uuid.NewString()

	port = config.Port
	save_cache = config.SaveCache
	cache_path = config.CachePath
	replace_key = config.ReplaceKey

	h2s := &http2.Server{}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	handler := http.HandlerFunc(TkvHandler)

	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(handler, h2s),
	}

	fmt.Printf("Listening %s...\n", addr)
	checkErr(server.ListenAndServe(), "while listening")
}

/*
func TkvRouteConnect(w http.ResponseWriter, r *http.Request) {
	if save_cache {
		connkey := r.URL.Query().Get("key")
		if _, err := os.Stat(cache_path); !os.IsNotExist(err) {
			res := make(map[string]interface{})
			file, err := ioutil.ReadFile(cache_path)
			if err != nil {
				log.Println(err)
			}

			if len(file) > 0 {
				cache, err := decrypt([]byte(connkey), string(file))
				if err != nil {
					log.Println(err)
				}

				if err = json.Unmarshal([]byte(cache), &res); err != nil {
					log.Println(err)
				}

				for key, value := range res {
					tkvdb.Store(key, value)
				}

				if len(cache) > 2 {
					global_private_key = []byte(connkey)
				}
			}
		}
	}
}
*/

func TkvHandler(w http.ResponseWriter, r *http.Request) {
	var response reqHTTPdataSave

	// check if request is not nil
	if r.Method == "POST" {
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			log.Println(err)
		}

		tkvdb = *response.Cache
		if save_cache {
			save_cache_file(&response)
		}
	} else if r.Method == "GET" {
		json.NewEncoder(w).Encode(tkvdb)
	}
}


func save_cache_file(response *reqHTTPdataSave) {
	j, err := json.Marshal(&response.Cache)
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(cache_path, j, 0744)
}

func checkErr(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("ERROR: %s: %s\n", msg, err)
	os.Exit(1)
}

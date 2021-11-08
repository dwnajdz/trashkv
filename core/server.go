package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/google/uuid"
	"golang.org/x/sync/syncmap"
)

// global database variable
var tkvdb syncmap.Map

// private key for server
var global_private_key []byte

// auth key is key for making database/server connection safer
// it is creating new uuid key
// everytime user is saving
//var auth_security_key = uuid.New().String()

// config
var (
	// used in 119 line in sync_with_servers() function
	// It is optional you can leave it blank
	SERVER_NAME = "node0"
	SERVER_URL  = fmt.Sprintf("http://localhost:%s", PORT)
	// declare servers and child servers names
	// !!!
	// always declare current server name first
	SERVERS_JSON = map[string]string{
		"node": fmt.Sprintf("http://localhost:%s", PORT),
		// example of second server
		//"child1": fmt.Sprintf("http://localhost:8894",),
	}
	SERVERS_JSON_PATH = "./servers.json"

	// port of server
	// you can set it to whatever port you want
	PORT = "80"

	CACHE_PATH = "./cache.tkv"
	// SAVE_IN_JSON as said its save your database in ./cache.json
	// if SAVE_IN_JSON is enabled all your data will not be lost
	// and restored when server will be started
	//
	// if you have disable it all your data when server will stop will be gone
	SAVE_CACHE = true

	// FALSE
	// whenever you will store new key in database
	// if this key exist it will not be changed
	// ---
	// TRUE
	// whenever you will store new key the old key will be repalced with the new one
	REPLACE_KEY = true
)

// port must have ':' before number
func Host(port string, server *http.ServeMux) {
	url := fmt.Sprintf("localhost%s", port)

	server.HandleFunc("/tkv_v1/connect", TkvRouteConnect)
	server.HandleFunc("/tkv_v1/save", TkvRouteCompareAndSave)
	server.HandleFunc("/tkv_v1/sync", TkvRouteSyncWithServers)
	server.HandleFunc("/tkv_v1/status", TkvRouteStatus)
	server.HandleFunc("/tkv_v1/servers.json", TkvRouteServersJson)
	http.PostForm(fmt.Sprintf("http://%s/tkv_v1/sync", url), nil)

	go func() {
		log.Printf("Listening on http://%s", url)
		http.ListenAndServe(url, server)
	}()

	defer log.Println("server stopped working")
}

func TkvRouteConnect(w http.ResponseWriter, r *http.Request) {
	if SAVE_CACHE {
		connkey := r.URL.Query().Get("key")
		if _, err := os.Stat(CACHE_PATH); !os.IsNotExist(err) {
			res := make(map[string]interface{})
			file, err := ioutil.ReadFile(CACHE_PATH)
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

	dataMap := make(map[string]interface{})
	tkvdb.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&dataMap)
	if err != nil {
		fmt.Println(err)
	}

	if global_private_key == nil {
		fmt.Fprint(w, string(j))
	} else {
		txt, err := encrypt(global_private_key, string(j))
		if err != nil {
			log.Println(err)
		}

		fmt.Fprint(w, txt)
	}
}

func TkvRouteCompareAndSave(w http.ResponseWriter, r *http.Request) {
	var response reqHTTPdataSave
	var newdb syncmap.Map

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Println(err)
	}

	// check if request is not nil
	if r.Method == "POST" {
		//if response.AuthKey != &auth_security_key
		for key, value := range *response.Cache {
			newdb.Store(key, value)
		}

		tkvdb = newdb
		if global_private_key == nil {
			global_private_key = *response.PrivateKey
		}

		if SAVE_CACHE {
			save_cache_file(&response)
		}
	}
}

func save_cache_file(response *reqHTTPdataSave) {
	j, err := json.Marshal(&response.Cache)
	if err != nil {
		log.Println(err)
	}

	txt, err := encrypt(global_private_key, string(j))
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(CACHE_PATH, []byte(txt), 0744)
}

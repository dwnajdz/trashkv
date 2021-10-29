package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/sync/syncmap"
)

// global database variable
var tkvdb syncmap.Map

// private key for server
var global_private_key []byte

// auth key is key for making database/server connection safer
// it is creating new uuid key
// everytime user is saving
var auth_security_key = uuid.New().String()

// config
var (
	// used in 119 line in sync_with_servers() function
	// It is optional you can leave it blank
	SERVER_NAME = "node0"
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

	CACHE_PATH = "./cache.json"
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
	REPLACE_KEY = false
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
		if _, err := os.Stat(CACHE_PATH); !os.IsNotExist(err) {
			res := make(map[string]interface{})
			file, err := ioutil.ReadFile(CACHE_PATH)
			if err != nil {
				log.Println(err)
			}

			if err := json.Unmarshal(file, &res); err != nil {
				log.Println(err)
			}

			for key, value := range res {
				tkvdb.Store(key, value)
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
	var response reqServerSave
	var newdb syncmap.Map

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Println(err)
	}

	// check if request is not nil
	if r.Method == "POST" {
		if response.AuthKey != &auth_security_key {
			for key, value := range *response.Cache {
				newdb.Store(key, value)
			}

			tkvdb = newdb
			if global_private_key == nil {
				global_private_key = *response.PrivateKey
			}

			// send request to make all servers synchronized
			http.PostForm(fmt.Sprintf("http://localhost:%s/tkv_v1/sync", PORT), nil)

			if SAVE_CACHE {
				j, err := json.Marshal(&response.Cache)
				if err != nil {
					log.Println(err)
				}
				ioutil.WriteFile(CACHE_PATH, j, 0644)
			}

			auth_security_key = uuid.NewString()
		}
	}
}

func TkvRouteStatus(w http.ResponseWriter, r *http.Request) {
	servers := ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)
	result := make(map[string]string)

	for key, value := range servers {
		_, err := Connect(value, "")
		if err == nil {
			result[key] = "active"
		} else {
			result[key] = "dead"
		}
	}

	jsonRes, err := json.MarshalIndent(&result, " ", " ")
	if err != nil {
		log.Println(err)
	}

	fmt.Fprint(w, string(jsonRes))
}

func TkvRouteSyncWithServers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("/sync request")
		jsonf := ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)

		for key, value := range jsonf {
			if key != SERVER_NAME {
				syncAllServers(tkvdb, value)
			}
		}
	}
}

func syncAllServers(inDatabase syncmap.Map, url string) {
	dataMap := make(map[string]interface{})
	inDatabase.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&dataMap)
	if err != nil {
		fmt.Println(err)
	}

	http.Post(fmt.Sprintf("%s/tkv_v1/sync/save", url), "application/json", bytes.NewBuffer(j))
}

func TkvRouteServersJson(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile(SERVERS_JSON_PATH)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprint(w, string(file))
}

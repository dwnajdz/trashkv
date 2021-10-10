package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/wspirrat/trashkv/core"
	"golang.org/x/sync/syncmap"
)

var (
	// global database variable
	tkvdb syncmap.Map
	// NODE_0 is name of first server
	// used in 153 line in sync_with_servers() function
	NODE_0 = "node"
	// declare servers and child servers names
	// !!!
	// always declare current server name first
	SERVERS_JSON = map[string]string{
		"node": fmt.Sprintf("http://localhost:%s", port),
		// example of second server
		//"child1": fmt.Sprintf("http://localhost:8894",),
	}
	SERVERS_JSON_PATH = "./servers.json"
)

const (
	// port of server
	// you can set it to whatever port you want
	port = "80"

	CACHE_PATH = "./cache.json"
	// SAVE_IN_JSON as said its save your database in ./cache.json
	// if SAVE_IN_JSON is enabled all your data will not be lost
	// and restored when server will be started
	//
	// if you have disable it all your data when server will stop will be gone
	SAVE_IN_JSON = true
)

func main() {
	// initialise variables
	tkvdb = syncmap.Map{}

	if SAVE_IN_JSON {
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

	http.HandleFunc("/tkv_v1/connect", connect)
	http.HandleFunc("/tkv_v1/save", compare_and_save)
	http.HandleFunc("/tkv_v1/sync", sync_with_servers)
	http.HandleFunc("/tkv_v1/status", status)
	http.HandleFunc("/tkv_v1/servers.json", servers_json)
	http.PostForm(fmt.Sprintf("http://localhost:%s/tkv_v1/sync", port), nil)

	http.ListenAndServe(":"+port, nil)
}

func connect(w http.ResponseWriter, r *http.Request) {
	dataMap := make(map[string]interface{})
	tkvdb.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&dataMap)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, string(j))
}

func compare_and_save(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	var newdb syncmap.Map

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
	}

	// check if request is not nil
	if r.Method == "POST" {
		for key, value := range request {
			newdb.Store(key, value)
		}

		tkvdb = newdb

		// send request to all servers be synced
		http.PostForm(fmt.Sprintf("http://localhost:%s/tkv_v1/sync", port), nil)

		if SAVE_IN_JSON {
			j, err := json.Marshal(&request)
			if err != nil {
				log.Println(err)
			}
			ioutil.WriteFile(CACHE_PATH, j, 0644)
		}
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	servers := core.ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)
	result := make(map[string]string)

	for key, value := range servers {
		_, err := core.Connect(value)
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

	fmt.Fprintf(w, string(jsonRes))
}

func sync_with_servers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("/sync request")
		jsonf := core.ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)

		for key, value := range jsonf {
			if key != NODE_0 {
				save(tkvdb, value)
			}
		}
	}
}

func save(inDatabase syncmap.Map, url string) {
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

func servers_json(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile(SERVERS_JSON_PATH)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprint(w, string(file))
}

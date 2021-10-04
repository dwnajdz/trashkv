package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wspirrat/trashkv/core"
	"golang.org/x/sync/syncmap"
)

var (
	db                  syncmap.Map
	dbJson              []byte
	serversLatestStatus map[string]bool
)

// port of server
const (
	port        = "4998"
	server_name = "server1"
)

func main() {
	db = syncmap.Map{}

	log.Printf("sever running on http://localhost:%s", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "working")
	})
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/save", compare_and_save)
	http.HandleFunc("/sync", sync_with_servers)

	http.ListenAndServe(":"+port, nil)
}

func connect(w http.ResponseWriter, r *http.Request) {
	dataMap := make(map[string]interface{})
	db.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&dataMap)
	if err != nil {
		fmt.Println(err)
	}

	dbJson = j

	fmt.Fprint(w, string(j))
}

func compare_and_save(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	var newdb syncmap.Map

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
	}

	// check if request is not nil
	if r.Method == "POST" {
		for key, value := range request {
			newdb.Store(key, value)
		}

		db = newdb
	}
}

func sync_with_servers(w http.ResponseWriter, r *http.Request) {
	jsonf := readSeversJson()

	for key, value := range jsonf {
		if key != server_name {
			connectToDb := core.Connect(value)
			server_db := connectToDb.Access()
			
			serverDataMap := make(map[string]interface{})
			server_db.Range(func(k interface{}, v interface{}) bool {
				serverDataMap[k.(string)] = v
				return true
			})

			serverDbJson, err := json.Marshal(&serverDataMap)
			if err != nil {
				log.Println(err)
			}

			if !bytes.Equal(serverDbJson, dbJson) {
				save(db, value)
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

	http.Post(fmt.Sprintf("%s/save", url), "application/json", bytes.NewBuffer(j))
}

func readSeversJson() map[string]string {
	var res map[string]string

	file, err := ioutil.ReadFile("./servers.json")
	if err != nil {
		log.Println(err)
	}

	if err := json.Unmarshal(file, &res); err != nil {
		log.Println(err)
	}

	return res
}

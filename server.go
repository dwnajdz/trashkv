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
	db     syncmap.Map
	dbJson []byte
	// heroku
	port = os.Getenv("PORT")
)

// port of server
const (
	//port        = "80"
	server_name = "main"
)

func main() {
	// initialise variables
	db = syncmap.Map{}

	if _, err := os.Stat("./db.json"); !os.IsNotExist(err) {
		res := make(map[string]interface{})
		file, err := ioutil.ReadFile("./db.json")
		if err != nil {
			log.Println(err)
		}

		if err := json.Unmarshal(file, &res); err != nil {
			log.Println(err)
		}

		for key, value := range res {
			db.Store(key, value)
		}
	}

	log.Printf("server running on http://localhost:%s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "working")
	})
	http.HandleFunc("/connect", connect)
	http.HandleFunc("/save", compare_and_save)
	http.HandleFunc("/sync", sync_with_servers)
	http.HandleFunc("/status", status)
	http.Get(fmt.Sprintf("http://localhost:%s/sync", port))

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
		log.Println(err)
	}

	// check if request is not nil
	if r.Method == "POST" {
		for key, value := range request {
			newdb.Store(key, value)
		}

		db = newdb

		j, err := json.Marshal(&request)
		if err != nil {
			log.Println(err)
		}
		ioutil.WriteFile("db.json", j, 0644)
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	servers := readSeversJson()

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
	log.Println("/sync")
	jsonf := readSeversJson()

	for key, value := range jsonf {
		if key != server_name {
			connectToDb, err := core.Connect(value)
			if err != nil {
				break
			}
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

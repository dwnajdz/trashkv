package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/sync/syncmap"
)

var (
	db syncmap.Map
	//heroku
	//port = os.Getenv("PORT")
)

// port of server
const (
	port        = "4990"
	server_name = "server1"
	// pass main server url here
	main_url = "http://localhost:80"

	SAVE_IN_JSON = false
)

func main() {
	// initialise variables
	db = syncmap.Map{}

	if SAVE_IN_JSON {
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
	}

	log.Printf("server running on http://localhost:%s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "working")
	})
	http.HandleFunc("/tkv_v1/connect", connect)
	http.HandleFunc("/tkv_v1/save", compare_and_save)
	http.HandleFunc("/tkv_v1/sync/save", sync_save)
	http.HandleFunc("/tkv_v1/sync", sync_with_servers)
	http.PostForm(fmt.Sprintf("http://localhost:%s/tkv_v1/sync", port), nil)

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

	fmt.Fprint(w, string(j))
}

func compare_and_save(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	var newdb syncmap.Map

	log.Println("/save")

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

		// send request to all servers be synced
		http.PostForm(fmt.Sprintf("http://localhost:%s/sync", port), nil)

		if SAVE_IN_JSON {
			j, err := json.Marshal(&request)
			if err != nil {
				log.Println(err)
			}
			ioutil.WriteFile("db.json", j, 0644)
		}
	}
}

// sync save is route for handling sync request
// from other servers
// why?
// without this route main server will send request to /save path
// and create infinite loop, because servers will send each other post requests
//
// maybe there is some other solution without need to create new route 
// but that's good for now
func sync_save(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	var newdb syncmap.Map

	log.Println("/save")

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

		if SAVE_IN_JSON {
			j, err := json.Marshal(&request)
			if err != nil {
				log.Println(err)
			}
			ioutil.WriteFile("db.json", j, 0644)
		}
	}
}

func sync_with_servers(w http.ResponseWriter, r *http.Request) {
	var jsonf map[string]string

	log.Println("/sync")
	resp, err := http.Get(fmt.Sprintf("%s/servers.json", main_url))

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &jsonf); err != nil {
		log.Println(err)
	}

	if r.Method == "POST" {
		for key, value := range jsonf {
			if key != server_name {
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

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// global database variable
var tkvdb map[string]interface{}

// private key for server
var global_private_key []byte

// config
var (
	CACHE_PATH  = "./cache.tkv"
	REPLACE_KEY = true
)

func Host(port, DatabaseFilePath string, server *http.ServeMux) {
	url := fmt.Sprintf("http://localhost:%s", port)
	CACHE_PATH = fmt.Sprintf("%s/cache.tkv", DatabaseFilePath)

	server.HandleFunc("/tkv_v1/connect", TkvRouteConnect)
	server.HandleFunc("/tkv_v1/save", TkvRouteCompareAndSave)

	log.Printf("Listening on http://%s", url)
	http.ListenAndServe(url, server)
}

func TkvRouteConnect(w http.ResponseWriter, r *http.Request) {
	connkey := r.URL.Query().Get("key")
	tkvdb = read_database_file(CACHE_PATH, []byte(connkey))

	j, err := json.Marshal(tkvdb)
	if err != nil {
		fmt.Println(err)
		return
	}

	txt, err := encrypt(global_private_key, string(j))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, txt)
}

func save_database_file(filepath string, response requestServerSave) {
	j, err := json.Marshal(&response.Cache)
	if err != nil {
		log.Println(err)
		return
	}

	txt, err := encrypt(global_private_key, string(j))
	if err != nil {
		log.Println(err)
		return
	}

	ioutil.WriteFile(filepath, []byte(txt), 0744)
}

func read_database_file(filepath string, decryption_key []byte) map[string]interface{} {
	res := make(map[string]interface{})
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		log.Println("File does not exist")
		return res
	}

	file, _ := ioutil.ReadFile(CACHE_PATH)
	if len(file) == 0 {
		return res
	}

	cache, err := decrypt(decryption_key, string(file))
	if err != nil {
		log.Println(err)
		return res
	}

	json.Unmarshal([]byte(cache), &res)
	return res
}

func TkvRouteCompareAndSave(w http.ResponseWriter, r *http.Request) {
	response := requestServerSave{}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		log.Println(err)
	}

	if r.Method != "POST" {
		return
	}

	tkvdb = response.Cache
	if global_private_key == nil {
		global_private_key = response.PrivateKey
	}

	save_database_file(CACHE_PATH, response)
}

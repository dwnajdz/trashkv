package core

import (
	"bytes"
	"crypto/sha256"
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
//var Http_security_key uuid.UUID

var saves = 0

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
	port string

	cache_path = "./cache.tkv"
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

func (config *TrashKvMuxConfig) Serve() http.Handler {
	server_name = uuid.NewString()

	port = config.Port
	save_cache = config.SaveCache
	cache_path = config.CachePath
	replace_key = config.ReplaceKey

	mux := http.NewServeMux()
	mux.HandleFunc("/tkv_v1/connect", TkvRouteConnect)
	mux.HandleFunc("/tkv_v1/save", TkvRouteCompareAndSave)
	mux.HandleFunc("/tkv_v1/sync", TkvRouteSyncWithServers)

	return mux
}

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
		if global_private_key == nil {
			global_private_key = *response.PrivateKey
		}

		// this is created for checking either private key is encrypted in sha256
		var check_global_private_key []byte
		// when its first insertion set global key to user key
		if saves == 0 {
			check_global_private_key = global_private_key
			saves++
		} else {
			check_global_private_key = makeSHA256(global_private_key)
		}

		if bytes.Equal(*response.PrivateKey, check_global_private_key) {
			for key, value := range *response.Cache {
				newdb.Store(key, value)
			}
			tkvdb = newdb
			if save_cache {
				save_cache_file(&response)
			}
		} else {
			http.Error(w, "aes: wrong key", http.StatusBadRequest)
		}

		//saves++
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

	ioutil.WriteFile(cache_path, []byte(txt), 0744)
}

func makeSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

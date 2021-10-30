package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/sync/syncmap"
)

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

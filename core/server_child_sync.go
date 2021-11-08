package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/syncmap"
)

func TkvRouteStatus(w http.ResponseWriter, r *http.Request) {
	servers := ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)
	result := make(map[string]string)
	keys := r.URL.Query()

	for key, value := range servers {
		_, err := Connect(value, keys.Get("sk"))
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
				log.Printf("synced with: %s", value)
				syncAllServers(tkvdb, value)
			}
		}
	}
}

func syncAllServers(inDatabase syncmap.Map, receiver string) {
	dataMap := make(map[string]interface{})
	inDatabase.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	request := reqHTTPdataSave{
		Sender:     &SERVER_URL,
		Receiver:   &receiver,
		Cache:      &dataMap,
		PrivateKey: &global_private_key,
	}
	readyToSendRequest, err := json.Marshal(&request)
	if err != nil {
		fmt.Println(err)
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 1 * time.Second,
	}
	client = &http.Client{Transport: tr}
	client.Post(fmt.Sprintf("%s/tkv_v1/save", receiver), "application/json", bytes.NewBuffer(readyToSendRequest))
}

func TkvRouteServersJson(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile(SERVERS_JSON_PATH)
	if err != nil {
		log.Println(err)
	}

	fmt.Fprint(w, string(file))
}

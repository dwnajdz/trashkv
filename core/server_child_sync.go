package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"os"

	"golang.org/x/sync/syncmap"
)

func TkvRouteStatus(secret_key string) map[string]bool {
	servers := ReadSeversJson(server_json_path, servers_json)
	// map[serverName]isActive
	result := make(map[string]bool)

	for key, value := range servers {
		_, err := Connect(value, secret_key)
		if err == nil {
			result[key] = true
		} else {
			result[key] = false
		}
	}

	return result
}

func TkvRouteSyncWithServers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("/sync request")
		jsonf := ReadSeversJson(server_json_path, servers_json)

		for key, value := range jsonf {
			if key != server_name {
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
		Sender:     &server_url,
		Receiver:   &receiver,
		Cache:      &dataMap,
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

// json function
func ReadSeversJson(path string, servers map[string]string) map[string]string {
	var res map[string]string

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}

		if err := json.Unmarshal(file, &res); err != nil {
			fmt.Println(err)
		}

		return res
	} else {
		file, _ := json.MarshalIndent(servers, "", " ")
		_ = ioutil.WriteFile(server_json_path, file, 0644)
	}
	return nil
}
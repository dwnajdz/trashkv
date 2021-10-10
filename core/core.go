package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/sync/syncmap"
)

const (
	DatabasePath = "core.db"
)

type Database struct {
	PrivateKey *string
	Url        string
	Syncmap    syncmap.Map
}

type Core interface {
	Store(key string, value interface{})
	Delete(key string)
	Load(key string) interface{}
	Save()
}

// funcs
func Connect(url string) (Core, error) {
	// dat is used for unmarshalling database from /connect
	// syncm is Syncmap passed in *Database
	// core is interface which is returned
	var dat map[string]interface{}
	var syncm syncmap.Map
	var core Core

	resp, err := http.Get(fmt.Sprintf("%s/tkv_v1/connect", url))
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &dat); err != nil {
		return nil, err
	}

	// add all keys from dat to syncm
	for key, value := range dat {
		syncm.Store(key, value)
	}

	resDb := &Database{
		PrivateKey: nil,
		Url:        url,
		Syncmap:    syncm,
	}

	core = resDb
	return core, nil
}

func (db *Database) Store(key string, value interface{}) {
	db.Syncmap.Store(key, value)
}

func (db *Database) Delete(key string) {
	db.Syncmap.Delete(key)
}

func (db *Database) Load(key string) interface{} {
	result, ok := db.Syncmap.Load(key)
	if ok {
		return result
	}

	return "this value does not exist in map"
}

// save function send request to server
// server compare and set var db *Database
// as database send in json request
func (db *Database) Save() {
	dataMap := make(map[string]interface{})
	db.Syncmap.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&dataMap)
	if err != nil {
		fmt.Println(err)
	}

	http.Post(fmt.Sprintf("%s/tkv_v1/save", db.Url), "application/json", bytes.NewBuffer(j))
}

func (db *Database) Access() syncmap.Map {
	return db.Syncmap
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
		_ = ioutil.WriteFile("test.json", file, 0644)
	}
	return nil
}

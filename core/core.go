package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/sync/syncmap"
)

const (
	DatabasePath = "core.db"
)

type Database struct {
	PrivateKey *string
	Syncmap    syncmap.Map
}

type Core interface {
	Store(key string, value interface{})
	Delete(key string)
	Load(key string) interface{}
	Save()
}

// funcs
func Connect() Core {
	// dat is used for unmarshalling database from /connect
	// syncm is Syncmap passed in *Database
	// core is interface which is returned
	var dat map[string]interface{}
	var syncm syncmap.Map
	var core Core

	resp, err := http.Get("http://localhost/connect")
	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &dat); err != nil {
		fmt.Println(err)
	}

	// add all keys from dat to syncm
	for key, value := range dat {
		syncm.Store(key, value)
	}

	resDb := &Database{
		PrivateKey: nil,
		Syncmap:    syncm,
	}

	core = resDb
	return core
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

	_, err = http.Post("http://localhost/save", "application/json", bytes.NewBuffer(j))
	if err != nil {
		fmt.Println(err)
	}
}

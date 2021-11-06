package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/syncmap"
)

const (
	DatabasePath = "core.db"
)

type Database struct {
	PrivateKey *[]byte
	Url        string
	Syncmap    syncmap.Map
}

type Core interface {
	Store(key string, value interface{})
	Delete(key string)
	Load(key string) (interface{}, bool)
	Save()
	Access() syncmap.Map // only for dashboard
}

// req = request
// request server save
type reqServerSave struct {
	AuthKey    *string                 `json:"AuthKey"`
	Cache      *map[string]interface{} `json:"Cache"`
	PrivateKey *[]byte                 `json:"PrivateKey"`
}

var client *http.Client

// funcs
func Connect(url, privateKey string) (Core, error) {
	// dat is used for unmarshalling database from /connect
	// syncm is Syncmap passed in *Database
	// core is interface which is returned
	var dat map[string]interface{}
	var syncm syncmap.Map
	var core Core

	resp, err := http.Get(fmt.Sprintf("%s/tkv_v1/connect?key=%s", url, privateKey))
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if len(body) <= 2 {
		if err := json.Unmarshal([]byte(body), &dat); err != nil {
			return nil, err
		}
	} else {
		txt, err := decrypt([]byte(privateKey), string(body))
		if err != nil {
			return nil, errors.New("private key is wrong")
		}
		// unmarshal decrypted text
		if err := json.Unmarshal([]byte(txt), &dat); err != nil {
			return nil, errors.New("private key is wrong")
		}
	}

	// add all keys from dat to syncm
	for key, value := range dat {
		syncm.Store(key, value)
	}

	dbpk := []byte(privateKey)
	resDb := &Database{
		PrivateKey: &dbpk,
		Url:        url,
		Syncmap:    syncm,
	}

	core = resDb

	return core, nil
}

func (db *Database) Store(key string, value interface{}) {
	if !REPLACE_KEY {
		_, exist := db.Syncmap.Load(key)
		if !exist {
			db.Syncmap.Store(key, value)
		}
	} else {
		db.Syncmap.Store(key, value)
		return
	}
}

func (db *Database) Delete(key string) {
	db.Syncmap.Delete(key)
}

// returns value
// and bool
// if object exist returns true
// else if object do not exist returns false
func (db *Database) Load(key string) (value interface{}, exist bool) {
	result, exist := db.Syncmap.Load(key)
	if exist {
		return result, true
	}

	return nil, false
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

	request := &reqServerSave{
		AuthKey:    &auth_security_key,
		Cache:      &dataMap,
		PrivateKey: db.PrivateKey,
	}

	j, err := json.Marshal(&request)
	if err != nil {
		fmt.Println(err)
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 1 * time.Second,
	}
	client = &http.Client{Transport: tr}

	client.Post(fmt.Sprintf("%s/tkv_v1/save", db.Url), "application/json", bytes.NewBuffer(j))
}

func (db *Database) Access() syncmap.Map {
	return db.Syncmap
}

// json function
// if you want only read json in second argument just pass nil
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
		_ = ioutil.WriteFile(SERVERS_JSON_PATH, file, 0644)
	}
	return nil
}

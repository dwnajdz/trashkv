package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Database struct {
	PrivateKey []byte
	Url        string
	Storage    map[string]interface{}
}

type Core interface {
	Store(key string, value interface{})
	Delete(key string)
	Load(key string) interface{}
	Save()
}

type requestServerSave struct {
	Cache      map[string]interface{} `json:"Cache"`
	PrivateKey []byte                 `json:"PrivateKey"`
}

func Connect(url string, privateKey []byte) (Core, error) {
	// connect to the server with key
	resp, err := http.Get(fmt.Sprintf("%s/tkv_v1/connect?key=%s", url, privateKey))
	if err != nil {
		return nil, err
	}

	// create data for local database
	data := make(map[string]interface{})

	body, _ := ioutil.ReadAll(resp.Body)
	if len(body) <= 2 {
		if err := json.Unmarshal([]byte(body), &data); err != nil {
			return nil, err
		}
	} else {
		txt, err := decrypt([]byte(privateKey), string(body))
		if err != nil {
			return nil, errors.New("private key is wrong")
		}
		// unmarshal decrypted text
		if err := json.Unmarshal([]byte(txt), &data); err != nil {
			return nil, errors.New("private key is wrong")
		}
	}

	core := &Database{
		PrivateKey: privateKey,
		Url:        url,
		Storage:    data,
	}
	return core, nil
}

func (db *Database) Store(key string, value interface{}) {
	_, exist := db.Storage[key]
	if exist && !REPLACE_KEY {
		return
	}
	db.Storage[key] = value
}

func (db *Database) Delete(key string) {
	db.Storage = nil
}

func (db *Database) Load(key string) interface{} {
	result, exist := db.Storage[key]
	if exist {
		return result
	}
	return exist
}

// save function send request to server
// server compare and set var db *Database
// as database send in json request
func (db *Database) Save() {
	request := requestServerSave{
		Cache:      db.Storage,
		PrivateKey: db.PrivateKey,
	}

	j, err := json.Marshal(request)
	if err != nil {
		return
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 1 * time.Second,
	}

	client := &http.Client{Transport: tr}
	client.Post(fmt.Sprintf("%s/tkv_v1/save", db.Url), "application/json", bytes.NewBuffer(j))
}

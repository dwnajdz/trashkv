
# trashkv

<p> trashkv is simple key-value store </p>

# Table of Contents
1. [Installing](#installing)
2. [Usage](#usage )
3. [License](#License)

## Installing

* Own files option
  1. Firstly install core of trashkv
    ``` go get github.com/wspirrat/traskhv_core ```

  2. Create **tkv_server.go** file in your project
  3. Paste this into **tkv_server.go** file
  ```go
  package main

  import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "github.com/wspirrat/trashkv/core"
    "golang.org/x/sync/syncmap"
  )

  var (
    // global database variable
    tkvdb syncmap.Map
    // NODE_0 is name of first server
    // used in 153 line in sync_with_servers() function
    NODE_0 = "node0"
    // declare servers and child servers names
    // !!!
    // always declare current server name first
    SERVERS_JSON = map[string]string{
      "node": fmt.Sprintf("http://localhost:%s", port),
      // example of second server
      //"child1": fmt.Sprintf("http://localhost:8894",),
    }
    SERVERS_JSON_PATH = "./servers.json"
  )

  const (
    CACHE_PATH = "./cache.json"
    // SAVE_IN_JSON as said its save your database in ./cache.json
    // if SAVE_IN_JSON is enabled all your data will not be lost
    // and restored when server will be started
    //
    // if you have disable it all your data when server will stop will be gone
    SAVE_IN_JSON = true
  )

  func connect(w http.ResponseWriter, r *http.Request) {
    dataMap := make(map[string]interface{})
    tkvdb.Range(func(k interface{}, v interface{}) bool {
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

      tkvdb = newdb

      // send request to all servers be synced
      http.PostForm(fmt.Sprintf("http://localhost:%s/tkv_v1/sync", port), nil)

      if SAVE_IN_JSON {
        j, err := json.Marshal(&request)
        if err != nil {
          log.Println(err)
        }
        ioutil.WriteFile(CACHE_PATH, j, 0644)
      }
    }
  }

  func status(w http.ResponseWriter, r *http.Request) {
    servers := core.ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)
    result := make(map[string]string)

    for key, value := range servers {
      _, err := core.Connect(value)
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

    fmt.Fprintf(w, string(jsonRes))
  }

  func sync_with_servers(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
      log.Println("/sync request")
      jsonf := core.ReadSeversJson(SERVERS_JSON_PATH, SERVERS_JSON)

      for key, value := range jsonf {
        if key != NODE_0 {
          save(tkvdb, value)
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

    http.Post(fmt.Sprintf("%s/tkv_v1/sync/save", url), "application/json", bytes.NewBuffer(j))
  }

  func servers_json(w http.ResponseWriter, r *http.Request) {
    file, err := ioutil.ReadFile(SERVERS_JSON_PATH)
    if err != nil {
      log.Println(err)
    }

    fmt.Fprint(w, string(file))
  }
  ```
## Usage

1. **Connect to database server**
  <p> To connect you just pass url for server into function and assign database variable to it.</p>

  ```go 
  import 	(
    core "github.com/wspirrat/trashkv_core"
  )

  db := core.Connect("http://localhost:80")
  ```


## Third Example
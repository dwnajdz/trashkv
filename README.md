
# trashkv

<p> trashkv is simple key-value store </p>

# Table of Contents
1. [Installing](#installing)
2. [Usage](#usage )
3. [License](#License)

## Installing

1. Firstly install trashkv.
  ``` go get github.com/wspirrat/traskhv ```

2. Next paste this into your file with all http routes.
    ```go
      http.HandleFunc("/tkv_v1/connect", core.TkvRouteConnect)
      http.HandleFunc("/tkv_v1/save", core.TkvRouteCompareAndSave)
      http.HandleFunc("/tkv_v1/sync", core.TkvRouteSyncWithServers)
      http.HandleFunc("/tkv_v1/status", core.TkvRouteStatus)
      http.HandleFunc("/tkv_v1/servers.json", core.TkvRouteServersJson)
      http.PostForm("http://localhost:80/tkv_v1/sync"), nil)
      // alternative with custom port 
      // import "fmt"
      // http.PostForm(fmt.Sprintf("http://localhost:%s/tkv_v1/sync", yourport), nil)
    ```

3. Configure. 
    All configue options avaliable in [server.go](https://github.com/wspirrat/trashkv/blob/master/core/server.go)
    ```go
    // config
    var (
      // used in 119 line in sync_with_servers() function
      // It is optional you can leave it blank
      SERVER_NAME = "node0"
      // declare servers and child servers names
      // !!!
      // always declare current server name first
      SERVERS_JSON = map[string]string{
        "node": fmt.Sprintf("http://localhost:%s", PORT),
        // example of second server
        //"child1": fmt.Sprintf("http://localhost:8894",),
      }
      SERVERS_JSON_PATH = "./servers.json"
      
      // port of server
      // you can set it to whatever port you are using
      PORT = "80"

      CACHE_PATH = "./cache.json"
      // SAVE_IN_JSON as said it is saving your database in ./cache.json
      // if SAVE_IN_JSON is enabled all your data will not be lost
      // and restored when server will be started
      //
      // if you have disable it all your data when server will stop will be gone
      SAVE_CACHE = true
    )
    ```
  4. Run server
    ``` go run main.go ```

## Usage

1. **Connect to database server**
    <p> To connect you just pass url for server into function and assign database variable to it.</p>

    ***Warning***: Pass url without ```/tkv_v1/connect```. Only link to your application
    
    ```go 
    import 	(
      "fmt"
    
      "github.com/wspirrat/trashkv/core"
    )

    db, err := core.Connect("http://localhost:80")
    if err != nil {
      fmt.Println(err)
    }

    // custom url connect
    db = core.Connect("https://urltomypagewithtrashkv.com")
    ```

2. **Store key in database**
    <p> Storing is very easy. You just pass in first argue your key name. In second argue just pass value of that key. </p>

    <p> TrashKv accepts all types of variables in golang. </p>

    ***Important***: *When key already exist in database it is replaced with the new* 

    ***Warning***: *All keys must be strings*

    ```go
    import 	(
        "github.com/wspirrat/trashkv/core"
    )

    func main() {
      // connecting to db
      db, _ := core.Connect("http://localhost:80")

      // storing string
      db.Store("mystring", "hello")

      // storing int
      db.Store("int", 1)

      // storing array
      db.Store("array", []string{"hello im array"})

      // storing byte
      db.Store("byte", []byte("hello im array")

      // storing struct
      type person struct {
        Name string
        Age int
        Childrens []string
        IsMarried bool
        Bank float64
      }
      db.Store("John", person{
        "John",
        102,
        []string{"John2"},
        true,
        123000.345,
      })

      db.Save()
    }
    ```

3. **Loading keys / getting keys**
  <p> To load key you just pass in first argue it's name.</p>
  
  ```go
  import 	(
    "fmt"

    "github.com/wspirrat/trashkv/core"
  )

  type person struct {
    Name string
    Age int
    Childrens []string
    IsMarried bool
    Bank float64
  }

  func main() {
    // connecting to db
    db, _ := core.Connect("http://localhost:80")

    db.Store("mystring", "hello")
    mystring := db.Load("mystring")
    fmt.Println(mystring)
    // result: hello

    db.Store("John", person{
      "John",
      102,
      []string{"John2"},
      true,
      123000.345,
    })
    // ACCESSING STRUCT
  	john := db.Load("John").(person)
    fmt.Println(john.Bank)
    // result: 123000.345

    db.Save()
  }
  ```

4. **Deleting keys**
  <p> To delete key just pass it's name. </p>

  ```go
  import "github.com/wspirrat/trashkv/core"

  func main() {
    // connecting to db
    db := core.Connect("http://localhost:80")

    db.Store("mystring", "hello")
    db.Delete("mystring")
    // done

    db.Save()
  }
  ```

5. **Save**
  <p> Always remember to save your database after *store/delete* functions. </p>

  ```go
  // To do that you just need to pass (*Database).Save() function

  import "github.com/wspirrat/trashkv/core"

  func main() {
    // connecting to db
    db := core.Connect("http://localhost:80")
    db.Store("mystring", "hello")

    db.Save()
  }
  ```
## License
[MIT](https://choosealicense.com/licenses/mit/)

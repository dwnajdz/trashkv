package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	fmt.Println("Enter url of database you want to sync with: ")
	var sync_url string
	fmt.Scanln(&sync_url)

	fmt.Println("Enter this database port: ")
	var port string
	fmt.Scanln(&port)

	core.SAVE_CACHE = false
	core.REPLACE_KEY = true
	http.HandleFunc("/tkv_v1/connect", core.TkvRouteConnect)
	http.HandleFunc("/tkv_v1/save", core.TkvRouteCompareAndSave)
	http.HandleFunc("/tkv_v1/sync", core.TkvRouteSyncWithServers)
	http.HandleFunc("/tkv_v1/status", core.TkvRouteStatus)
	http.HandleFunc("/tkv_v1/servers.json", core.TkvRouteServersJson)

	log.Println("server working...")
	http.ListenAndServe(":" + port, nil)
}

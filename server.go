package main

import (
	"fmt"
	"net/http"

	"github.com/wspirrat/trashkv/core"
)

func main() {	
	core.SAVE_CACHE = false
	http.HandleFunc("/tkv_v1/connect", core.TkvRouteConnect)
	http.HandleFunc("/tkv_v1/save", core.TkvRouteCompareAndSave)
	http.HandleFunc("/tkv_v1/sync", core.TkvRouteSyncWithServers)
	http.HandleFunc("/tkv_v1/status", core.TkvRouteStatus)
	http.HandleFunc("/tkv_v1/servers.json", core.TkvRouteServersJson)
	http.PostForm("http://localhost:80/tkv_v1/sync", nil)

	db, _ := core.Connect("http://localhost:80")
	fmt.Println(db)

	http.ListenAndServe(":80", nil)
} 
package main

import (
	"log"
	"net/http"
	//"os"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	//port := os.Getenv("PORT")

	trashkv := core.TrashKvMuxConfig{
		Port:       "80",
		SaveCache:  false,
		CachePath:  "./cache.tkv",
		ReplaceKey: true,
	}

	handler := trashkv.Serve()

	log.Println("Server started on: http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", handler))
}

package main

import (
	"log"
	"net/http"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	trashkv := core.TrashKvMuxConfig{
		Port:       "80",
		SaveCache:  true,
		CachePath:  "./cache.tkv",
		ReplaceKey: true,
	}

	handler := trashkv.Serve()

	log.Println("Server started on: http://localhost:80")
	http.ListenAndServe("localhost:80", handler)
}

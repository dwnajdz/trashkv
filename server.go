package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	port := os.Getenv("PORT")

	trashkv := core.TrashKvMuxConfig{
		Port:       port,
		SaveCache:  true,
		CachePath:  "./cache.tkv",
		ReplaceKey: true,
	}

	handler := trashkv.Serve()

	log.Println("Server started on: http://localhost:80")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

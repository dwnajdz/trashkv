/*
req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
  -keyout "C:\Users\Ja\Desktop\kv_database\rootca.key" -out "C:\Users\Ja\Desktop\kv_database\rootca.crt" -subj "/CN=trashkv.com" \
  -addext "subjectAltName=DNS:example.com,DNS:www.example.net,IP:10.0.0.1"
*/

// For generating SSL certs
//genrsa -out "C:\Users\Ja\Desktop\kv_database\rootca.key" 2048
//req -new -out "C:\Users\Ja\Desktop\kv_database\rootca.crt" -key "C:\Users\Ja\Desktop\kv_database\rootca.key" -subj "/C=US/CN=TrashKeyValueStore-CA"
//req -x509 -newkey rsa:4096 -keyout "C:\Users\Ja\Desktop\kv_database\rootca.key" -out "C:\Users\Ja\Desktop\kv_database\rootca.crt" -sha256 -days 365

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
		SaveCache:  true,
		CachePath:  "./cache.tkv",
		ReplaceKey: true,
	}

	handler := trashkv.Serve()

	//	log.Println("Server started on: http://localhost:80")
	// log.Fatal(http.ListenAndServe(":"+port, handler))
	err := http.ListenAndServeTLS(":80", "./rootca.crt", "./rootca.key", handler)
	log.Fatal(err)
}


package main

import (
	//"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	start := time.Now()

	db, _ := core.Connect("http://localhost:80")
	for i := 0; i < 500000; i+=10 {
		db.Store("k"+strconv.Itoa(i), i)
	}

	db.Save()
	elapsed := time.Since(start)

	b, err := ioutil.ReadFile("output.txt")
	if err != nil {
		panic(err)
	}

	b_ctx := string(b) + elapsed.String() + ","

	// write the whole body at once
	err = ioutil.WriteFile("output.txt", []byte(b_ctx), 0644)
	if err != nil {
		panic(err)
	}

	//fmt.Println(elapsed)
}

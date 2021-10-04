package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	start := time.Now()

	db := core.Connect("http://localhost:80")
	for i:=0; i<100000; i++ {
		db.Store("k"+strconv.Itoa(i), i)
	}

 	db.Save()

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

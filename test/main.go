package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wspirrat/trashkv/core"
)

type MyStruct struct {
	Id     int
	Name   string
	Idname string
}

func main() {
	start := time.Now()
	db, _ := core.Connect("http://localhost:1010", "hello")
	
	for i := 0; i < 1000000; i++ {
		db.Store("k"+strconv.Itoa(i), i)
	}

	db.Save()
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

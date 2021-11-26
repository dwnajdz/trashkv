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
	db, _ := core.Connect("http://localhost:80", "mykey:)" )
	
	for i := 0; i < 1; i++ {
		save := MyStruct{
			Id:     i,
			Name:   "k" + strconv.Itoa(i),
			Idname: strconv.Itoa(i),
		}

		db.Store("k"+strconv.Itoa(i), save)
	}

	db.Save()

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

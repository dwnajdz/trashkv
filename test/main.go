package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wspirrat/trashkv/core"
)

type db_save struct {
	id     int
	name   string
	idname string
}

func main() {
	start := time.Now()
	core.REPLACE_KEY = true

	db, err := core.Connect("http://localhost:80", "hello")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 501; i++ {
		save := db_save{
			id:     i,
			name:   "k" + strconv.Itoa(i),
			idname: strconv.Itoa(i),
		}
		db.Store("k"+strconv.Itoa(i), save)
	}

	db.Save()
	db.Store("k500", "changed :(")
	answer, exist := db.Load("k500")
	fmt.Println(answer, ",", exist)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

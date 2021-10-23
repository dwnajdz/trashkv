package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	start := time.Now()

	db, err := core.Connect("http://localhost:80", "hello")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 500000; i++ {
		db.Store("k"+strconv.Itoa(i), i)
	}
	
	db.Save()
	answer, _ := db.Load("k10")
	fmt.Println(answer)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

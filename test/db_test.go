package main

import (
	"strconv"
	"testing"

	"github.com/wspirrat/trashkv/core"
)

//ok  	github.com/wspirrat/trashkv/test	2.047s
func BenchmarkStoreSave(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db, _ := core.Connect("http://localhost:1010", "hello")
		db.Store("k"+strconv.Itoa(i), i)
		db.Save()
	}
}

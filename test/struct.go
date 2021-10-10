package main

import (
	"fmt"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	db, err := core.Connect("http://localhost:80")
	if err != nil {
		fmt.Println(err)
	}

	type person struct {
    Name string
    Age int
    Childrens []string
    IsMarried bool
    Bank float64
  }
  db.Store("John", person{
    "John",
    102,
    []string{"John2"},
    true,
    123000.345,
  })

	john := db.Load("John").(person)

  fmt.Println(john.Bank)
  db.Save()

}
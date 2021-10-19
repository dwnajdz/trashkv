package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/wspirrat/trashkv/core"
)

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)

	start := time.Now()

	db, err := core.Connect("http://localhost:80", "mykey")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 500000; i++ {
		db.Store("k"+strconv.Itoa(i), i)
	}

	db.Save(ctx)

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

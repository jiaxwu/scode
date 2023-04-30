package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jiaxwu/scode"
	"github.com/jiaxwu/scode/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	// store method for short code
	store := store.NewRedisStore(&redis.Options{
		Addr: "localhost:6379",
	}, time.Hour*24)
	// scode client
	scode := scode.New(store, "0123456789", scode.WithTryTimes(3))
	// allocate a short code
	code, err := scode.Allocate(context.Background(), 9)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("short code %s\n", code)
	// release short code
	if err := scode.Release(context.Background(), code); err != nil {
		log.Fatal(err)
	}
}

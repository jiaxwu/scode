package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func getRedisStore() Store {
	return NewRedisStore(&redis.Options{
		Addr: "localhost:6379",
	}, time.Second*10)
}

func TestRedisStore_SetIfNotExists(t *testing.T) {
	s := getRedisStore()
	set, err := s.SetIfNotExists(context.Background(), "abc")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(set)
}

func TestRedisStore_Delete(t *testing.T) {
	s := getRedisStore()
	if err := s.Delete(context.Background(), "abc"); err != nil {
		t.Error(err)
	}
}

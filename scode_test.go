package scode

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jiaxwu/scode/store"
	"github.com/redis/go-redis/v9"
)

func getRedisStore() store.Store {
	return store.NewRedisStore(&redis.Options{
		Addr: "localhost:6379",
	}, time.Hour*24)
}

func getScode() *Scode {
	return New(getRedisStore(), "0123456789", WithTryTimes(3))
}

func TestScode_Allocate(t *testing.T) {
	s := getScode()
	shortCode, err := s.Allocate(context.Background(), 9)
	if err != nil {
		t.Error(err)
	}
	t.Log(shortCode)
}

func TestScode_AllocateN(t *testing.T) {
	testAllocate(t, 1000000, 10)
}

func testAllocate(t *testing.T, n, size int) {
	s := getScode()
	start := time.Now()
	workers := 100
	codes := make([]map[string]bool, workers)
	for i := 0; i < workers; i++ {
		codes[i] = map[string]bool{}
	}
	perWorkN := n / workers
	var wg sync.WaitGroup
	wg.Add(workers)
	for j := 0; j < workers; j++ {
		go func(j int) {
			defer wg.Done()
			for i := perWorkN * j; i < perWorkN*(j+1); i++ {
				code, err := s.Allocate(context.Background(), size)
				if err != nil {
					t.Error(err)
				}
				codes[j][code] = true
			}
		}(j)
	}
	wg.Wait()
	allocateTime := time.Since(start)

	start = time.Now()
	wg.Add(workers)
	for j := 0; j < workers; j++ {
		go func(j int) {
			defer wg.Done()
			for code := range codes[j] {
				s.Release(context.Background(), code)
			}
		}(j)
	}
	wg.Wait()
	releaseTime := time.Since(start)

	fmt.Printf("| allocate | %f |\n", float64(n)/allocateTime.Seconds())
	fmt.Printf("| release | %f |\n", float64(n)/releaseTime.Seconds())
}

func TestScode_Collision(t *testing.T) {
	testCollision(t, 8, 10000, 10000)
	testCollision(t, 8, 100000, 10000)
	testCollision(t, 9, 100000, 10000)
	testCollision(t, 9, 1000000, 10000)
	testCollision(t, 10, 1000000, 10000)
	testCollision(t, 10, 10000000, 10000)
}

func testCollision(t *testing.T, size, n, allocateTimes int) {
	s := getScode()
	codes := map[string]bool{}
	for i := 0; i < n; i++ {
		code, err := s.Allocate(context.Background(), size)
		if err != nil {
			t.Error(err)
		}
		codes[code] = true
	}
	s = getScode()

	for i := 0; i < allocateTimes; i++ {
		code, err := s.Allocate(context.Background(), size)
		if err != nil {
			t.Error(err)
		}
		s.Release(context.Background(), code)
	}

	for code := range codes {
		s.Release(context.Background(), code)
	}

	stats := s.Stats()
	fmt.Printf("| %d | %d | %f | %f |\n", size, n,
		(float64(stats.AllocateTryTimes)/float64(stats.AllocateTimes)-1)*100,
		(float64(stats.AllocateFailureTimes)/float64(stats.AllocateTimes))*100)
}

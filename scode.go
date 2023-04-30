package scode

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/jiaxwu/scode/store"
)

var ErrTryTooManyTimes = errors.New("try too many times")

type OptionFunc func(score *Scode)

func WithAllowFunc(allowFunc AllowFunc) OptionFunc {
	return func(score *Scode) {
		score.allowFunc = allowFunc
	}
}

func WithTryTimes(tryTimes int) OptionFunc {
	return func(score *Scode) {
		score.tryTimes = tryTimes
	}
}

type Stats struct {
	AllocateTimes        uint64
	AllocateTryTimes     uint64
	AllocateFailureTimes uint64
}

type Scode struct {
	store     store.Store
	allowFunc AllowFunc
	charset   string // charset for short code
	rand      *rand.Rand
	tryTimes  int // times of try to set
	stats     Stats
}

func New(store store.Store, charset string, options ...OptionFunc) *Scode {
	score := &Scode{
		store:   store,
		charset: charset,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, option := range options {
		option(score)
	}
	if score.tryTimes < 1 {
		score.tryTimes = 1
	}
	if score.allowFunc == nil {
		score.allowFunc = AllowAll
	}
	return score
}

// Allocate allocate an short code
func (s *Scode) Allocate(ctx context.Context, size int) (string, error) {
	atomic.AddUint64(&s.stats.AllocateTimes, 1)
	for i := 0; i < s.tryTimes; i++ {
		atomic.AddUint64(&s.stats.AllocateTryTimes, 1)
		code := make([]byte, size)
		for i := 0; i < size; i++ {
			code[i] = s.charset[rand.Intn(size)]
		}
		scode := string(code)
		if !s.allowFunc(scode) {
			continue
		}
		set, err := s.store.SetIfNotExists(ctx, scode)
		if err != nil {
			log.Printf("set code error|%s", err.Error())
			continue
		}
		if !set {
			continue
		}
		return scode, nil
	}
	atomic.AddUint64(&s.stats.AllocateFailureTimes, 1)
	return "", ErrTryTooManyTimes
}

// Release release an short code
func (s *Scode) Release(ctx context.Context, code string) error {
	return s.store.Delete(ctx, code)
}

func (s *Scode) Stats() Stats {
	return s.stats
}

package store

import (
	"context"
)

type Store interface {
	SetIfNotExists(ctx context.Context, code string) (bool, error)
	Delete(ctx context.Context, code string) error
}

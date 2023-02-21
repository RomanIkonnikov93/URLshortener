package storage

import (
	"context"

	"github.com/RomanIkonnikov93/URLshortner/internal/model"
)

type Storage interface {
	Add(ctx context.Context, short, long, id string) error
	Get(ctx context.Context, short string) (string, error)
	GetByUserID(ctx context.Context, id string) (map[string]string, error)
	GetShort(ctx context.Context, long string) (string, error)
	BatchDelete(batch model.UserRequest) error
}

package repository

import (
	"github.com/RomanIkonnikov93/URLshortner/cmd/config"
	"github.com/RomanIkonnikov93/URLshortner/internal/repository/storage"
	"github.com/RomanIkonnikov93/URLshortner/internal/repository/users"
)

type Pool struct {
	Users   *users.Repository
	Storage *storage.Repository
	Ping    *Ping
}

func NewReps(cfg config.Config) (*Pool, error) {

	u, err := users.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	s, err := storage.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	p, err := NewPing(cfg)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Users:   u,
		Storage: s,
		Ping:    p,
	}, nil
}

package repository

import (
	"context"

	"github.com/RomanIkonnikov93/URLshortner/cmd/config"
	"github.com/RomanIkonnikov93/URLshortner/internal/conn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Pinger interface {
	PingDB() error
}

type Ping struct {
	pool *pgxpool.Pool
}

func NewPing(cfg config.Config) (*Ping, error) {

	pool := conn.NewConnection(cfg)

	return &Ping{
		pool: pool,
	}, nil
}

func (p *Ping) PingDB() error {

	pool := p.pool
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := pool.Ping(ctx)

	return err
}

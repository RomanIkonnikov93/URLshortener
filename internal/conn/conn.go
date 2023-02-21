package conn

import (
	"context"

	"github.com/RomanIkonnikov93/URLshortner/cmd/config"
	"github.com/RomanIkonnikov93/URLshortner/internal/model"
	"github.com/RomanIkonnikov93/URLshortner/logging"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewConnection(cfg config.Config) *pgxpool.Pool {

	logger := logging.GetLogger()

	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, cfg.DSN)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}

	return pool
}

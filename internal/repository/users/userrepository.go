package users

import (
	"context"

	"github.com/RomanIkonnikov93/URLshortner/cmd/config"
	"github.com/RomanIkonnikov93/URLshortner/internal/conn"
	"github.com/RomanIkonnikov93/URLshortner/internal/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(cfg config.Config) (*Repository, error) {

	pool := conn.NewConnection(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()

	if _, err := pool.Exec(ctx, `
	create table if not exists users (	    
	    user_id varchar(16) unique
	)

`); err != nil {
		return nil, err
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (p *Repository) AddUserID(user string) error {

	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()

	if _, err := p.pool.Exec(ctx, `insert into users (user_id) values ($1)`, user); err != nil {
		return err
	}

	return nil
}

func (p *Repository) CheckUserID(user string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()

	if _, err := p.pool.Exec(ctx, `insert into users (user_id) values ($1)`, user); err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				return true, nil
			}
		}
	}
	return false, nil
}

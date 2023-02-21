package storage

import (
	"context"
	"strconv"

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
	create table if not exists urls (
	    short varchar(5),
	    long text unique,
	    user_id varchar(16),
	    del_flag boolean
	)

`); err != nil {
		return nil, err
	}

	p := &Repository{
		pool: pool,
	}

	return p, nil
}

func (p *Repository) Add(ctx context.Context, short, long, id string) error {

	flag := false
	if _, err := p.pool.Exec(ctx, `insert into urls (short, long, user_id, del_flag) values ($1, $2 ,$3,$4)`, short, long, id, flag); err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				return model.ErrConflict
			}
		}

		return err
	}

	return nil
}

func (p *Repository) Get(ctx context.Context, short string) (string, error) {

	rows, err := p.pool.Query(ctx, `select long, del_flag from urls where short = $1`, short)
	if err != nil {
		return "", err
	}
	var out string
	var flag bool
	for rows.Next() {
		if err := rows.Scan(&out, &flag); err != nil {
			return "", err
		}
	}
	if flag {
		return "", model.ErrDelFlag
	}
	return out, nil
}

func (p *Repository) GetByUserID(ctx context.Context, id string) (map[string]string, error) {

	rows, err := p.pool.Query(ctx, `select short, long from urls where user_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]string)

	for rows.Next() {
		var short, long string
		if err := rows.Scan(&short, &long); err != nil {
			return nil, err
		}
		m[short] = long
	}

	return m, nil
}

func (p *Repository) GetShort(ctx context.Context, long string) (string, error) {

	row := p.pool.QueryRow(ctx, `select short from urls where long = $1`, long)
	var out string
	if err := row.Scan(&out); err != nil {
		return "", err
	}

	return out, nil
}

func (p *Repository) BatchDelete(batch model.UserRequest) error {

	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()

	query := `update urls set del_flag='true' where `
	args := make([]interface{}, 0)
	i := 1
	for _, val := range batch.UserUrls {
		query += `((user_id=$` + strconv.Itoa(i) + `) and (short=$` + strconv.Itoa(i+1) + `)) OR`
		i += 2
		args = append(args, batch.UserID, val)
	}
	query = query[:len(query)-2]
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

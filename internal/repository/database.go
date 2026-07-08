package repository


import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
    Pool *pgxpool.Pool
}

func NewDatabase(
    ctx context.Context,
    conn string,
) (*Database, error) {

    pool, err := pgxpool.New(ctx, conn)
    if err != nil {
        return nil, err
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, err
    }

    return &Database{
        Pool: pool,
    }, nil
}

func (d *Database) Close() {
    d.Pool.Close()
}
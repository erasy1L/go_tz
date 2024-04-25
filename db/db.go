package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func NewDatabase(ctx context.Context, connString string) (*Database, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		conn.Close(ctx)
	}()

	return &Database{Conn: conn}, nil
}

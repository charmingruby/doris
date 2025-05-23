package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Database interface {
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Preparex(query string) (*sqlx.Stmt, error)
}

type Client struct {
	Conn   *sqlx.DB
	logger *instrumentation.Logger
}

type ConnectionInput struct {
	User         string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SSL          string
}

func New(logger *instrumentation.Logger, in ConnectionInput) (*Client, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		in.User,
		in.Password,
		in.Host,
		in.Port,
		in.DatabaseName,
		in.SSL,
	)

	dbDriver := "postgres"

	db, err := sqlx.Connect(dbDriver, connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Client{Conn: db, logger: logger}, nil
}

func (c *Client) Close(ctx context.Context) error {
	c.Conn.Close()

	select {
	case <-ctx.Done():
		c.logger.Error("failed to close postgres connection", "error", ctx.Err())
		return ctx.Err()
	default:
		c.logger.Debug("postgres connection closed")
		return nil
	}
}

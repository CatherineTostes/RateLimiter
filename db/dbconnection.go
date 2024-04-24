package db

import "context"

type (
	DBConnection interface {
		Connect()
		Get(ctx context.Context, key string) (int64, error)
		Set(ctx context.Context, key string, value string) error
		Incr(ctx context.Context, key string) (uint64, error)
	}

	Connection struct {
		Db DBConnection
	}
)

func (c *Connection) Connect() {
	c.Db.Connect()
}

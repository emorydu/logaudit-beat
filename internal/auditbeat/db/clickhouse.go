package db

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

// ClickhouseOptions defines options for the Clickhouse database.
type ClickhouseOptions struct {
	Host                  []string
	Database              string
	Username              string
	Password              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
}

// NewClickhouse creates a new instance of the gorm database with the given options.
func NewClickhouse(opts *ClickhouseOptions) (driver.Conn, error) {
	if opts.MaxIdleConnections == 0 {
		opts.MaxIdleConnections = 5
	}

	if opts.MaxConnectionLifeTime == 0 {
		opts.MaxConnectionLifeTime = time.Hour
	}

	if opts.MaxOpenConnections == 0 {
		opts.MaxOpenConnections = 10
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: opts.Host,
		Auth: clickhouse.Auth{
			Database: opts.Database,
			Username: opts.Username,
			Password: opts.Password,
		},

		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		MaxOpenConns:    opts.MaxOpenConnections,
		MaxIdleConns:    opts.MaxIdleConnections,
		ConnMaxLifetime: opts.MaxConnectionLifeTime,
	})

	if err != nil {
		return nil, err
	}
	return conn, nil
}

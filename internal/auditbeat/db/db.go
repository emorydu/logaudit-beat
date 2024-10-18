// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLOptions defines options for the MySQL database.
type MySQLOptions struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
}

// DSN returns the DSN from MySQLOptions.
func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")
}

// NewMySQL creates a new instance of the gorm database with the given options.
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	logLevel := logger.Silent
	if opts.LogLevel != 0 {
		logLevel = logger.LogLevel(opts.LogLevel)
	}
	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns set the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime setting the maximum time a connection can be reused
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns setting the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

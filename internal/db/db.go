package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Open(driver, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("empty DB_DSN")
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// 连接池参数
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

// Ping Ping用于readiness：确认数据库可连通且可用
func Ping(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

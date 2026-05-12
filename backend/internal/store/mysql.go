package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

// Store 封装 MySQL 连接池，是所有数据库操作的入口。
type Store struct {
	db *sql.DB
}

// NewMySQLStore 创建数据库连接池并验证连通性。
// 连接池参数可通过环境变量调优（DB_MAX_OPEN_CONNS 等），适应不同部署规模。
func NewMySQLStore(ctx context.Context, dsn string) (*Store, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	// 强制 ParseTime + UTC：确保时间字段被正确解析，避免时区混乱
	cfg.ParseTime = true
	cfg.Loc = time.UTC

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	// 连接池配置：从环境变量读取，使用合理默认值
	db.SetMaxOpenConns(envIntOrDefault("DB_MAX_OPEN_CONNS", 20))
	db.SetMaxIdleConns(envIntOrDefault("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(envDurationOrDefault("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	db.SetConnMaxIdleTime(envDurationOrDefault("DB_CONN_MAX_IDLE_TIME", 5*time.Minute))

	// 启动时验证连接，快速发现配置错误
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() {
	_ = s.db.Close()
}

// Ping 健康检查探针使用，验证数据库连接是否存活。
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func envIntOrDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}

	return n
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return fallback
	}

	return d
}
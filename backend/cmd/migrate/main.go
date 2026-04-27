package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

func main() {
	timeout := flag.Duration("timeout", 30*time.Second, "migration timeout")
	flag.Parse()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	files := flag.Args()
	if len(files) == 0 {
		files = []string{
			"migrations/001_init.sql",
			"migrations/002_add_campaign_status.sql",
		}
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("parse dsn: %v", err)
	}
	cfg.ParseTime = true
	cfg.Loc = time.UTC

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping mysql: %v", err)
	}

	for _, file := range files {
		if err := runFile(ctx, db, file); err != nil {
			log.Fatalf("run %s: %v", file, err)
		}
		log.Printf("applied %s", file)
	}
}

func runFile(ctx context.Context, db *sql.DB, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	statements := splitStatements(normalizeSQL(string(raw)))
	for _, stmt := range statements {
		if stmt == "" {
			continue
		}
		compatible := makeCompatible(stmt)
		if _, err := db.ExecContext(ctx, compatible); err != nil {
			if ignorableMySQLError(err) {
				continue
			}
			return fmt.Errorf("exec sql from %s: %w", filepath.Base(path), err)
		}
	}

	return nil
}

func normalizeSQL(in string) string {
	lines := strings.Split(in, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func splitStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt != "" {
			out = append(out, stmt)
		}
	}
	return out
}

func makeCompatible(stmt string) string {
	stmt = strings.ReplaceAll(stmt, "CREATE INDEX IF NOT EXISTS", "CREATE INDEX")
	stmt = strings.ReplaceAll(stmt, "ADD COLUMN IF NOT EXISTS", "ADD COLUMN")
	return stmt
}

func ignorableMySQLError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}

	switch mysqlErr.Number {
	case 1060, 1061:
		return true
	default:
		return false
	}
}

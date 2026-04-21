package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crowdfunding/backend/internal/chain"
	"crowdfunding/backend/internal/config"
	"crowdfunding/backend/internal/indexer"
	"crowdfunding/backend/internal/store"
)

func main() {
	cfgPath := envOrDefault("CHAIN_CONFIG_PATH", "config/chain.testnet.example.json")
	dbDSN := os.Getenv("DATABASE_URL")
	if dbDSN == "" {
		log.Fatal("DATABASE_URL is required")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := chain.DialHTTP(startupCtx, cfg.RPCHTTPURL)
	if err != nil {
		log.Fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	st, err := store.NewPostgresStore(startupCtx, dbDSN)
	if err != nil {
		log.Fatalf("new store: %v", err)
	}
	defer st.Close()

	svc, err := indexer.New(cfg, client, st)
	if err != nil {
		log.Fatalf("new indexer: %v", err)
	}

	log.Printf("indexer started for contract %s", cfg.ContractAddress)
	if err := svc.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("indexer stopped with error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"crowdfunding/backend/internal/api"
	"crowdfunding/backend/internal/chain"
	"crowdfunding/backend/internal/config"
	"crowdfunding/backend/internal/store"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfgPath := envOrDefault("CHAIN_CONFIG_PATH", "config/chain.testnet.example.json")
	dbDSN := os.Getenv("DATABASE_URL")
	if dbDSN == "" {
		log.Fatal("DATABASE_URL is required")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client, err := chain.DialHTTP(ctx, cfg.RPCHTTPURL)
	if err != nil {
		log.Fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	st, err := store.NewMySQLStore(ctx, dbDSN)
	if err != nil {
		log.Fatalf("new store: %v", err)
	}
	defer st.Close()

	reader, err := chain.NewCrowdFundReader(common.HexToAddress(cfg.ContractAddress), client)
	if err != nil {
		log.Fatalf("new crowd reader: %v", err)
	}

	h := api.NewHandler(st, reader)
	addr := envOrDefault("API_ADDR", ":8080")
	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatalf("serve api: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

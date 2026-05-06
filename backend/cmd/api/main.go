package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crowdfunding/backend/internal/api"
	"crowdfunding/backend/internal/chain"
	"crowdfunding/backend/internal/config"
	"crowdfunding/backend/internal/store"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

	client, err := chain.DialHTTP(startupCtx, cfg.RPCHTTPURL)
	if err != nil {
		log.Fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	st, err := store.NewMySQLStore(startupCtx, dbDSN)
	if err != nil {
		log.Fatalf("new store: %v", err)
	}
	defer st.Close()

	reader, err := chain.NewCrowdFundReader(common.HexToAddress(cfg.ContractAddress), client)
	if err != nil {
		log.Fatalf("new crowd reader: %v", err)
	}

	h := api.NewHandler(st, reader, cfg.Public())
	addr := envOrDefault("API_ADDR", ":8080")
	log.Printf("api listening on %s", addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           h.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown api: %v", err)
		}
	}()

	if err := <-errCh; err != nil && err != http.ErrServerClosed {
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
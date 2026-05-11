package main

import (
	"context"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfgPath := envOrDefault("CHAIN_CONFIG_PATH", "config/chain.testnet.example.json")
	dbDSN := os.Getenv("DATABASE_URL")
	if dbDSN == "" {
		logger.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Error("load config failed", "error", err)
		os.Exit(1)
	}

	client, err := chain.DialHTTP(startupCtx, cfg.RPCHTTPURL)
	if err != nil {
		logger.Error("dial rpc failed", "error", err, "url", cfg.RPCHTTPURL)
		os.Exit(1)
	}
	defer client.Close()

	st, err := store.NewMySQLStore(startupCtx, dbDSN)
	if err != nil {
		logger.Error("connect database failed", "error", err)
		os.Exit(1)
	}
	defer st.Close()

	reader, err := chain.NewCrowdFundReader(common.HexToAddress(cfg.ContractAddress), client)
	if err != nil {
		logger.Error("create chain reader failed", "error", err)
		os.Exit(1)
	}

	h := api.NewHandler(st, reader, cfg.Public())
	addr := envOrDefault("API_ADDR", ":8080")

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

	logger.Info("api server started",
		"addr", addr,
		"chain", cfg.ChainName,
		"contract", cfg.ContractAddress,
	)

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown error", "error", err)
		}
	}()

	if err := <-errCh; err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
	logger.Info("api server stopped gracefully")
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
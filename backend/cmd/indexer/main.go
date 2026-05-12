// Indexer 服务入口。负责持续监听链上 CrowdFund 合约事件并同步到 MySQL。
// 与 API 服务分离部署：同步失败不影响查询服务，两者可独立扩容。
package main

import (
	"context"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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

	svc, err := indexer.New(cfg, client, st)
	if err != nil {
		logger.Error("create indexer failed", "error", err)
		os.Exit(1)
	}

	logger.Info("indexer started",
		"contract", cfg.ContractAddress,
		"chain", cfg.ChainName,
		"startBlock", cfg.DeploymentStartBlock,
	)

	// Run 阻塞直到 context 取消或遇到不可恢复的错误
	if err := svc.Run(ctx); err != nil && err != context.Canceled {
		logger.Error("indexer stopped with error", "error", err)
		os.Exit(1)
	}
	logger.Info("indexer stopped gracefully")
}

func envOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
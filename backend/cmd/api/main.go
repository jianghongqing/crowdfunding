// API 服务入口。负责启动 HTTP 服务器，连接数据库和链 RPC。
// 启动流程：加载配置 -> 连接 RPC -> 连接 MySQL -> 注册路由 -> 监听端口。
// 优雅关闭：收到 SIGINT/SIGTERM 后等待在途请求完成（最多 10s）再退出。
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

	// 监听系统信号，用于优雅关闭
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 启动阶段的超时保护：15s 内必须完成所有连接初始化
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

	// 链读取器供 API 在数据库未命中时回退到链上查询
	reader, err := chain.NewCrowdFundReader(common.HexToAddress(cfg.ContractAddress), client)
	if err != nil {
		logger.Error("create chain reader failed", "error", err)
		os.Exit(1)
	}

	h := api.NewHandler(st, reader, cfg.Public())
	addr := envOrDefault("API_ADDR", ":8080")

	// 超时配置防止慢客户端占用连接
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

	// 优雅关闭：等待 context 取消后 shutdown，给在途请求最多 10s 完成
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
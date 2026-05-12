// Package api 提供 HTTP API 处理层，负责请求校验、路由分发和响应格式化。
// 所有业务数据优先从 MySQL 读取，未命中时回退到链上合约调用。
package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"crowdfunding/backend/internal/chain"
	"crowdfunding/backend/internal/config"
	"crowdfunding/backend/internal/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler 聚合 API 所需的所有依赖：数据库、链读取器、公共配置。
type Handler struct {
	store  *store.Store
	reader *chain.CrowdFundReader
	cfg    config.PublicChainConfig
	log    *slog.Logger
}

// NewHandler 创建 Handler 实例并初始化结构化日志（JSON 格式输出到 stdout）。
func NewHandler(st *store.Store, reader *chain.CrowdFundReader, cfg config.PublicChainConfig) *Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return &Handler{store: st, reader: reader, cfg: cfg, log: logger}
}

// Routes 组装完整的中间件链和路由表。
// 中间件执行顺序：RequestID -> RealIP -> slog日志 -> panic恢复 -> 超时 -> 安全头 -> CORS -> 限流。
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// 每 IP 30 req/s，峰值 burst 60，防止滥用
	rl := newRateLimiter(30, 60)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(slogMiddleware(h.log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(securityHeaders)
	r.Use(corsMiddleware)
	r.Use(rateLimitMiddleware(rl))

	// 基础设施端点（不带版本前缀）
	r.Get("/healthz", h.health)
	r.Get("/config", h.getConfig)

	// v1 版本化 API
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/campaigns", h.listCampaigns)
		api.Get("/campaigns/{id}", h.getCampaign)
		api.Get("/campaigns/{id}/contributions/{address}", h.getContributionForAddress)
		api.Get("/campaigns/{id}/contributions", h.listContributions)
		api.Get("/stats", h.getStats)
	})

	// 向后兼容：保留不带前缀的旧路由，便于已有前端平滑迁移
	r.Get("/campaigns", h.listCampaigns)
	r.Get("/campaigns/{id}", h.getCampaign)
	r.Get("/campaigns/{id}/contributions/{address}", h.getContributionForAddress)
	r.Get("/campaigns/{id}/contributions", h.listContributions)

	return r
}

// health 健康检查端点，验证数据库连通性。容器编排和负载均衡依赖此端点判断服务可用性。
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.store.Ping(ctx); err != nil {
		h.log.Error("health check failed", "error", err)
		respondError(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database unavailable")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// getConfig 返回脱敏后的公共链配置（不含 RPC URL 等敏感信息），供前端自动加载合约地址。
func (h *Handler) getConfig(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.cfg)
}

// getStats 返回平台级聚合统计数据（活动总数、各状态数量、总筹款额）。
func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetStats(r.Context())
	if err != nil {
		h.log.Error("get stats failed", "error", err)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to load statistics")
		return
	}
	respondJSON(w, http.StatusOK, stats)
}

// listCampaigns 分页查询活动列表，支持按 status 筛选。
// 返回值包含 pagination 元数据（total / hasMore），方便前端实现翻页。
func (h *Handler) listCampaigns(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	status := r.URL.Query().Get("status")
	items, total, err := h.store.ListCampaignsWithCount(r.Context(), limit, offset, status)
	if err != nil {
		h.log.Error("list campaigns failed", "error", err)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list campaigns")
		return
	}

	respondPaginated(w, items, total, limit, offset)
}

// getCampaign 查询单个活动详情。
// 策略：先查 MySQL（indexer 同步的快照），未命中则回退到链上 eth_call 兜底。
func (h *Handler) getCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PARAMS", "invalid campaign id")
		return
	}

	item, err := h.store.GetCampaign(r.Context(), id)
	if err == nil {
		respondJSON(w, http.StatusOK, item)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		h.log.Error("get campaign failed", "error", err, "campaignId", id)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get campaign")
		return
	}

	// 数据库没有则回链上读，可能是 indexer 尚未同步到该区块
	view, chainErr := h.reader.GetCampaign(r.Context(), id)
	if chainErr != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "campaign not found")
		return
	}

	respondJSON(w, http.StatusOK, view)
}

// getContributionForAddress 查询某地址在某活动中的当前捐款额（直接读链上 mapping）。
func (h *Handler) getContributionForAddress(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PARAMS", "invalid campaign id")
		return
	}

	addr := chi.URLParam(r, "address")
	if !common.IsHexAddress(addr) {
		respondError(w, http.StatusBadRequest, "INVALID_PARAMS", "invalid ethereum address")
		return
	}

	amount, err := h.reader.GetContribution(r.Context(), id, common.HexToAddress(addr))
	if err != nil {
		h.log.Error("get contribution failed", "error", err, "campaignId", id, "address", addr)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get contribution")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"campaignId": strconv.FormatUint(id, 10),
		"address":    addr,
		"amountWei":  amount,
	})
}

// listContributions 分页查询某活动的捐款明细（从 MySQL 读取，由 indexer 同步写入）。
func (h *Handler) listContributions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PARAMS", "invalid campaign id")
		return
	}

	limit, offset, err := parsePagination(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PARAMS", err.Error())
		return
	}

	items, total, err := h.store.ListContributionsWithCount(r.Context(), id, limit, offset)
	if err != nil {
		h.log.Error("list contributions failed", "error", err, "campaignId", id)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list contributions")
		return
	}

	respondPaginated(w, items, total, limit, offset)
}

// ---- 统一响应格式 ----

// APIError 标准化错误响应体，code 为机器可读错误码，message 为人类可读描述。
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PaginatedResponse 分页列表的标准响应格式，data + pagination 分离。
type PaginatedResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta 分页元数据，hasMore 方便前端判断是否还有下一页。
type PaginationMeta struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"hasMore"`
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, httpCode int, errCode, message string) {
	respondJSON(w, httpCode, map[string]any{
		"error": APIError{Code: errCode, Message: message},
	})
}

// respondErr 简化版错误响应，供中间件调用（中间件不关心具体错误码）。
func respondErr(w http.ResponseWriter, code int, msg string) {
	respondError(w, code, "ERROR", msg)
}

func respondPaginated(w http.ResponseWriter, data any, total, limit, offset int) {
	resp := PaginatedResponse{
		Data: data,
		Pagination: PaginationMeta{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: offset+limit < total,
		},
	}
	respondJSON(w, http.StatusOK, resp)
}

// parsePagination 从 query string 解析分页参数，限制 limit 上限为 100 防止大查询。
func parsePagination(r *http.Request) (limit, offset int, err error) {
	limit = 20
	offset = 0

	if raw := r.URL.Query().Get("limit"); raw != "" {
		limit, err = strconv.Atoi(raw)
		if err != nil || limit <= 0 || limit > 100 {
			return 0, 0, errors.New("limit must be between 1 and 100")
		}
	}

	if raw := r.URL.Query().Get("offset"); raw != "" {
		offset, err = strconv.Atoi(raw)
		if err != nil || offset < 0 {
			return 0, 0, errors.New("offset must be a non-negative integer")
		}
	}

	return limit, offset, nil
}

// slogMiddleware 替代 chi 内置 Logger，输出 JSON 结构化日志，便于日志收集系统解析。
func slogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration_ms", time.Since(start).Milliseconds(),
				"ip", r.RemoteAddr,
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
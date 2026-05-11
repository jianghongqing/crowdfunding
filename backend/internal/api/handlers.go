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

type Handler struct {
	store  *store.Store
	reader *chain.CrowdFundReader
	cfg    config.PublicChainConfig
	log    *slog.Logger
}

func NewHandler(st *store.Store, reader *chain.CrowdFundReader, cfg config.PublicChainConfig) *Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return &Handler{store: st, reader: reader, cfg: cfg, log: logger}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	rl := newRateLimiter(30, 60)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(slogMiddleware(h.log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(securityHeaders)
	r.Use(corsMiddleware)
	r.Use(rateLimitMiddleware(rl))

	r.Get("/healthz", h.health)
	r.Get("/config", h.getConfig)

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/campaigns", h.listCampaigns)
		api.Get("/campaigns/{id}", h.getCampaign)
		api.Get("/campaigns/{id}/contributions/{address}", h.getContributionForAddress)
		api.Get("/campaigns/{id}/contributions", h.listContributions)
		api.Get("/stats", h.getStats)
	})

	r.Get("/campaigns", h.listCampaigns)
	r.Get("/campaigns/{id}", h.getCampaign)
	r.Get("/campaigns/{id}/contributions/{address}", h.getContributionForAddress)
	r.Get("/campaigns/{id}/contributions", h.listContributions)

	return r
}

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

func (h *Handler) getConfig(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.cfg)
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetStats(r.Context())
	if err != nil {
		h.log.Error("get stats failed", "error", err)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to load statistics")
		return
	}
	respondJSON(w, http.StatusOK, stats)
}

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

	view, chainErr := h.reader.GetCampaign(r.Context(), id)
	if chainErr != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "campaign not found")
		return
	}

	respondJSON(w, http.StatusOK, view)
}

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

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PaginatedResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

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
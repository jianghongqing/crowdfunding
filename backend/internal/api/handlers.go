package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
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
}

func NewHandler(store *store.Store, reader *chain.CrowdFundReader, cfg config.PublicChainConfig) *Handler {
	return &Handler{store: store, reader: reader, cfg: cfg}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(corsMiddleware)
	r.Use(jsonContentType)
	r.Get("/config", h.getConfig)
	r.Get("/healthz", h.health)
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
		respondErr(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listCampaigns(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.store.ListCampaigns(r.Context(), limit, offset)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondErr(w, http.StatusBadRequest, "invalid campaign id")
		return
	}

	item, err := h.store.GetCampaign(r.Context(), id)
	if err == nil {
		respondJSON(w, http.StatusOK, item)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	view, chainErr := h.reader.GetCampaign(r.Context(), id)
	if chainErr != nil {
		respondErr(w, http.StatusNotFound, chainErr.Error())
		return
	}

	respondJSON(w, http.StatusOK, view)
}

func (h *Handler) getContributionForAddress(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondErr(w, http.StatusBadRequest, "invalid campaign id")
		return
	}

	addr := chi.URLParam(r, "address")
	if !common.IsHexAddress(addr) {
		respondErr(w, http.StatusBadRequest, "invalid address")
		return
	}

	amount, err := h.reader.GetContribution(r.Context(), id, common.HexToAddress(addr))
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
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
		respondErr(w, http.StatusBadRequest, "invalid campaign id")
		return
	}

	limit, offset, err := parsePagination(r)
	if err != nil {
		respondErr(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.store.ListContributions(r.Context(), id, limit, offset)
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, items)
}

func (h *Handler) getConfig(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, h.cfg)
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func respondErr(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, code, map[string]string{"error": msg})
}

func jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
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
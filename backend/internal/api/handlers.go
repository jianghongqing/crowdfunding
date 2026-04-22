package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"crowdfunding/backend/internal/chain"
	"crowdfunding/backend/internal/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store  *store.Store
	reader *chain.CrowdFundReader
}

func NewHandler(store *store.Store, reader *chain.CrowdFundReader) *Handler {
	return &Handler{store: store, reader: reader}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", h.health)
	r.Get("/campaigns", h.listCampaigns)
	r.Get("/campaigns/{id}", h.getCampaign)
	r.Get("/campaigns/{id}/contributions/{address}", h.getContributionForAddress)
	return r
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listCampaigns(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.ListCampaigns(r.Context(), 50, 0)
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
	respondJSON(w, http.StatusOK, map[string]string{"campaignId": strconv.FormatUint(id, 10), "address": addr, "amountWei": amount})
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func respondErr(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, code, map[string]string{"error": msg})
}

package auction

import (
	"encoding/json"
	"errors"
	"net/http"
)

type handler struct {
	svc Service
}

func NewHandler(svc Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) listAuctions(wr http.ResponseWriter, req *http.Request) {
	auctions, err := h.svc.ListAuctions(req.Context(), true)
	if err != nil {
		// log.Printf("listAuctions error: %v", err)
		http.Error(wr, "internal server error", http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(auctions)
}

func (h *handler) getAuction(wr http.ResponseWriter, req *http.Request) {
	address := req.PathValue("address")

	auction, err := h.svc.GetAuction(req.Context(), address)
	if err != nil {
		if errors.Is(err, ErrAuctionNotFound) {
			http.Error(wr, "auction not found", http.StatusNotFound)
			return
		}
		http.Error(wr, "internal server error", http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(auction)
}

func (h *handler) listBids(wr http.ResponseWriter, req *http.Request) {
	address := req.PathValue("address")

	bids, err := h.svc.ListBids(req.Context(), address)
	if err != nil {
		http.Error(wr, "internal server error", http.StatusInternalServerError)
		return
	}

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(bids)
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auctions", h.listAuctions)
	mux.HandleFunc("GET /auctions/{address}", h.getAuction)
	mux.HandleFunc("GET /auctions/{address}/bids", h.listBids)
}

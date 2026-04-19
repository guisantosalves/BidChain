package auction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/guisantosalves/bidchain/internal/blockchain"
)

type handler struct {
	svc    Service
	caller *blockchain.Caller
}

func NewHandler(svc Service, caller *blockchain.Caller) *handler {
	return &handler{svc: svc, caller: caller}
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

func (h *handler) createAuction(wr http.ResponseWriter, req *http.Request) {
	var createAucReq CreateAuctionRequest
	if err := json.NewDecoder(req.Body).Decode(&createAucReq); err != nil {
		http.Error(wr, "Invalid request body", http.StatusBadRequest)
		return
	}

	if createAucReq.Description == "" || createAucReq.IPFSHash == "" || createAucReq.DurationSeconds == 0 {
		http.Error(wr, "missing required fields", http.StatusBadRequest)
		return
	}

	go func() {
		if err := h.caller.CreateAuction(
			context.Background(),
			createAucReq.Description,
			createAucReq.IPFSHash,
			createAucReq.DurationSeconds,
		); err != nil {
			log.Printf("createAuction: blockchain error: %v", err)
		}
	}()

	wr.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(wr, `{"status": "auction creation submitted"}`)
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auctions", h.listAuctions)
	mux.HandleFunc("GET /auctions/{address}", h.getAuction)
	mux.HandleFunc("GET /auctions/{address}/bids", h.listBids)
	mux.HandleFunc("POST /auctions", h.createAuction)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/guisantosalves/bidchain/internal/auction"
	"github.com/guisantosalves/bidchain/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("Failed to get database url")
	}

	db, err := database.New(dsn)
	if err != nil {
		log.Fatalf("Failed to get db instance")
	}
	defer db.Close() // make sure to do this

	auctionRepo := auction.NewAuctionRepository(db)
	bidRepo := auction.NewBidRepository(db)
	svc := auction.NewService(auctionRepo, bidRepo)
	h := auction.NewHandler(svc)

	fmt.Println("database connected")

	// opening and managing server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(wr http.ResponseWriter, req *http.Request) {
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte(`{"status": "ok"}`))
	})

	// Bid | auctions
	h.RegisterRoutes(mux)

	log.Println("servver running on 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %w", err)
	}
}

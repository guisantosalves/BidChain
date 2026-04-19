package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guisantosalves/bidchain/internal/auction"
	"github.com/guisantosalves/bidchain/internal/blockchain"
	"github.com/guisantosalves/bidchain/internal/database"
	"github.com/joho/godotenv"
)

func startServer(mux *http.ServeMux, cancel context.CancelFunc) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit // bloqueia até ctrl c

	log.Printf("Shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped")
}

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
	caller, err := blockchain.NewCaller(
		os.Getenv("CALLER_RPC_URL"),
		os.Getenv("PRIVATE_KEY"),
		blockchain.FactoryAddress,
	)
	if err != nil {
		log.Fatalf("failed to create caller: %v", err)
	}
	h := auction.NewHandler(svc, caller)

	fmt.Println("database connected")

	// opening and managing server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(wr http.ResponseWriter, req *http.Request) {
		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte(`{"status": "ok"}`))
	})

	// Bid | auctions
	h.RegisterRoutes(mux)

	// listener blockchain
	listener, err := blockchain.NewListener(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	// workers
	pool := blockchain.NewWorkerPool(5, listener.Events, func(ctx context.Context, event blockchain.AuctionCreatedEvent) error {
		// everytime the event receive value on chanel it will be called, the worker will do that
		auction := &auction.Auction{
			Address:     event.AuctionAddress.Hex(),
			Seller:      event.Seller.Hex(),
			Description: event.Descrition,
		}

		return svc.CreateAuction(ctx, auction)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go listener.Start(ctx)
	go pool.Start(ctx)

	startServer(mux, cancel)
}

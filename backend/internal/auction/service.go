package auction

import "context"

type Service interface {
	CreateAuction(ctx context.Context, a *Auction) error
	GetAuction(ctx context.Context, address string) (*Auction, error)
	ListAuctions(ctx context.Context, onlyActivce bool) ([]*Auction, error)
	CreateBid(ctx context.Context, b *Bid) error
	ListBids(ctx context.Context, auctionAddress string) ([]*Bid, error)
}

type service struct {
	auctionRepo AuctionRepository
	bidRepo     BidRepository
}

// CreateAuction implements [Service].
func (s *service) CreateAuction(ctx context.Context, a *Auction) error {
	return s.CreateAuction(ctx, a)
}

// CreateBid implements [Service].
func (s *service) CreateBid(ctx context.Context, b *Bid) error {
	return s.CreateBid(ctx, b)
}

// GetAuction implements [Service].
func (s *service) GetAuction(ctx context.Context, address string) (*Auction, error) {
	return s.GetAuction(ctx, address)
}

// ListAuctions implements [Service].
func (s *service) ListAuctions(ctx context.Context, onlyActivce bool) ([]*Auction, error) {
	return s.auctionRepo.Lists(ctx, onlyActivce)
}

// ListBids implements [Service].
func (s *service) ListBids(ctx context.Context, auctionAddress string) ([]*Bid, error) {
	return s.bidRepo.ListByAction(ctx, auctionAddress)
}

func NewService(auctionRepo AuctionRepository, bidRepo BidRepository) Service {
	return &service{auctionRepo: auctionRepo, bidRepo: bidRepo}
}

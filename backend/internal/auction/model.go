package auction

import "time"

type Auction struct {
	ID          int64
	Address     string
	Seller      string
	Description string
	IPFSHash    string
	CreatedAt   time.Time
	EndsAt      time.Time
	Finalized   bool
}

type Bid struct {
	ID             int64
	AuctionAddress string
	Bidder         string
	Amount         string
	TxHash         string
	BlockNumber    uint64
	CreatedAt      time.Time
}

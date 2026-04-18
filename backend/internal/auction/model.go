package auction

import "time"

type Auction struct {
	ID          int64     `json:"id"`
	Address     string    `json:"address"`
	Seller      string    `json:"seller"`
	Description string    `json:"sescription"`
	IPFSHash    string    `json:"IPFSHash"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	EndsAt      time.Time `json:"endsAt"`
	Finalized   bool      `json:"finalized"`
}

type Bid struct {
	ID             int64     `json:"id"`
	AuctionAddress string    `json:"auctionAddress"`
	Bidder         string    `json:"bidder"`
	Amount         string    `json:"amount"`
	TxHash         string    `json:"txHash"`
	BlockNumber    uint64    `json:"blockNumber"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

package auction

import (
	"context"
	"database/sql"
	"fmt"
)

type BidRepository interface {
	Create(ctx context.Context, b *Bid) error
	ListByAction(ctx context.Context, auctionAddress string) ([]*Bid, error)
}

type postgresBidRepo struct {
	db *sql.DB
}

// Create implements [BidRepository].
func (p *postgresBidRepo) Create(ctx context.Context, b *Bid) error {
	query := `
		INSERT INTO bids (auction_address, bidder, amount, tx_hash, block_number)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
	`

	err := p.db.QueryRowContext(ctx, query, b.AuctionAddress,
		b.Bidder,
		b.Amount,
		b.TxHash,
		b.BlockNumber).Scan(&b.ID, &b.CreatedAt)

	if err != nil {
		return fmt.Errorf("bidRepo.Create: %w", err)
	}

	return nil
}

// ListByAction implements [BidRepository].
func (p *postgresBidRepo) ListByAction(ctx context.Context, auctionAddress string) ([]*Bid, error) {
	query := `
		SELECT id, auction_address, bidder, amount, tx_hash, block_number, created_at
		FROM bids
		WHERE auction_address = $1
		ORDER BY block_number DESC
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("bidRepo.listByAction: %w", err)
	}
	defer rows.Close()

	var bids []*Bid
	for rows.Next() {
		b := &Bid{}
		if err := rows.Scan(&b.ID,
			&b.AuctionAddress,
			&b.Bidder,
			&b.Amount,
			&b.TxHash,
			&b.BlockNumber,
			&b.CreatedAt); err != nil {
			return nil, fmt.Errorf("bidRepo.listByAction: scan: %w", err)
		}

		bids = append(bids, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("bidRepo.listByAction: rows: %w", err)
	}

	return bids, nil
}

func NewBidRepository(db *sql.DB) BidRepository {
	return &postgresBidRepo{db: db}
}

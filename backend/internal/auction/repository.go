package auction

import (
	"context"
	"database/sql"
	"fmt"
)

type AuctionRepository interface {
	Create(ctx context.Context, a *Auction) error
	GetByAddress(ctx context.Context, address string) (*Auction, error)
	Lists(ctx context.Context, onlyActive bool) ([]*Auction, error)
	Finalize(ctx context.Context, address string) error
}

type postgresAuctionRepo struct {
	db *sql.DB
}

func (r *postgresAuctionRepo) Create(ctx context.Context, a *Auction) error {
	query := `
		INSERT INTO auctions (address, seller, description, ipfs_hash, ends_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, a.Address,
		a.Seller,
		a.Description,
		a.IPFSHash,
		a.EndsAt,
	).Scan(&a.ID, &a.CreatedAt)

	if err != nil {
		return fmt.Errorf("auctionRepo.Create: %w", err)
	}

	return nil
}

func (r *postgresAuctionRepo) GetByAddress(ctx context.Context, address string) (*Auction,
	error) {
	query := `
        SELECT id, address, seller, description, ipfs_hash, ends_at, finalized, created_at
        FROM auctions
        WHERE address = $1
	`
	a := &Auction{}

	err := r.db.QueryRowContext(ctx, query, address).Scan(&a.ID,
		&a.Address,
		&a.Seller,
		&a.Description,
		&a.IPFSHash,
		&a.EndsAt,
		&a.Finalized,
		&a.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAuctionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("auctionRepo.GetByAddress: %w", err)
	}

	return a, nil
}

func (r *postgresAuctionRepo) Lists(ctx context.Context, onlyActive bool) ([]*Auction, error) {
	query := `
		SELECT id, address, seller, description, ipfs_hash, ends_at, finalized, created_at
        FROM auctions
        WHERE ($1 = false OR (finalized = false AND ends_at > NOW()))
        ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, onlyActive)
	if err != nil {
		return nil, fmt.Errorf("auctionRepo.List: %w", err)
	}
	defer rows.Close()

	var auctions []*Auction
	for rows.Next() {
		a := &Auction{}
		if err := rows.Scan(&a.ID,
			&a.Address,
			&a.Seller,
			&a.Description,
			&a.IPFSHash,
			&a.EndsAt,
			&a.Finalized,
			&a.CreatedAt); err != nil {
			return nil, fmt.Errorf("auctionRepo.List: scan: %w", err)
		}
		auctions = append(auctions, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("auctionRepo.List: rows: %w", err)
	}

	return auctions, err
}

func (r *postgresAuctionRepo) Finalize(ctx context.Context, address string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("auctionRepo.Finalize: begin: %w", err)
	}
	defer tx.Rollback() // runs only if the commit doesnt run

	query := `
		UPDATE auctions
		SET finalized = true
		WHERE address = $1 AND finalized = false
	`

	result, err := tx.ExecContext(ctx, query, address)

	if err != nil {
		return fmt.Errorf("auctionRepo.Finalize: exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("auctionRepo.Finalize: rows affected: %w", err)
	}
	if rows == 0 {
		return ErrAuctionFinalized
	}

	return tx.Commit()
}

func NewAuctionRepository(db *sql.DB) AuctionRepository {
	return &postgresAuctionRepo{db: db}
}

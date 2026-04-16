package auction

import "errors"

var (
	ErrAuctionNotFound  = errors.New("auction not found")
	ErrAuctionFinalized = errors.New("auction already finalized")
	ErrBidTooLow        = errors.New("bid amount too low")
)

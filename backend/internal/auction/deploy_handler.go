package auction

type CreateAuctionRequest struct {
	Description     string `json:"description"`
	IPFSHash        string `json:"ipfs_hash"`
	DurationSeconds uint64 `json:"duration_seconds"`
}

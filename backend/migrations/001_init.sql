CREATE TABLE IF NOT EXISTS users (
      id         SERIAL PRIMARY KEY,
      address    TEXT NOT NULL UNIQUE,
      username   TEXT,
      email      TEXT,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );

  CREATE TABLE IF NOT EXISTS auth_nonces (
      address    TEXT PRIMARY KEY,
      nonce      TEXT NOT NULL,
      expires_at TIMESTAMPTZ NOT NULL
  );

  CREATE TABLE IF NOT EXISTS auctions (
      id          SERIAL PRIMARY KEY,
      address     TEXT NOT NULL UNIQUE,
      seller      TEXT NOT NULL REFERENCES users(address),
      description TEXT NOT NULL,
      ipfs_hash   TEXT NOT NULL,
      ends_at     TIMESTAMPTZ NOT NULL,
      finalized   BOOLEAN NOT NULL DEFAULT FALSE,
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );

  CREATE TABLE IF NOT EXISTS bids (
      id              SERIAL PRIMARY KEY,
      auction_address TEXT NOT NULL REFERENCES auctions(address),
      bidder          TEXT NOT NULL REFERENCES users(address),
      amount          TEXT NOT NULL,
      tx_hash         TEXT NOT NULL UNIQUE,
      block_number    BIGINT NOT NULL,
      created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );

  CREATE INDEX IF NOT EXISTS idx_bids_auction ON bids(auction_address);
  CREATE INDEX IF NOT EXISTS idx_auctions_seller ON auctions(seller);
  CREATE INDEX IF NOT EXISTS idx_auctions_finalized ON auctions(finalized);
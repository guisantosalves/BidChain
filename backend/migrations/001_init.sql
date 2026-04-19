CREATE TABLE IF NOT EXISTS auctions (
      id          SERIAL PRIMARY KEY,
      address     TEXT NOT NULL UNIQUE,
      seller      TEXT NOT NULL,
      description TEXT NOT NULL,
      ipfs_hash   TEXT NOT NULL,
      ends_at     TIMESTAMPTZ NOT NULL,
      finalized   BOOLEAN NOT NULL DEFAULT FALSE,
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );

  CREATE TABLE IF NOT EXISTS bids (
      id              SERIAL PRIMARY KEY,
      auction_address TEXT NOT NULL,
      bidder          TEXT NOT NULL,
      amount          TEXT NOT NULL,
      tx_hash         TEXT NOT NULL UNIQUE,
      block_number    BIGINT NOT NULL,
      created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );
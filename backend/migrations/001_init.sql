CREATE TABLE IF NOT EXISTS campaigns (
    campaign_id BIGINT PRIMARY KEY,
    creator TEXT NOT NULL,
    title TEXT NOT NULL,
    goal_wei TEXT NOT NULL,
    pledged_wei TEXT NOT NULL,
    deadline BIGINT NOT NULL,
    withdrawn BOOLEAN NOT NULL DEFAULT FALSE,
    created_block BIGINT NOT NULL DEFAULT 0,
    created_tx_hash TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS contributions (
    id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL,
    funder TEXT NOT NULL,
    amount_wei TEXT NOT NULL,
    tx_hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    log_index BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tx_hash, log_index)
);

CREATE TABLE IF NOT EXISTS refunds (
    id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL,
    funder TEXT NOT NULL,
    amount_wei TEXT NOT NULL,
    tx_hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    log_index BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tx_hash, log_index)
);

CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL,
    creator TEXT NOT NULL,
    amount_wei TEXT NOT NULL,
    tx_hash TEXT NOT NULL,
    block_number BIGINT NOT NULL,
    log_index BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tx_hash, log_index)
);

CREATE TABLE IF NOT EXISTS indexer_checkpoints (
    worker_name TEXT PRIMARY KEY,
    last_scanned_block BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contributions_campaign_id ON contributions (campaign_id);
CREATE INDEX IF NOT EXISTS idx_refunds_campaign_id ON refunds (campaign_id);
CREATE INDEX IF NOT EXISTS idx_withdrawals_campaign_id ON withdrawals (campaign_id);

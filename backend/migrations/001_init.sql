CREATE TABLE IF NOT EXISTS campaigns (
    campaign_id BIGINT UNSIGNED PRIMARY KEY,
    creator VARCHAR(42) NOT NULL,
    title VARCHAR(255) NOT NULL,
    goal_wei VARCHAR(78) NOT NULL,
    pledged_wei VARCHAR(78) NOT NULL,
    deadline BIGINT UNSIGNED NOT NULL,
    withdrawn BOOLEAN NOT NULL DEFAULT FALSE,
    created_block BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_tx_hash VARCHAR(66) NOT NULL DEFAULT '',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS contributions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    campaign_id BIGINT UNSIGNED NOT NULL,
    funder VARCHAR(42) NOT NULL,
    amount_wei VARCHAR(78) NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    block_number BIGINT UNSIGNED NOT NULL,
    log_index BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_contributions_tx_log (tx_hash, log_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS refunds (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    campaign_id BIGINT UNSIGNED NOT NULL,
    funder VARCHAR(42) NOT NULL,
    amount_wei VARCHAR(78) NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    block_number BIGINT UNSIGNED NOT NULL,
    log_index BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_refunds_tx_log (tx_hash, log_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    campaign_id BIGINT UNSIGNED NOT NULL,
    creator VARCHAR(42) NOT NULL,
    amount_wei VARCHAR(78) NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    block_number BIGINT UNSIGNED NOT NULL,
    log_index BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uniq_withdrawals_tx_log (tx_hash, log_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS indexer_checkpoints (
    worker_name VARCHAR(64) PRIMARY KEY,
    last_scanned_block BIGINT UNSIGNED NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX IF NOT EXISTS idx_contributions_campaign_id ON contributions (campaign_id);
CREATE INDEX IF NOT EXISTS idx_refunds_campaign_id ON refunds (campaign_id);
CREATE INDEX IF NOT EXISTS idx_withdrawals_campaign_id ON withdrawals (campaign_id);

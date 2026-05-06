package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Store) InsertContribution(
	ctx context.Context, campaignID uint64, funder, amountWei, txHash string, blockNumber uint64, logIndex uint,
) error {
	const q = `
INSERT IGNORE INTO contributions (campaign_id, funder, amount_wei, tx_hash, block_number, log_index)
VALUES (?,?,?,?,?,?)`

	_, err := s.db.ExecContext(ctx, q, campaignID, funder, amountWei, txHash, blockNumber, logIndex)
	if err != nil {
		return fmt.Errorf("insert contribution: %w", err)
	}

	return nil
}

func (s *Store) InsertRefund(
	ctx context.Context, campaignID uint64, funder, amountWei, txHash string, blockNumber uint64, logIndex uint,
) error {
	const q = `
INSERT IGNORE INTO refunds (campaign_id, funder, amount_wei, tx_hash, block_number, log_index)
VALUES (?,?,?,?,?,?)`

	_, err := s.db.ExecContext(ctx, q, campaignID, funder, amountWei, txHash, blockNumber, logIndex)
	if err != nil {
		return fmt.Errorf("insert refund: %w", err)
	}

	return nil
}

func (s *Store) InsertWithdrawal(
	ctx context.Context, campaignID uint64, creator, amountWei, txHash string, blockNumber uint64, logIndex uint,
) error {
	const q = `
INSERT IGNORE INTO withdrawals (campaign_id, creator, amount_wei, tx_hash, block_number, log_index)
VALUES (?,?,?,?,?,?)`

	_, err := s.db.ExecContext(ctx, q, campaignID, creator, amountWei, txHash, blockNumber, logIndex)
	if err != nil {
		return fmt.Errorf("insert withdrawal: %w", err)
	}

	return nil
}

func (s *Store) GetCheckpoint(ctx context.Context, worker string, fallback uint64) (uint64, error) {
	const q = `SELECT last_scanned_block FROM indexer_checkpoints WHERE worker_name = ?`

	var block uint64
	err := s.db.QueryRowContext(ctx, q, worker).Scan(&block)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fallback, nil
		}
		return 0, fmt.Errorf("get checkpoint: %w", err)
	}

	return block, nil
}

type ContributionRecord struct {
	ID          uint64 `json:"id"`
	CampaignID  uint64 `json:"campaignId"`
	Funder      string `json:"funder"`
	AmountWei   string `json:"amountWei"`
	TxHash      string `json:"txHash"`
	BlockNumber uint64 `json:"blockNumber"`
	LogIndex    uint64 `json:"logIndex"`
	CreatedAt   string `json:"createdAt"`
}

func (s *Store) ListContributions(ctx context.Context, campaignID uint64, limit, offset int) ([]ContributionRecord, error) {
	const q = `
SELECT id, campaign_id, funder, amount_wei, tx_hash, block_number, log_index, created_at
FROM contributions
WHERE campaign_id = ?
ORDER BY id DESC
LIMIT ? OFFSET ?`

	rows, err := s.db.QueryContext(ctx, q, campaignID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list contributions query: %w", err)
	}
	defer rows.Close()

	var out []ContributionRecord
	for rows.Next() {
		var r ContributionRecord
		var createdAt any
		if err := rows.Scan(&r.ID, &r.CampaignID, &r.Funder, &r.AmountWei, &r.TxHash, &r.BlockNumber, &r.LogIndex, &createdAt); err != nil {
			return nil, fmt.Errorf("scan contribution: %w", err)
		}
		if t, ok := createdAt.(string); ok {
			r.CreatedAt = t
		} else {
			r.CreatedAt = fmt.Sprintf("%v", createdAt)
		}
		out = append(out, r)
	}

	return out, rows.Err()
}

func (s *Store) UpsertCheckpoint(ctx context.Context, worker string, block uint64) error {
	const q = `
INSERT INTO indexer_checkpoints (worker_name, last_scanned_block)
VALUES (?,?)
ON DUPLICATE KEY UPDATE last_scanned_block = VALUES(last_scanned_block), updated_at = CURRENT_TIMESTAMP`

	_, err := s.db.ExecContext(ctx, q, worker, block)
	if err != nil {
		return fmt.Errorf("upsert checkpoint: %w", err)
	}

	return nil
}
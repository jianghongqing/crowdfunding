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

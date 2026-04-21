package store

import (
	"context"
	"fmt"
)

func (s *PostgresStore) InsertContribution(
	ctx context.Context, campaignID uint64, funder, amountWei, txHash string, blockNumber uint64, logIndex uint,
) error {
	const q = `
INSERT INTO contributions (campaign_id, funder, amount_wei, tx_hash, block_number, log_index)
VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (tx_hash, log_index) DO NOTHING`
	_, err := s.pool.Exec(ctx, q, campaignID, funder, amountWei, txHash, blockNumber, logIndex)
	if err != nil {
		return fmt.Errorf("insert contribution: %w", err)
	}
	return nil
}

func (s *PostgresStore) InsertRefund(
	ctx context.Context, campaignID uint64, funder, amountWei, txHash string, blockNumber uint64, logIndex uint,
) error {
	const q = `
INSERT INTO refunds (campaign_id, funder, amount_wei, tx_hash, block_number, log_index)
VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (tx_hash, log_index) DO NOTHING`
	_, err := s.pool.Exec(ctx, q, campaignID, funder, amountWei, txHash, blockNumber, logIndex)
	if err != nil {
		return fmt.Errorf("insert refund: %w", err)
	}
	return nil
}

func (s *PostgresStore) InsertWithdrawal(
	ctx context.Context, campaignID uint64, creator, amountWei, txHash string, blockNumber uint64, logIndex uint,
) error {
	const q = `
INSERT INTO withdrawals (campaign_id, creator, amount_wei, tx_hash, block_number, log_index)
VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (tx_hash, log_index) DO NOTHING`
	_, err := s.pool.Exec(ctx, q, campaignID, creator, amountWei, txHash, blockNumber, logIndex)
	if err != nil {
		return fmt.Errorf("insert withdrawal: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetCheckpoint(ctx context.Context, worker string, fallback uint64) (uint64, error) {
	const q = `SELECT last_scanned_block FROM indexer_checkpoints WHERE worker_name = $1`
	var block uint64
	err := s.pool.QueryRow(ctx, q, worker).Scan(&block)
	if err != nil {
		return fallback, nil
	}
	return block, nil
}

func (s *PostgresStore) UpsertCheckpoint(ctx context.Context, worker string, block uint64) error {
	const q = `
INSERT INTO indexer_checkpoints (worker_name, last_scanned_block)
VALUES ($1,$2)
ON CONFLICT (worker_name) DO UPDATE SET last_scanned_block = EXCLUDED.last_scanned_block, updated_at = NOW()`
	_, err := s.pool.Exec(ctx, q, worker, block)
	if err != nil {
		return fmt.Errorf("upsert checkpoint: %w", err)
	}
	return nil
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// InsertContribution 写入捐款明细。
// INSERT IGNORE 保证幂等：indexer 重扫同一区块时不会插入重复记录（以 tx_hash + log_index 为唯一约束）。
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

// InsertRefund 写入退款明细，同样通过 INSERT IGNORE 保证幂等。
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

// InsertWithdrawal 写入提现明细。
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

// GetCheckpoint 读取 indexer 的扫块进度。
// 首次启动时 worker 不存在，返回 fallback（通常为合约部署区块号）。
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

// ContributionRecord 对应 contributions 表的一行。
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

// ListContributions 向后兼容的查询接口。
func (s *Store) ListContributions(ctx context.Context, campaignID uint64, limit, offset int) ([]ContributionRecord, error) {
	items, _, err := s.ListContributionsWithCount(ctx, campaignID, limit, offset)
	return items, err
}

// ListContributionsWithCount 分页查询某活动的捐款记录并返回总数。
func (s *Store) ListContributionsWithCount(ctx context.Context, campaignID uint64, limit, offset int) ([]ContributionRecord, int, error) {
	var total int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM contributions WHERE campaign_id = ?`, campaignID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count contributions: %w", err)
	}

	const q = `
SELECT id, campaign_id, funder, amount_wei, tx_hash, block_number, log_index, created_at
FROM contributions
WHERE campaign_id = ?
ORDER BY id DESC
LIMIT ? OFFSET ?`

	rows, err := s.db.QueryContext(ctx, q, campaignID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list contributions query: %w", err)
	}
	defer rows.Close()

	var out []ContributionRecord
	for rows.Next() {
		var r ContributionRecord
		// created_at 可能是 time.Time 或 string，取决于驱动配置，这里兼容两种
		var createdAt any
		if err := rows.Scan(&r.ID, &r.CampaignID, &r.Funder, &r.AmountWei, &r.TxHash, &r.BlockNumber, &r.LogIndex, &createdAt); err != nil {
			return nil, 0, fmt.Errorf("scan contribution: %w", err)
		}
		if t, ok := createdAt.(string); ok {
			r.CreatedAt = t
		} else {
			r.CreatedAt = fmt.Sprintf("%v", createdAt)
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return out, total, nil
}

// UpsertCheckpoint 更新 indexer 扫块进度，供下次启动时断点续扫。
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
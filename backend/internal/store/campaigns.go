// Package store 封装 MySQL 数据访问层。
// campaigns 表存储由 indexer 同步的链上活动快照，是 API 读取的主要数据源。
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CampaignRecord 对应 campaigns 表的一行，金额字段用字符串避免 uint256 溢出。
type CampaignRecord struct {
	CampaignID   uint64    `json:"campaignId"`
	Creator      string    `json:"creator"`
	Title        string    `json:"title"`
	GoalWei      string    `json:"goalWei"`
	PledgedWei   string    `json:"pledgedWei"`
	Deadline     uint64    `json:"deadline"`
	Withdrawn    bool      `json:"withdrawn"`
	Status       string    `json:"status"`
	CreatedBlock uint64    `json:"createdBlock"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// StatsResponse 平台级聚合统计，供 /api/v1/stats 端点返回。
type StatsResponse struct {
	TotalCampaigns int    `json:"totalCampaigns"`
	ActiveCount    int    `json:"activeCampaigns"`
	SuccessCount   int    `json:"succeededCampaigns"`
	FailedCount    int    `json:"failedCampaigns"`
	TotalPledged   string `json:"totalPledgedWei"`
}

// UpsertCampaign 插入或更新活动快照。
// 使用 ON DUPLICATE KEY UPDATE 实现幂等写入：indexer 重扫同一区块不会产生脏数据。
func (s *Store) UpsertCampaign(ctx context.Context, c CampaignRecord, txHash string) error {
	const q = `
INSERT INTO campaigns (campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, status, created_block, created_tx_hash)
VALUES (?,?,?,?,?,?,?,?,?,?)
ON DUPLICATE KEY UPDATE
    creator = VALUES(creator),
    title = VALUES(title),
    goal_wei = VALUES(goal_wei),
    pledged_wei = VALUES(pledged_wei),
    deadline = VALUES(deadline),
    withdrawn = VALUES(withdrawn),
    status = VALUES(status),
    updated_at = CURRENT_TIMESTAMP`

	_, err := s.db.ExecContext(
		ctx,
		q,
		c.CampaignID,
		c.Creator,
		c.Title,
		c.GoalWei,
		c.PledgedWei,
		c.Deadline,
		c.Withdrawn,
		c.Status,
		c.CreatedBlock,
		txHash,
	)
	if err != nil {
		return fmt.Errorf("upsert campaign: %w", err)
	}

	return nil
}

// ListCampaigns 向后兼容的列表查询（不返回总数），内部代理到 ListCampaignsWithCount。
func (s *Store) ListCampaigns(ctx context.Context, limit, offset int) ([]CampaignRecord, error) {
	items, _, err := s.ListCampaignsWithCount(ctx, limit, offset, "")
	return items, err
}

// ListCampaignsWithCount 分页查询活动列表并返回总数，支持按 status 筛选。
// 使用两次查询（COUNT + SELECT）而非 SQL_CALC_FOUND_ROWS，兼容性更好。
func (s *Store) ListCampaignsWithCount(ctx context.Context, limit, offset int, status string) ([]CampaignRecord, int, error) {
	var total int
	var countErr error

	if status != "" {
		countErr = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campaigns WHERE status = ?`, status).Scan(&total)
	} else {
		countErr = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campaigns`).Scan(&total)
	}
	if countErr != nil {
		return nil, 0, fmt.Errorf("count campaigns: %w", countErr)
	}

	var q string
	var args []any
	if status != "" {
		q = `
SELECT campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, status, created_block, updated_at
FROM campaigns
WHERE status = ?
ORDER BY campaign_id DESC
LIMIT ? OFFSET ?`
		args = []any{status, limit, offset}
	} else {
		q = `
SELECT campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, status, created_block, updated_at
FROM campaigns
ORDER BY campaign_id DESC
LIMIT ? OFFSET ?`
		args = []any{limit, offset}
	}

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list campaigns query: %w", err)
	}
	defer rows.Close()

	var out []CampaignRecord
	for rows.Next() {
		var r CampaignRecord
		if err := rows.Scan(
			&r.CampaignID,
			&r.Creator,
			&r.Title,
			&r.GoalWei,
			&r.PledgedWei,
			&r.Deadline,
			&r.Withdrawn,
			&r.Status,
			&r.CreatedBlock,
			&r.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan campaign: %w", err)
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return out, total, nil
}

func (s *Store) GetCampaign(ctx context.Context, campaignID uint64) (CampaignRecord, error) {
	const q = `
SELECT campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, status, created_block, updated_at
FROM campaigns
WHERE campaign_id = ?`

	var r CampaignRecord
	err := s.db.QueryRowContext(ctx, q, campaignID).Scan(
		&r.CampaignID,
		&r.Creator,
		&r.Title,
		&r.GoalWei,
		&r.PledgedWei,
		&r.Deadline,
		&r.Withdrawn,
		&r.Status,
		&r.CreatedBlock,
		&r.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CampaignRecord{}, sql.ErrNoRows
		}
		return CampaignRecord{}, fmt.Errorf("get campaign: %w", err)
	}

	return r, nil
}

// GetStats 执行多次 COUNT 聚合查询获取平台统计。
// 子状态查询失败不影响整体（忽略错误），只有总数查询失败才报错。
func (s *Store) GetStats(ctx context.Context) (StatsResponse, error) {
	var stats StatsResponse

	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campaigns`).Scan(&stats.TotalCampaigns)
	if err != nil {
		return stats, fmt.Errorf("count total: %w", err)
	}

	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campaigns WHERE status = 'active'`).Scan(&stats.ActiveCount)
	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campaigns WHERE status = 'succeeded_withdrawn'`).Scan(&stats.SuccessCount)
	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM campaigns WHERE status = 'failed_refundable'`).Scan(&stats.FailedCount)

	// pledged_wei 是字符串存储的大整数，CAST AS UNSIGNED 后求和
	var pledged sql.NullString
	_ = s.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(CAST(pledged_wei AS UNSIGNED)), 0) FROM campaigns`).Scan(&pledged)
	if pledged.Valid {
		stats.TotalPledged = pledged.String
	} else {
		stats.TotalPledged = "0"
	}

	return stats, nil
}
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

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

func (s *Store) ListCampaigns(ctx context.Context, limit, offset int) ([]CampaignRecord, error) {
	const q = `
SELECT campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, status, created_block, updated_at
FROM campaigns
ORDER BY campaign_id DESC
LIMIT ? OFFSET ?`

	rows, err := s.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list campaigns query: %w", err)
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
			return nil, fmt.Errorf("scan campaign: %w", err)
		}
		out = append(out, r)
	}

	return out, rows.Err()
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

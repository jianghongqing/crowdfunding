package store

import (
	"context"
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
	CreatedBlock uint64    `json:"createdBlock"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (s *PostgresStore) UpsertCampaign(ctx context.Context, c CampaignRecord, txHash string) error {
	const q = `
INSERT INTO campaigns (campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, created_block, created_tx_hash)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
ON CONFLICT (campaign_id) DO UPDATE
SET creator = EXCLUDED.creator,
    title = EXCLUDED.title,
    goal_wei = EXCLUDED.goal_wei,
    pledged_wei = EXCLUDED.pledged_wei,
    deadline = EXCLUDED.deadline,
    withdrawn = EXCLUDED.withdrawn,
    updated_at = NOW()`
	_, err := s.pool.Exec(ctx, q,
		c.CampaignID, c.Creator, c.Title, c.GoalWei, c.PledgedWei, c.Deadline, c.Withdrawn, c.CreatedBlock, txHash,
	)
	if err != nil {
		return fmt.Errorf("upsert campaign: %w", err)
	}
	return nil
}

func (s *PostgresStore) ListCampaigns(ctx context.Context, limit, offset int) ([]CampaignRecord, error) {
	const q = `
SELECT campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, created_block, updated_at
FROM campaigns
ORDER BY campaign_id DESC
LIMIT $1 OFFSET $2`
	rows, err := s.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list campaigns query: %w", err)
	}
	defer rows.Close()

	var out []CampaignRecord
	for rows.Next() {
		var r CampaignRecord
		if err := rows.Scan(
			&r.CampaignID, &r.Creator, &r.Title, &r.GoalWei, &r.PledgedWei, &r.Deadline, &r.Withdrawn, &r.CreatedBlock, &r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan campaign: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *PostgresStore) GetCampaign(ctx context.Context, campaignID uint64) (CampaignRecord, error) {
	const q = `
SELECT campaign_id, creator, title, goal_wei, pledged_wei, deadline, withdrawn, created_block, updated_at
FROM campaigns
WHERE campaign_id = $1`
	var r CampaignRecord
	err := s.pool.QueryRow(ctx, q, campaignID).Scan(
		&r.CampaignID, &r.Creator, &r.Title, &r.GoalWei, &r.PledgedWei, &r.Deadline, &r.Withdrawn, &r.CreatedBlock, &r.UpdatedAt,
	)
	if err != nil {
		return CampaignRecord{}, fmt.Errorf("get campaign: %w", err)
	}
	return r, nil
}

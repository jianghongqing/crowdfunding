ALTER TABLE campaigns
    ADD COLUMN IF NOT EXISTS status VARCHAR(32) NOT NULL DEFAULT 'active' AFTER withdrawn;

UPDATE campaigns
SET status = CASE
    WHEN withdrawn = TRUE THEN 'succeeded_withdrawn'
    WHEN CAST(pledged_wei AS DECIMAL(65,0)) >= CAST(goal_wei AS DECIMAL(65,0)) THEN 'goal_reached_pending_withdraw'
    WHEN deadline <= UNIX_TIMESTAMP() THEN 'failed_refundable'
    ELSE 'active'
END
WHERE status = '' OR status IS NULL OR status = 'active';

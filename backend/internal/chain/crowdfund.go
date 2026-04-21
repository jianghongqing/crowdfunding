package chain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"crowdfunding/backend/contracts/crowdfund"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type CampaignView struct {
	ID         uint64 `json:"id"`
	Creator    string `json:"creator"`
	Title      string `json:"title"`
	GoalWei    string `json:"goalWei"`
	PledgedWei string `json:"pledgedWei"`
	Deadline   uint64 `json:"deadline"`
	Withdrawn  bool   `json:"withdrawn"`
	Status     string `json:"status"`
}

type CrowdFundReader struct {
	contract *crowdfund.CrowdFundCaller
}

func NewCrowdFundReader(address common.Address, backend bind.ContractCaller) (*CrowdFundReader, error) {
	c, err := crowdfund.NewCrowdFundCaller(address, backend)
	if err != nil {
		return nil, fmt.Errorf("new crowdfund caller: %w", err)
	}
	return &CrowdFundReader{contract: c}, nil
}

func (r *CrowdFundReader) GetCampaign(ctx context.Context, campaignID uint64) (CampaignView, error) {
	raw, err := r.contract.GetCampaign(&bind.CallOpts{Context: ctx}, new(big.Int).SetUint64(campaignID))
	if err != nil {
		return CampaignView{}, fmt.Errorf("get campaign from chain: %w", err)
	}

	status := deriveStatus(raw.Pledged, raw.Goal, raw.Deadline, raw.Withdrawn)
	return CampaignView{
		ID:         raw.Id.Uint64(),
		Creator:    raw.Creator.Hex(),
		Title:      raw.Title,
		GoalWei:    raw.Goal.String(),
		PledgedWei: raw.Pledged.String(),
		Deadline:   raw.Deadline.Uint64(),
		Withdrawn:  raw.Withdrawn,
		Status:     status,
	}, nil
}

func (r *CrowdFundReader) GetContribution(ctx context.Context, campaignID uint64, user common.Address) (string, error) {
	value, err := r.contract.Contributions(&bind.CallOpts{Context: ctx}, new(big.Int).SetUint64(campaignID), user)
	if err != nil {
		return "", fmt.Errorf("get contribution from chain: %w", err)
	}
	return value.String(), nil
}

func deriveStatus(pledged, goal, deadline *big.Int, withdrawn bool) string {
	if withdrawn {
		return "succeeded_withdrawn"
	}
	if pledged != nil && goal != nil && pledged.Cmp(goal) >= 0 {
		return "goal_reached_pending_withdraw"
	}
	if deadline == nil {
		return "active"
	}
	if uint64(time.Now().Unix()) >= deadline.Uint64() {
		return "failed_refundable"
	}
	return "active"
}

// Package chain 封装与以太坊 RPC 的交互，提供合约数据读取能力。
// API 在数据库未命中时调用此包回退到链上查询。
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

// CampaignView 链上查询结果的 API 友好视图。
// 金额用字符串序列化，避免 JSON 大整数精度丢失。
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

// CrowdFundReader 只读合约调用器，通过 eth_call 读取链上状态（不消耗 gas）。
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

// GetCampaign 从链上读取活动完整状态并推导业务状态（active/failed/succeeded）。
func (r *CrowdFundReader) GetCampaign(ctx context.Context, campaignID uint64) (CampaignView, error) {
	raw, err := r.contract.GetCampaign(&bind.CallOpts{Context: ctx}, new(big.Int).SetUint64(campaignID))
	if err != nil {
		return CampaignView{}, fmt.Errorf("get campaign from chain: %w", err)
	}

	status := DeriveStatus(raw.Pledged, raw.Goal, raw.Deadline, raw.Withdrawn)
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

// GetContribution 读取某地址在某活动中的累计捐款额（对应合约 contributions mapping）。
func (r *CrowdFundReader) GetContribution(ctx context.Context, campaignID uint64, user common.Address) (string, error) {
	value, err := r.contract.Contributions(&bind.CallOpts{Context: ctx}, new(big.Int).SetUint64(campaignID), user)
	if err != nil {
		return "", fmt.Errorf("get contribution from chain: %w", err)
	}
	return value.String(), nil
}

// DeriveStatus 根据链上字段推导活动业务状态。
// 合约本身没有 status 字段，状态需要从 pledged/goal/deadline/withdrawn 四个维度计算。
// 判定优先级：已提款 > 已达标 > 已过期 > 进行中。
func DeriveStatus(pledged, goal, deadline *big.Int, withdrawn bool) string {
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
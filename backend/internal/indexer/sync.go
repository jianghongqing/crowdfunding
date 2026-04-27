package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"crowdfunding/backend/contracts/crowdfund"
	"crowdfunding/backend/internal/chain"
	"crowdfunding/backend/internal/config"
	"crowdfunding/backend/internal/store"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const workerName = "crowdfund-indexer"
const blockBatchSize uint64 = 1000

type Service struct {
	cfg      config.ChainConfig
	client   *ethclient.Client
	contract *crowdfund.CrowdFundFilterer
	caller   *crowdfund.CrowdFundCaller
	store    *store.Store

	topicCampaignCreated common.Hash
	topicFunded          common.Hash
	topicRefunded        common.Hash
	topicWithdrawn       common.Hash
}

func New(cfg config.ChainConfig, client *ethclient.Client, store *store.Store) (*Service, error) {
	address := common.HexToAddress(cfg.ContractAddress)
	filterer, err := crowdfund.NewCrowdFundFilterer(address, client)
	if err != nil {
		return nil, fmt.Errorf("new filterer: %w", err)
	}

	caller, err := crowdfund.NewCrowdFundCaller(address, client)
	if err != nil {
		return nil, fmt.Errorf("new caller: %w", err)
	}

	return &Service{
		cfg:                  cfg,
		client:               client,
		contract:             filterer,
		caller:               caller,
		store:                store,
		topicCampaignCreated: common.HexToHash("0xdc26653af5b99b2da33e2ad69ee6600d9aeccc82b034501db4338309615ca238"),
		topicFunded:          common.HexToHash("0x38c48552690c96ec2872092ac1db6c19fb59f5a8c5b49bbf41ed4886d0ca6926"),
		topicRefunded:        common.HexToHash("0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4"),
		topicWithdrawn:       common.HexToHash("0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372"),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	for {
		if err := s.syncOnce(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(8 * time.Second):
		}
	}
}

func (s *Service) syncOnce(ctx context.Context) error {
	head, err := s.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("read chain head: %w", err)
	}
	if head <= s.cfg.Confirmations {
		return nil
	}

	safeHead := head - s.cfg.Confirmations
	from, err := s.store.GetCheckpoint(ctx, workerName, s.cfg.DeploymentStartBlock)
	if err != nil {
		return err
	}
	if from > safeHead {
		return nil
	}

	for batchFrom := from; batchFrom <= safeHead; {
		batchTo := min(batchFrom+blockBatchSize-1, safeHead)
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(0).SetUint64(batchFrom),
			ToBlock:   big.NewInt(0).SetUint64(batchTo),
			Addresses: []common.Address{common.HexToAddress(s.cfg.ContractAddress)},
		}

		logs, err := s.client.FilterLogs(ctx, query)
		if err != nil {
			return fmt.Errorf("filter logs %d-%d: %w", batchFrom, batchTo, err)
		}

		for _, lg := range logs {
			if err := s.handleLog(ctx, lg); err != nil {
				return err
			}
		}

		if err := s.store.UpsertCheckpoint(ctx, workerName, batchTo+1); err != nil {
			return err
		}

		batchFrom = batchTo + 1
	}

	return nil
}

func (s *Service) handleLog(ctx context.Context, lg types.Log) error {
	if len(lg.Topics) == 0 {
		return nil
	}

	switch lg.Topics[0] {
	case s.topicCampaignCreated:
		ev, err := s.contract.ParseCampaignCreated(lg)
		if err != nil {
			return fmt.Errorf("parse CampaignCreated: %w", err)
		}

		onchain, err := s.caller.GetCampaign(&bind.CallOpts{Context: ctx}, ev.CampaignId)
		if err != nil {
			return fmt.Errorf("read campaign after create: %w", err)
		}

		return s.store.UpsertCampaign(ctx, store.CampaignRecord{
			CampaignID:   ev.CampaignId.Uint64(),
			Creator:      ev.Creator.Hex(),
			Title:        onchain.Title,
			GoalWei:      onchain.Goal.String(),
			PledgedWei:   onchain.Pledged.String(),
			Deadline:     onchain.Deadline.Uint64(),
			Withdrawn:    onchain.Withdrawn,
			Status:       chain.DeriveStatus(onchain.Pledged, onchain.Goal, onchain.Deadline, onchain.Withdrawn),
			CreatedBlock: lg.BlockNumber,
		}, lg.TxHash.Hex())
	case s.topicFunded:
		ev, err := s.contract.ParseFunded(lg)
		if err != nil {
			return fmt.Errorf("parse Funded: %w", err)
		}

		if err := s.store.InsertContribution(
			ctx, ev.CampaignId.Uint64(), ev.Funder.Hex(), ev.Amount.String(), lg.TxHash.Hex(), lg.BlockNumber, lg.Index,
		); err != nil {
			return err
		}

		return s.refreshCampaign(ctx, ev.CampaignId.Uint64())
	case s.topicRefunded:
		ev, err := s.contract.ParseRefunded(lg)
		if err != nil {
			return fmt.Errorf("parse Refunded: %w", err)
		}

		if err := s.store.InsertRefund(
			ctx, ev.CampaignId.Uint64(), ev.Funder.Hex(), ev.Amount.String(), lg.TxHash.Hex(), lg.BlockNumber, lg.Index,
		); err != nil {
			return err
		}

		return s.refreshCampaign(ctx, ev.CampaignId.Uint64())
	case s.topicWithdrawn:
		ev, err := s.contract.ParseWithdrawn(lg)
		if err != nil {
			return fmt.Errorf("parse Withdrawn: %w", err)
		}

		if err := s.store.InsertWithdrawal(
			ctx, ev.CampaignId.Uint64(), ev.Creator.Hex(), ev.Amount.String(), lg.TxHash.Hex(), lg.BlockNumber, lg.Index,
		); err != nil {
			return err
		}

		return s.refreshCampaign(ctx, ev.CampaignId.Uint64())
	default:
		return nil
	}
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func (s *Service) refreshCampaign(ctx context.Context, campaignID uint64) error {
	raw, err := s.caller.GetCampaign(&bind.CallOpts{Context: ctx}, big.NewInt(0).SetUint64(campaignID))
	if err != nil {
		return fmt.Errorf("refresh campaign from chain: %w", err)
	}

	return s.store.UpsertCampaign(ctx, store.CampaignRecord{
		CampaignID:   raw.Id.Uint64(),
		Creator:      raw.Creator.Hex(),
		Title:        raw.Title,
		GoalWei:      raw.Goal.String(),
		PledgedWei:   raw.Pledged.String(),
		Deadline:     raw.Deadline.Uint64(),
		Withdrawn:    raw.Withdrawn,
		Status:       chain.DeriveStatus(raw.Pledged, raw.Goal, raw.Deadline, raw.Withdrawn),
		CreatedBlock: 0,
	}, "")
}

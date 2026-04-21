package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

type ChainConfig struct {
	ChainName            string `json:"chainName"`
	ChainID              uint64 `json:"chainId"`
	RPCHTTPURL           string `json:"rpcHttpUrl"`
	RPCWSURL             string `json:"rpcWsUrl"`
	ContractAddress      string `json:"contractAddress"`
	DeploymentStartBlock uint64 `json:"deploymentStartBlock"`
	Confirmations        uint64 `json:"confirmations"`
}

func Load(path string) (ChainConfig, error) {
	var cfg ChainConfig

	raw, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	if cfg.RPCHTTPURL == "" {
		return cfg, fmt.Errorf("rpcHttpUrl is required")
	}
	if !common.IsHexAddress(cfg.ContractAddress) {
		return cfg, fmt.Errorf("contractAddress is invalid")
	}
	if cfg.Confirmations == 0 {
		cfg.Confirmations = 5
	}

	return cfg, nil
}

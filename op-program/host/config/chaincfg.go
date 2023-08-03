package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var enabledFromBedrockBlock = uint64(0)

var OPGoerliChainConfig, OPSepoliaChainConfig, OPMainnetChainConfig *params.ChainConfig

func init() {
	mustLoadConfig := func(chainID uint64) *params.ChainConfig {
		cfg, err := params.LoadOPStackChainConfig(chainID)
		if err != nil {
			panic(err)
		}
		return cfg
	}
	OPGoerliChainConfig = mustLoadConfig(420)
	OPSepoliaChainConfig = mustLoadConfig(11155420)
	OPSepoliaChainConfig = mustLoadConfig(10)
}

var L2ChainConfigsByName = map[string]*params.ChainConfig{
	"goerli":  OPGoerliChainConfig,
	"sepolia": OPSepoliaChainConfig,
	"mainnet": OPMainnetChainConfig,
}

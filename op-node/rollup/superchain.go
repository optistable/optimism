package rollup

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/superchain-registry/superchain"
)

const (
	opMainnet  = 10
	opGoerli   = 420
	opSepolia  = 11155420
	baseGoerli = 84531
	baseMainnet = 8453
	pgnMainnet = 424
	pgnSepolia = 58008
	zoraGoerli = 999
	zoraMainnet = 7777777
)

var genesisSysConfigs = map[uint64]eth.SystemConfig{
	opMainnet: {
		BatcherAddr: common.HexToAddress("0x6887246668a3b87f54deb3b94ba47a6f63f32985"),
		Overhead:    eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0")),
		GasLimit:    30_000_000,
	},
	opGoerli: {
		BatcherAddr: common.HexToAddress("0x7431310e026B69BFC676C0013E12A1A11411EEc9"),
		Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
		GasLimit:    25_000_000,
	},
	opSepolia: {
		BatcherAddr: common.HexToAddress("0x7431310e026b69bfc676c0013e12a1a11411eec9"),
		Overhead:    eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0")),
		GasLimit:    30_000_000,
	},
	baseGoerli: {
		BatcherAddr: common.HexToAddress("0x2d679b567db6187c0c8323fa982cfb88b74dbcc7"),
		Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
		GasLimit:    25_000_000,
	},
	baseMainnet: {
		BatcherAddr: common.HexToAddress("0x5050f69a9786f081509234f1a7f4684b5e5b76c9"),
		Overhead:    eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0")),
		GasLimit:    30_000_000,
	},
	pgnSepolia: {
		BatcherAddr: common.HexToAddress("0x7224e05E6cF6E07aFBE1eFa09a3fA23A637DD485"),
		Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
		GasLimit:    30_000_000,
	},
	pgnMainnet: {
		BatcherAddr: common.HexToAddress("0x99526b0e49A95833E734EB556A6aBaFFAb0Ee167"),
		Overhead:    eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0")),
		GasLimit:    30_000_000,
	},
	zoraGoerli: {
		BatcherAddr: common.HexToAddress("0x427c9a666d3b27873111cE3894712Bf64C6343A0"),
		Overhead:    eth.Bytes32(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000834")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000f4240")),
		GasLimit:    30_000_000,
	},
	zoraMainnet: {
		BatcherAddr: common.HexToAddress("0x625726c858dBF78c0125436C943Bf4b4bE9d9033"),
		Overhead:    eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bc")),
		Scalar:      eth.Bytes32(common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000a6fe0")),
		GasLimit:    30_000_000,
	},
}

var depositContractAddrs = map[uint64]common.Address{
	opMainnet: common.HexToAddress("0xbEb5Fc579115071764c7423A4f12eDde41f106Ed"),
	opGoerli:  common.HexToAddress("0x5b47E1A08Ea6d985D6649300584e6722Ec4B1383"),
	opSepolia: common.HexToAddress("0x8f6452d842438c4e22ba18baa21652ff65530df4"),
	baseGoerli: common.HexToAddress("0xe93c8cd0d409341205a592f8c4ac1a5fe5585cfa"),
	baseMainnet: common.HexToAddress("0x49048044d57e1c92a77f79988d21fa8faf74e97e"),
	pgnSepolia: common.HexToAddress("0xF04BdD5353Bb0EFF6CA60CfcC78594278eBfE179"),
	pgnMainnet: common.HexToAddress("0xb26Fd985c5959bBB382BAFdD0b879E149e48116c"),
	zoraGoerli: common.HexToAddress("0xDb9F51790365e7dc196e7D072728df39Be958ACe"),
	zoraMainnet: common.HexToAddress("0x1a0ad011913A150f69f6A19DF447A0CfD9551054"),
}


// SystemConfigProvider is an interface to retrieve SystemConfig variables, as registered in the onchain smart contract.
type SystemConfigProvider interface {
	GenesisSystemConfig() eth.SystemConfig
	DepositContractAddress() common.Address
}

// LoadOPStackRollupConfig loads the rollup configuration of the requested chain ID from the superchain-registry.
// Some chains may require a SystemConfigProvider to retrieve any values not part of the registry.
func LoadOPStackRollupConfig(sysConfig SystemConfigProvider, chainID uint64) (*Config, error) {
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}

	superChain, ok := superchain.Superchains[chConfig.Superchain]
	if !ok {
		return nil, fmt.Errorf("unknown superchain: %q", chConfig.Superchain)
	}

	var genesisSysConfig eth.SystemConfig
	if sysCfg, ok := genesisSysConfigs[chainID]; ok {
		genesisSysConfig = sysCfg
	} else if sysConfig != nil {
		genesisSysConfig = sysConfig.GenesisSystemConfig()
	} else {
		return nil, fmt.Errorf("unable to retrieve genesis SystemConfig")
	}

	var depositContractAddress common.Address
	if addr, ok := depositContractAddrs[chainID]; ok {
		depositContractAddress = addr
	} else if sysConfig != nil {
		depositContractAddress = sysConfig.DepositContractAddress()
	} else {
		return nil, fmt.Errorf("unable to retrieve deposit contract address")
	}

	regolithTime := uint64(0)
	// two goerli testnets test-ran Bedrock and later upgraded to Regolith.
	// All other OP-Stack chains have Regolith enabled from the start.
	switch chainID {
	case baseGoerli:
		regolithTime = 1683219600
	case opGoerli:
		regolithTime = 1679079600
	}

	cfg := &Config{
		Genesis: Genesis{
			L1: eth.BlockID{
				Hash:   common.Hash(chConfig.Genesis.L1.Hash),
				Number: chConfig.Genesis.L1.Number,
			},
			L2: eth.BlockID{
				Hash:   common.Hash(chConfig.Genesis.L2.Hash),
				Number: chConfig.Genesis.L2.Number,
			},
			L2Time:       chConfig.Genesis.L2Time,
			SystemConfig: genesisSysConfig,
		},
		// The below chain parameters can be different per OP-Stack chain,
		// but since none of the superchain chains differ, it's not represented in the superchain-registry yet.
		// This restriction on superchain-chains may change in the future.
		// Test/Alt configurations can still load custom rollup-configs when necessary.
		BlockTime:              2,
		MaxSequencerDrift:      600,
		SeqWindowSize:          3600,
		ChannelTimeout:         300,
		L1ChainID:              new(big.Int).SetUint64(superChain.Config.L1.ChainID),
		L2ChainID:              new(big.Int).SetUint64(chConfig.ChainID),
		RegolithTime:           &regolithTime,
		BatchInboxAddress:      common.Address(chConfig.BatchInboxAddr),
		DepositContractAddress: depositContractAddress,
		L1SystemConfigAddress:  common.Address(chConfig.SystemConfigAddr),
	}
	return cfg, nil
}

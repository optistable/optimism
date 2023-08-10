package derive

// this file defines a new transaction type
// and how it should be encoded and decoded

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-service/solabi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	aggregatorv3 "github.com/ethereum-optimism/optimism/op-node/aggregatorv3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ChainlinkReportFuncSignature = "recordPrice(uint256,uint256)"
	ChainlinkReportArguments     = 2
	ChainlinkReportLen           = 4 + 32*ChainlinkReportArguments
)

var (
	ChainlinkReportFuncBytes4 = crypto.Keccak256([]byte(ChainlinkReportFuncSignature))[:4]
	ChainlinkReportAddress    = common.HexToAddress("0x0ba4D449658758Dae8a8793e8182acc7b5f8976f")
)

type ChainlinkInfo struct {
	Number *big.Int
	Price  *big.Int
}

func (info *ChainlinkInfo) MarshalBinary() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, ChainlinkReportLen))
	if err := solabi.WriteSignature(w, ChainlinkReportFuncBytes4); err != nil {
		return nil, err
	}
	if err := solabi.WriteUint256(w, info.Number); err != nil {
		return nil, err
	}
	if err := solabi.WriteUint256(w, info.Price); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (info *ChainlinkInfo) UnmarshalBinary(data []byte) error {
	if len(data) != ChainlinkReportLen {
		return fmt.Errorf("data is unexpected length: %d", len(data))
	}
	reader := bytes.NewReader(data)

	var err error
	if _, err := solabi.ReadAndValidateSignature(reader, ChainlinkReportFuncBytes4); err != nil {
		return err
	}
	if info.Number, err = solabi.ReadUint256(reader); err != nil {
		return err
	}
	if info.Price, err = solabi.ReadUint256(reader); err != nil {
		return err
	}

	return nil
}

func ChainlinkInfoDepositTxData(data []byte) (ChainlinkInfo, error) {
	var info ChainlinkInfo
	err := info.UnmarshalBinary(data)
	return info, err
}

func makeChainlinkCall() *big.Int {
	// Fetch the rpc_url.
	rpcUrl := "https://eth-sepolia.g.alchemy.com/v2/MdbUH8ez_zjMYPZBkIhj8FQhBhlIw1wt"

	// Assign default values to feedAddress, and update value if a feed address was passed in the command line.
	feedAddress := "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43"

	// Initialize client instance using the rpcUrl.
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		fmt.Println(err)
	}

	chainlinkPriceFeedProxyAddress := common.HexToAddress(feedAddress)
	chainlinkPriceFeedProxy, err := aggregatorv3.NewAggregatorV3Interface(chainlinkPriceFeedProxyAddress, client)
	if err != nil {
		fmt.Println(err)
	}

	roundData, err := chainlinkPriceFeedProxy.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}

	decimals, err := chainlinkPriceFeedProxy.Decimals(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}

	description, err := chainlinkPriceFeedProxy.Description(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
	}

	// Compute a big.int which is 10**decimals.
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	fmt.Printf("%v Price feed address is  %v\n", description, chainlinkPriceFeedProxyAddress)
	fmt.Printf("Round id is %v\n", roundData.RoundId)
	fmt.Printf("Answer is %v\n", roundData.Answer)
	fmt.Printf("Formatted answer is %v\n", divideBigInt(roundData.Answer, divisor))
	fmt.Printf("Started at %v\n", roundData.StartedAt)
	fmt.Printf("Updated at %v\n", roundData.UpdatedAt)
	fmt.Printf("Answered in round %v\n", roundData.AnsweredInRound)

	return roundData.Answer
	// return roundData.Answer.Uint64()
}

func divideBigInt(num1 *big.Int, num2 *big.Int) *big.Float {
	if num2.BitLen() == 0 {
		panic("cannot divide by zero.")
	}
	num1BigFloat := new(big.Float).SetInt(num1)
	num2BigFloat := new(big.Float).SetInt(num2)
	result := new(big.Float).Quo(num1BigFloat, num2BigFloat)
	return result
}

func ChainlinkInfoDeposit(seqNumber uint64, block eth.BlockInfo, sysCfg eth.SystemConfig) (*types.DepositTx, error) {
	// L1 info:
	// TODO
	// record the (L2) block where we checked for a price update
	// record the L1 block when the price was last updated
	// record the price, source and network (if on chain)

	chainLinkPriceU256 := makeChainlinkCall()
	infoDat := ChainlinkInfo{
		Number: big.NewInt(0).SetUint64(block.NumberU64()),
		Price:  chainLinkPriceU256,
	}

	data, err := infoDat.MarshalBinary()
	if err != nil {
		return nil, err
	}

	source := L1InfoDepositSource{
		L1BlockHash: block.Hash(),
		// qq: what is seqnumber used for
		SeqNumber: seqNumber,
	}

	fmt.Println("==== send ChainlinkReport info ", L1InfoDepositerAddress, ChainlinkReportAddress, infoDat.Number, infoDat.Price)

	return &types.DepositTx{
		SourceHash:          source.SourceHash(),
		From:                L1InfoDepositerAddress,
		To:                  &ChainlinkReportAddress,
		Mint:                nil,
		Value:               big.NewInt(0),
		Gas:                 150_000_000,
		IsSystemTransaction: false,
		Data:                data,
	}, nil
}

func ChainlinkInfoDepositBytes(seqNumber uint64, l1Info eth.BlockInfo, sysCfg eth.SystemConfig) ([]byte, error) {
	dep, err := ChainlinkInfoDeposit(seqNumber, l1Info, sysCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create l1 burn tx: %w", err)
	}

	l1Tx := types.NewTx(dep)
	opaqueL1Tx, err := l1Tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to encode l1 burn tx: %w", err)
	}
	return opaqueL1Tx, nil

}

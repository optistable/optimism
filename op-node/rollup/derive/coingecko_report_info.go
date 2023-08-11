package derive

// this file defines a new transaction type
// and how it should be encoded and decoded

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-service/solabi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	CoingeckoReportFuncSignature = "recordPrice(uint256,uint256)"
	CoingeckoReportArguments     = 2
	CoingeckoReportLen           = 4 + 32*ChainlinkReportArguments
)

var (
	CoingeckoReportFuncBytes4 = crypto.Keccak256([]byte(ChainlinkReportFuncSignature))[:4]
	CoingeckoReportAddress    = common.HexToAddress("0xd197a45De818f46781e59267Fb86F026D34F884d")
)

type PriceResponse struct {
	USDCoin struct {
		USD float64 `json:"usd"`
	} `json:"usd-coin"`
}

type CoingeckoInfo struct {
	Number *big.Int
	Price  *big.Int
}

func (info *CoingeckoInfo) MarshalBinary() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, CoingeckoReportLen))
	if err := solabi.WriteSignature(w, CoingeckoReportFuncBytes4); err != nil {
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

func (info *CoingeckoInfo) UnmarshalBinary(data []byte) error {
	if len(data) != CoingeckoReportLen {
		return fmt.Errorf("data is unexpected length: %d", len(data))
	}
	reader := bytes.NewReader(data)

	var err error
	if _, err := solabi.ReadAndValidateSignature(reader, CoingeckoReportFuncBytes4); err != nil {
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

func CoingeckoInfoDepositTxData(data []byte) (CoingeckoInfo, error) {
	var info CoingeckoInfo
	err := info.UnmarshalBinary(data)
	return info, err
}

var lastCoingeckoPrice float64

func makeCoingeckoCall() *big.Int {

	url := "https://api.coingecko.com/api/v3/simple/price?ids=usd-coin&vs_currencies=usd"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	var price PriceResponse
	err = json.Unmarshal(body, &price)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	fmt.Println("original coingecko price", price.USDCoin.USD)
	if price.USDCoin.USD == 0 {
		price.USDCoin.USD = lastCoingeckoPrice
	} else {
		lastCoingeckoPrice = price.USDCoin.USD
	}

	power := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	bf := new(big.Float).SetFloat64(price.USDCoin.USD)

	bf.Mul(bf, power)
	bi := new(big.Int)
	bf.Int(bi)

	return bi
}

func CoingeckoInfoDeposit(seqNumber uint64, block eth.BlockInfo, sysCfg eth.SystemConfig, regolith bool) (*types.DepositTx, error) {
	coingeckoPriceU256 := makeCoingeckoCall()
	infoDat := CoingeckoInfo{
		Number: big.NewInt(0).SetUint64(block.NumberU64()),
		Price:  coingeckoPriceU256,
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

	fmt.Println("==== send CoingeckoReport info ", L1InfoDepositerAddress, CoingeckoReportAddress, infoDat.Number, infoDat.Price)

	out := &types.DepositTx{
		SourceHash:          source.SourceHash(),
		From:                L1InfoDepositerAddress,
		To:                  &CoingeckoReportAddress,
		Mint:                nil,
		Value:               big.NewInt(0),
		Gas:                 150_000_000,
		IsSystemTransaction: true,
		Data:                data,
	}

	if regolith {
		out.IsSystemTransaction = false
		out.Gas = RegolithSystemTxGas
	}

	return out, nil
}

func CoingeckoInfoDepositBytes(seqNumber uint64, l1Info eth.BlockInfo, sysCfg eth.SystemConfig, regolith bool) ([]byte, error) {
	dep, err := CoingeckoInfoDeposit(seqNumber, l1Info, sysCfg, regolith)
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

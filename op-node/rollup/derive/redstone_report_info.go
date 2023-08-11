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
	RedstoneReportFuncSignature = "recordPrice(uint256,uint256)"
	RedstoneReportArguments     = 2
	RedstoneReportLen           = 4 + 32*RedstoneReportArguments
)

var (
	RedstoneReportFuncBytes4 = crypto.Keccak256([]byte(RedstoneReportFuncSignature))[:4]
	RedstoneReportAddress    = common.HexToAddress("0x5FdEd0D534D0D880760394fdF83A45aCFAD3ca99")
)

// Structure to match the JSON response from Redstone API
type RedstonePriceResponse struct {
	Symbol string  `json:"symbol"`
	Value  float64 `json:"value"`
}

type RedstoneInfo struct {
	Number *big.Int
	Price  *big.Int
}

func (info *RedstoneInfo) MarshalBinary() ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, RedstoneReportLen))
	if err := solabi.WriteSignature(w, RedstoneReportFuncBytes4); err != nil {
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

func (info *RedstoneInfo) UnmarshalBinary(data []byte) error {
	if len(data) != RedstoneReportLen {
		return fmt.Errorf("data is unexpected length: %d", len(data))
	}
	reader := bytes.NewReader(data)

	var err error
	if _, err := solabi.ReadAndValidateSignature(reader, RedstoneReportFuncBytes4); err != nil {
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

func RedstoneInfoDepositTxData(data []byte) (RedstoneInfo, error) {
	var info RedstoneInfo
	err := info.UnmarshalBinary(data)
	return info, err
}

var lastRedstonePrice float64

func makeRedstoneCall() *big.Int {
	url := "https://api.redstone.finance/prices/?symbol=USDC&provider=redstone&limit=1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	var prices []RedstonePriceResponse
	err = json.Unmarshal(body, &prices)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	fmt.Println("original redstone price", prices)
	if len(prices) == 0 || prices[0].Value == 0 {
		prices[0].Value = lastRedstonePrice
	} else {
		lastRedstonePrice = prices[0].Value
	}

	power := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	bf := new(big.Float).SetFloat64(prices[0].Value)

	bf.Mul(bf, power)
	bi := new(big.Int)
	bf.Int(bi)

	return bi
}

func RedstoneInfoDeposit(seqNumber uint64, block eth.BlockInfo, sysCfg eth.SystemConfig, regolith bool) (*types.DepositTx, error) {
	redstonePriceU256 := makeRedstoneCall()
	infoDat := RedstoneInfo{
		Number: big.NewInt(0).SetUint64(block.NumberU64()),
		Price:  redstonePriceU256,
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

	fmt.Println("==== send RedstoneReport info ", L1InfoDepositerAddress, RedstoneReportAddress, infoDat.Number, infoDat.Price)

	out := &types.DepositTx{
		SourceHash:          source.SourceHash(),
		From:                L1InfoDepositerAddress,
		To:                  &RedstoneReportAddress,
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

func RedstoneInfoDepositBytes(seqNumber uint64, l1Info eth.BlockInfo, sysCfg eth.SystemConfig, regolith bool) ([]byte, error) {
	dep, err := RedstoneInfoDeposit(seqNumber, l1Info, sysCfg, regolith)
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

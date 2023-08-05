package derive

// this file defines a new transaction type
// and how it should be encoded and decoded

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	L1BurnFuncSignature = "report(uint64,uint64)"
	L1BurnArguments     = 2
	L1BurnLen           = 4 + 32*L1BurnArguments
)

var (
	L1BurnFuncBytes4 = crypto.Keccak256([]byte(L1BurnFuncSignature))[:4]
	L1BurnAddress    = common.HexToAddress("0x4081101F39205EdD2eE7aA2756D01bb2fFBe56e6")
)

type L1BurnInfo struct {
	Number uint64
	Burn   uint64
}

func (info *L1BurnInfo) MarshalBinary() ([]byte, error) {
	data := make([]byte, L1BurnLen)
	offset := 0
	copy(data[offset:4], L1BurnFuncBytes4)
	offset += 4
	// qq: why do we leave first 24 bytes empty?
	binary.BigEndian.PutUint64(data[offset+24:offset+32], info.Number)
	offset += 32
	// qq: why do we leave first 24 bytes empty? padding??
	binary.BigEndian.PutUint64(data[offset+24:offset+32], info.Burn)
	return data, nil
}

func (info *L1BurnInfo) UnmarshalBinary(data []byte) error {
	if len(data) != L1BurnLen {
		return fmt.Errorf("data is unexpected length: %d", len(data))
	}
	var padding [24]byte
	offset := 4
	info.Number = binary.BigEndian.Uint64(data[offset+24 : offset+32])
	if !bytes.Equal(data[offset:offset+24], padding[:]) {
		return fmt.Errorf("l1 burn tx number exceeds uint64 bounds: %x", data[offset:offset+32])
	}
	offset += 32
	info.Burn = binary.BigEndian.Uint64(data[offset+24 : offset+32])
	if !bytes.Equal(data[offset:offset+24], padding[:]) {
		return fmt.Errorf("l1 burn tx number exceeds uint64 bounds: %x", data[offset:offset+32])
	}

	return nil
}

func L1BurnDepositTxData(data []byte) (L1BurnInfo, error) {
	var info L1BurnInfo
	err := info.UnmarshalBinary(data)
	return info, err
}

func L1BurnDeposit(seqNumber uint64, block eth.BlockInfo, sysCfg eth.SystemConfig) (*types.DepositTx, error) {
	infoDat := L1BurnInfo{
		Number: block.NumberU64(),
		Burn:   block.BaseFee().Uint64() * block.GasUsed(),
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

	fmt.Println("==== send L1Burn info ", L1BurnAddress, infoDat.Burn)

	return &types.DepositTx{
		SourceHash:          source.SourceHash(),
		From:                L1InfoDepositerAddress,
		To:                  &L1BurnAddress,
		Mint:                nil,
		Value:               big.NewInt(0),
		Gas:                 150_000_000,
		IsSystemTransaction: true,
		Data:                data,
	}, nil
}

func L1BurnDepositBytes(seqNumber uint64, l1Info eth.BlockInfo, sysCfg eth.SystemConfig) ([]byte, error) {
	dep, err := L1BurnDeposit(seqNumber, l1Info, sysCfg)
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

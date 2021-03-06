// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package syscontract

import (
	"errors"
	"fmt"
	"github.com/AsimovNetwork/asimov/asiutil"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/protos"
	"github.com/AsimovNetwork/asimov/vm/fvm"
	"github.com/AsimovNetwork/asimov/vm/fvm/core/vm"
	"github.com/AsimovNetwork/asimov/vm/fvm/params"
	"math/big"
)

// FeeItem defines assets and their effective heights
type FeeItem struct {
	Assets *protos.Assets
	Height int32
}

// GetFees returns an asset list and their valid heights
// assets in list are used as transaction fees
// heights describe the assets formally effective
// the asset list is submitted by members of validator committee
func (m *Manager) GetFees(
	block *asiutil.Block,
	stateDB vm.StateDB,
	chainConfig *params.ChainConfig) (map[protos.Assets]int32, error, uint64) {

	officialAddr := chaincfg.OfficialAddress
	contract := m.GetActiveContractByHeight(block.Height(), common.ValidatorCommittee)
	if contract == nil {
		errStr := fmt.Sprintf("Failed to get active contract %s, %d", common.ValidatorCommittee, block.Height())
		log.Error(errStr)
		panic(errStr)
	}
	proxyAddr, abi := vm.ConvertSystemContractAddress(common.ValidatorCommittee), contract.AbiInfo

	feeListFunc := common.ContractValidatorCommittee_GetAssetFeeListFunction()

	runCode, err := fvm.PackFunctionArgs(abi, feeListFunc)
	result, leftOverGas, err := fvm.CallReadOnlyFunction(officialAddr, block, m.chain, stateDB, chainConfig,
		common.SystemContractReadOnlyGas, proxyAddr, runCode)
	if err != nil {
		log.Errorf("Get contract templates failed, error: %s", err)
		return nil, err, leftOverGas
	}
	assets := make([]*big.Int, 0)
	height := make([]*big.Int, 0)
	outData := []interface{}{
		&assets,
		&height,
	}
	err = fvm.UnPackFunctionResult(abi, &outData, feeListFunc, result)
	if err != nil {
		log.Errorf("Get fee list failed, error: %s", err)
		return nil, err, leftOverGas
	}
	if len(assets) != len(height) {
		errStr := "get fee list failed, length of assets does not match length of height"
		log.Errorf(errStr)
		return nil, errors.New(errStr), leftOverGas
	}

	fees := make(map[protos.Assets]int32)
	for i := 0; i < len(assets); i++ {
		pAssets := protos.AssetFromBytes(assets[i].Bytes())
		fees[*pAssets] = int32(height[i].Int64())
	}

	return fees, nil, leftOverGas
}

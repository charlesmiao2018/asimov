// Copyright (c) 2018-2020. The asimov developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package syscontract

import (
	"encoding/json"
	"github.com/AsimovNetwork/asimov/ainterface"
	"github.com/AsimovNetwork/asimov/chaincfg"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/vm/fvm"
)

// Manager defines an contract manager that manages multiple system contracts and
// implements the blockchain.ContractManager interface so it can be seamlessly
// plugged into normal chain processing.
type Manager struct {
	chain fvm.ChainContext
	//  genesis transaction data cache
	genesisDataCache map[common.ContractCode][]chaincfg.ContractInfo
}

// Init manager by genesis data.
func (m *Manager) Init(chain fvm.ChainContext, dataBytes [] byte) error {
	var cMap map[common.ContractCode][]chaincfg.ContractInfo
	err := json.Unmarshal(dataBytes, &cMap)
	if err != nil {
		return err
	}
	m.chain = chain
	m.genesisDataCache = cMap
	return nil
}

// Get latest contract by height.
func (m *Manager) GetActiveContractByHeight(height int32, delegateAddr common.ContractCode) *chaincfg.ContractInfo {
	contracts, ok := m.genesisDataCache[delegateAddr]
	if !ok {
		return nil
	}
	for i := len(contracts) - 1; i >= 0; i-- {
		if height >= contracts[i].BlockHeight {
			return &contracts[i]
		}
	}
	return nil
}

// NewContractManager returns an empty struct of Manager
func NewContractManager() ainterface.ContractManager {
    return &Manager {}
}

// Ensure the Manager type implements the blockchain.ContractManager interface.
var _ ainterface.ContractManager = (*Manager)(nil)
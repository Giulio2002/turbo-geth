package commands

import (
	"context"
	"fmt"

	"github.com/ledgerwatch/turbo-geth/common"
)

// Coinbase is the address that mining rewards will be sent to
func (api *APIImpl) Coinbase(_ context.Context) (common.Address, error) {
	if api.ethBackend == nil {
		// We're running in --chaindata mode or otherwise cannot get the backend
		return common.Address{}, fmt.Errorf(NotAvailableChainData, "eth_coinbase")
	}
	return api.ethBackend.Etherbase()
}

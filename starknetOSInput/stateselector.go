package osinput

import (
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
)

// StateSelector contains the set of contract addresses and class hashes that a transaction
// touches during execution / state updates.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/storage/storage.py#L403
// Note: each txn type used to have it's own get_state_selecto() method, but now that info is in
// the execution info https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/business_logic/execution/deprecated_objects.py#L382
type StateSelector struct {
	contractAddresses []felt.Felt
	classHashes       []felt.Felt
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/objects.py#L1421
func get_state_selector_invoke(txn *core.InvokeTransaction, config *StarknetGeneralConfig) (*StateSelector, error) {
	contractAddresses := []felt.Felt{}
	contractAddresses = append(contractAddresses, *txn.SenderAddress)

	if !txn.MaxFee.IsZero() {
		contractAddresses = append(contractAddresses, config.StarknetOsConfig.FeeTokenAddress)
	}

	return &StateSelector{
		contractAddresses: contractAddresses,
		classHashes:       nil,
	}, nil
}

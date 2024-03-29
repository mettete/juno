package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
)

// StateSelector contains the set of contract addresses and class hashes that a transaction
// touches during execution / state updates.
// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/storage/storage.py#L403
// Note: each txn type used to have it's own get_state_selecto() method, but now that info is in
// the execution info https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/business_logic/execution/deprecated_objects.py#L382
type StateSelector struct {
	ContractAddresses []felt.Felt
	ClassHashes       []felt.Felt
}

func get_state_selector_transactions(txns []*core.Transaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {

	contractAddresses := []felt.Felt{}
	classHashes := []felt.Felt{}

	for _, txn := range txns {
		stateSelector, err := get_state_selector_transaction(txn, generalConfig)
		if err != nil {
			return nil, err
		}

		for _, addr := range stateSelector.ContractAddresses {
			contractAddresses = append(contractAddresses, addr)
		}
		for _, hash := range stateSelector.ClassHashes {
			classHashes = append(classHashes, hash)
		}
	}
	return &StateSelector{
		ContractAddresses: contractAddresses,
		ClassHashes:       classHashes,
	}, nil
}

func get_state_selector_transaction(txn *core.Transaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	// Todo: push get_state_selectors to transaction methods??
	switch t := (*txn).(type) {
	case *core.InvokeTransaction:
		return get_state_selector_invoke(t, generalConfig)
	case *core.DeployTransaction:
		return get_state_selector_deploy(t, generalConfig)
	case *core.DeclareTransaction:
		return get_state_selector_declare(t, generalConfig)
	case *core.L1HandlerTransaction:
		return get_state_selector_l1handler(t, generalConfig)
	case *core.DeployAccountTransaction:
		return get_state_selector_deploy_account(t, generalConfig)
	default:
		return nil, errors.New("Unknown transaction type")
	}

}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/objects.py#L1421
func get_state_selector_invoke(txn *core.InvokeTransaction, config *StarknetGeneralConfig) (*StateSelector, error) {
	contractAddresses := []felt.Felt{}
	contractAddresses = append(contractAddresses, *txn.SenderAddress)

	if !txn.MaxFee.IsZero() {
		contractAddresses = append(contractAddresses, config.StarknetOsConfig.FeeTokenAddress)
	}

	return &StateSelector{
		ContractAddresses: contractAddresses,
		ClassHashes:       nil,
	}, nil
}

package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/rpc"
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
type TransactionExecutionInfo struct {
	ValidateInfo    *CallInfo
	CallInfo        *CallInfo
	FeeTransferInfo *CallInfo
	ActualFee       int64
	ActualResources ResourcesMapping
	TxType          *rpc.TransactionType
	RevertError     *string
}

type ResourcesMapping map[string]int

type CallInfo struct {
	CallerAddress       int64
	CallType            *CallType
	ContractAddress     int64
	ClassHash           *int64
	EntryPointSelector  *int64
	EntryPointType      *EntryPointType
	Calldata            []int64
	GasConsumed         int64
	FailureFlag         int64
	Retdata             []int64
	ExecutionResources  core.ExecutionResources
	Events              []OrderedEvent
	L2ToL1Messages      []OrderedL2ToL1Message
	StorageReadValues   []int64
	AccessedStorageKeys map[int64]struct{}
	InternalCalls       []*CallInfo
	CodeAddress         *int64
}
type CallType int

const (
	Call CallType = iota
	Delegate
)

type EntryPointType int

const (
	External EntryPointType = iota
	L1Handler
	Constructor
)

type OrderedEvent struct {
	Order int
	Keys  []int64
	Data  []int64
}

type OrderedL2ToL1Message struct {
	Order     int
	ToAddress int64
	Payload   []int64
}

func get_state_seelctor_execution_info(executionInfo []TransactionExecutionInfo, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	panic("Not implemented")
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/state_objects.py#L37
func get_state_selector_transactions(txns []*core.Transaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {

	contractAddresses := []felt.Felt{}
	classHashes := []felt.Felt{}

	for _, txn := range txns {
		stateSelector, err := get_state_selector_transaction(txn, generalConfig)
		if err != nil {
			return nil, err
		}
		contractAddresses = append(contractAddresses, stateSelector.ContractAddresses...)
		classHashes = append(classHashes, stateSelector.ClassHashes...)

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
	contractAddresses := []felt.Felt{*txn.SenderAddress}

	if !txn.MaxFee.IsZero() {
		contractAddresses = append(contractAddresses, config.StarknetOsConfig.FeeTokenAddress)
	}

	return &StateSelector{
		ContractAddresses: contractAddresses,
		ClassHashes:       nil,
	}, nil
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/objects.py#L1094
func get_state_selector_deploy(txn *core.DeployTransaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	return &StateSelector{
		ContractAddresses: []felt.Felt{*txn.ContractAddress},
		ClassHashes:       []felt.Felt{},
	}, nil
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/objects.py#L665-L671
func get_state_selector_declare(txn *core.DeclareTransaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	if txn.Version.Is(0) {
		return &StateSelector{}, nil
	}

	return &StateSelector{
		ContractAddresses: []felt.Felt{*txn.SenderAddress},
		ClassHashes:       []felt.Felt{*txn.ClassHash},
	}, nil
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/objects.py#L1611
func get_state_selector_l1handler(txn *core.L1HandlerTransaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	return &StateSelector{
		ContractAddresses: []felt.Felt{*txn.ContractAddress},
		ClassHashes:       []felt.Felt{},
	}, nil
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/objects.py#L856
func get_state_selector_deploy_account(txn *core.DeployAccountTransaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	return &StateSelector{
		ContractAddresses: []felt.Felt{*txn.ContractAddress},
		ClassHashes:       []felt.Felt{*txn.ClassHash},
	}, nil
}

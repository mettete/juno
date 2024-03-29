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

func (s *StateSelector) union(other *StateSelector) *StateSelector {
	contractAddresses := append(s.ContractAddresses, other.ContractAddresses...)
	classHashes := append(s.ClassHashes, other.ClassHashes...)

	return &StateSelector{
		ContractAddresses: contractAddresses,
		ClassHashes:       classHashes,
	}
}

// Todo: fix, this is not the correct way to do this
func (s *StateSelector) deleteContractAddress(address felt.Felt) {
	for i, a := range s.ContractAddresses {
		if a == address {
			s.ContractAddresses = append(s.ContractAddresses[:i], s.ContractAddresses[i+1:]...)
			break
		}
	}
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
	ContractAddress     felt.Felt
	ClassHash           *felt.Felt
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
	CodeAddress         *felt.Felt
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

func get_os_state_selector(
	txns []core.Transaction,
	executionInfos []TransactionExecutionInfo,
	generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	transactionStateSelector, err := get_state_selector_transactions(txns, generalConfig)
	if err != nil {
		return nil, err
	}
	executionStateSelector, err := get_state_seelctor_execution_info(executionInfos, generalConfig)
	if err != nil {
		return nil, err
	}
	// Include reserved contract addresses
	reservedStateSelector := &StateSelector{
		ContractAddresses: []felt.Felt{*new(felt.Felt).SetUint64(0), *new(felt.Felt).SetUint64(1)},
		ClassHashes:       []felt.Felt{},
	}
	// return union
	return transactionStateSelector.union(executionStateSelector).union(reservedStateSelector), nil
}

func get_state_seelctor_execution_info(executionInfos []TransactionExecutionInfo, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	contractAddresses := []felt.Felt{}
	classHashes := []felt.Felt{}
	for _, execInfo := range executionInfos {
		stateSelector, err := get_state_selector_exec_info(&execInfo, generalConfig)
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

func get_state_selector_exec_info(execInfo *TransactionExecutionInfo, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	nonOptionalCalls := non_optional_calls(execInfo)
	return get_state_selector_calls(nonOptionalCalls, generalConfig)
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/execution/objects.py#L324
func get_state_selector_calls(calls []*CallInfo, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {
	DEFAULT_DECLARE_SENDER_ADDRESS := new(felt.Felt).SetUint64(1)
	combinedSelector := &StateSelector{
		ContractAddresses: []felt.Felt{},
		ClassHashes:       []felt.Felt{},
	}

	for _, call := range calls {
		callAddress := call.ContractAddress
		if call.CodeAddress != nil {
			callAddress = *call.CodeAddress
		}
		if call.ClassHash == nil {
			return nil, errors.New("Class hash is missing from call info")
		}
		callSelector := StateSelector{
			ContractAddresses: []felt.Felt{call.ContractAddress, callAddress},
			ClassHashes:       []felt.Felt{*call.ClassHash},
		}
		// Todo: make cleaner
		callSelector.deleteContractAddress(*DEFAULT_DECLARE_SENDER_ADDRESS)

		interalSelector, err := get_state_selector_calls(call.InternalCalls, generalConfig)
		if err != nil {
			return nil, err
		}
		combinedSelector = combinedSelector.union(&callSelector)
		combinedSelector = combinedSelector.union(interalSelector)
	}

	return combinedSelector, nil
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/execution/objects.py#L493
func non_optional_calls(execInfo *TransactionExecutionInfo) []*CallInfo {
	var orderedOptionalCalls []*CallInfo
	if *execInfo.TxType == rpc.TxnDeployAccount {
		orderedOptionalCalls = []*CallInfo{execInfo.CallInfo, execInfo.ValidateInfo, execInfo.FeeTransferInfo}
	} else {
		orderedOptionalCalls = []*CallInfo{execInfo.ValidateInfo, execInfo.CallInfo, execInfo.FeeTransferInfo}
	}

	var nonOptionalCalls []*CallInfo
	for _, call := range orderedOptionalCalls {
		if call != nil {
			nonOptionalCalls = append(nonOptionalCalls, call)
		}
	}
	return nonOptionalCalls
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/v0.12.3/src/starkware/starknet/business_logic/transaction/state_objects.py#L37
func get_state_selector_transactions(txns []core.Transaction, generalConfig *StarknetGeneralConfig) (*StateSelector, error) {

	contractAddresses := []felt.Felt{}
	classHashes := []felt.Felt{}

	for _, txn := range txns {
		stateSelector, err := get_state_selector_transaction(&txn, generalConfig)
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

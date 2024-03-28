package osinput

import (
	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
)

// StarknetOsInput defines the input required to execute the OS.
// Note: this type will need adapted to the specific Runner being used (eg cairo-lang, snos, etc)
// starknet ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/os_input.py#L29
// snos ref: https://github.com/keep-starknet-strange/snos/blob/main/src/io/input.rs#L18
type StarknetOsInput struct {
	// Todo: Set of nodes that are changed in the contract/class trie from updating state?
	ContractStateCommitmentInfo CommitmentInfo `json:"contract_state_commitment_info"`
	ContractClassCommitmentInfo CommitmentInfo `json:"contract_class_commitment_info"`

	// New classes to be declared in the block
	DeprecatedCompiledClasses map[felt.Felt]core.Cairo0Class   `json:"deprecated_compiled_classes"`
	CompiledClasses           map[felt.Felt]core.CompiledClass `json:"compiled_classes"`

	// Todo ??
	CompiledClassVisitedPcs map[felt.Felt][]felt.Felt `json:"compiled_class_visited_pcs"`

	// Contract data associated with every contract that the batch of transactions require
	Contracts map[felt.Felt]ContractState `json:"contracts"`

	// Mapping from contract-tries class_hash to class-tries compiled_class_hash
	ClassHashToCompiledClassHash map[felt.Felt]felt.Felt `json:"class_hash_to_compiled_class_hash"`

	// Fixed Starknet Config
	GeneralConfig StarknetGeneralConfig `json:"general_config"`

	// New transactions defined in the block
	Transactions []core.Transaction `json:"transactions"`

	// Todo: initial, or final blockhash?
	BlockHash felt.Felt `json:"block_hash"`
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/storage/starknet_storage.py#L29
// Context: When updating a trie (class trie, contract trie), we keep track
// of the nodes, and thier (length,path,value) values, that change.
type CommitmentInfo struct {
	PreviousRoot    felt.Felt                 `json:"previous_root"`
	UpdatedRoot     felt.Felt                 `json:"updated_root"`
	TreeHeight      uint                      `json:"tree_height"`
	CommitmentFacts map[felt.Felt][]felt.Felt `json:"commitment_facts"`
}

type StarknetGeneralConfig struct {
	StarknetOsConfig        StarknetOsConfig   `json:"starknet_os_config"`
	GasPriceBounds          GasPriceBounds     `json:"gas_price_bounds"`
	InvokeTxMaxNSteps       int                `json:"invoke_tx_max_n_steps"`
	ValidateMaxNSteps       int                `json:"validate_max_n_steps"`
	DefaultEthPriceInFri    felt.Felt          `json:"default_eth_price_in_fri"`
	ConstantGasPrice        bool               `json:"constant_gas_price"`
	SequencerAddress        felt.Felt          `json:"sequencer_address"`
	CairoResourceFeeWeights map[string]float64 `json:"cairo_resource_fee_weights"`
	EnforceL1HandlerFee     bool               `json:"enforce_l1_handler_fee"`
	UseKzgDa                bool               `json:"use_kzg_da"`
}

type StarknetOsConfig struct {
	ChainID                   felt.Felt `json:"chain_id"`
	DeprecatedFeeTokenAddress felt.Felt `json:"deprecated_fee_token_address"`
	FeeTokenAddress           felt.Felt `json:"fee_token_address"`
}

type ContractState struct {
	ContractHash          []byte       `json:"contract_hash"`
	StorageCommitmentTree PatriciaTree `json:"storage_commitment_tree"`
	Nonce                 felt.Felt    `json:"nonce"`
}

type PatriciaTree struct {
	Root   string `json:"root"`
	Height uint   `json:"height"`
}

type GasPriceBounds struct {
	MinWeiL1GasPrice     felt.Felt `json:"min_wei_l1_gas_price"`
	MinFriL1GasPrice     felt.Felt `json:"min_fri_l1_gas_price"`
	MaxFriL1GasPrice     felt.Felt `json:"max_fri_l1_gas_price"`
	MinWeiL1DataGasPrice felt.Felt `json:"min_wei_l1_data_gas_price"`
	MinFriL1DataGasPrice felt.Felt `json:"min_fri_l1_data_gas_price"`
	MaxFriL1DataGasPrice felt.Felt `json:"max_fri_l1_data_gas_price"`
}

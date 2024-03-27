package osinput

import "github.com/NethermindEth/juno/starknet"

// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/core/os/os_input.py#L29
type StarknetOsInput struct {
	// Todo: ??
	ContractStateCommitmentInfo CommitmentInfo `json:"contract_state_commitment_info"`
	ContractClassCommitmentInfo CommitmentInfo `json:"contract_class_commitment_info"`

	// New classes to be declared in the block
	DeprecatedCompiledClasses map[int]starknet.Cairo0Definition `json:"deprecated_compiled_classes"`
	CompiledClasses           map[int]starknet.CompiledClass    `json:"compiled_classes"`

	// Todo ??
	CompiledClassVisitedPcs map[int][]int `json:"compiled_class_visited_pcs"`

	// Mapping from contract address to ContractState
	Contracts map[int]ContractState `json:"contracts"`

	ClassHashToCompiledClassHash map[int]int `json:"class_hash_to_compiled_class_hash"`

	// Fixed Starknet Config
	GeneralConfig StarknetGeneralConfig `json:"general_config"`

	// New transactions defined in the block
	Transactions []starknet.Transaction `json:"transactions"`

	// Todo: initial, or final blockhash?
	BlockHash int `json:"block_hash"`
}

// ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/starknet/storage/starknet_storage.py#L29
type CommitmentInfo struct {
	PreviousRoot    int           `json:"previous_root"`
	UpdatedRoot     int           `json:"updated_root"`
	TreeHeight      int           `json:"tree_height"`
	CommitmentFacts map[int][]int `json:"commitment_facts"`
}

type StarknetGeneralConfig struct {
	StarknetOsConfig        StarknetOsConfig   `json:"starknet_os_config"`
	GasPriceBounds          GasPriceBounds     `json:"gas_price_bounds"`
	InvokeTxMaxNSteps       int                `json:"invoke_tx_max_n_steps"`
	ValidateMaxNSteps       int                `json:"validate_max_n_steps"`
	DefaultEthPriceInFri    int                `json:"default_eth_price_in_fri"`
	ConstantGasPrice        bool               `json:"constant_gas_price"`
	SequencerAddress        int                `json:"sequencer_address"`
	CairoResourceFeeWeights map[string]float64 `json:"cairo_resource_fee_weights"`
	EnforceL1HandlerFee     bool               `json:"enforce_l1_handler_fee"`
	UseKzgDa                bool               `json:"use_kzg_da"`
}

type StarknetOsConfig struct {
	ChainID                   int `json:"chain_id"`
	DeprecatedFeeTokenAddress int `json:"deprecated_fee_token_address"`
	FeeTokenAddress           int `json:"fee_token_address"`
}

type ContractState struct {
	ContractHash          []byte       `json:"contract_hash"`
	StorageCommitmentTree PatriciaTree `json:"storage_commitment_tree"`
	Nonce                 int          `json:"nonce"`
}

type PatriciaTree struct {
	Root   string `json:"root"`
	Height int    `json:"height"`
}

type GasPriceBounds struct {
	MinWeiL1GasPrice     int `json:"min_wei_l1_gas_price"`
	MinFriL1GasPrice     int `json:"min_fri_l1_gas_price"`
	MaxFriL1GasPrice     int `json:"max_fri_l1_gas_price"`
	MinWeiL1DataGasPrice int `json:"min_wei_l1_data_gas_price"`
	MinFriL1DataGasPrice int `json:"min_fri_l1_data_gas_price"`
	MaxFriL1DataGasPrice int `json:"max_fri_l1_data_gas_price"`
}

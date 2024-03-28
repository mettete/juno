package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/rpc"
)

// Note: assumes the block and state have already been stored
func PopulateOSInput(block core.Block, handler *rpc.Handler) (*StarknetOsInput, error) {
	osinput := &StarknetOsInput{}

	// Todo
	contractAddresses := getContractDataThatTxsUse(handler, osinput.Transactions)

	// Todo: commitment info
	contractStateCommitmentInfo := getContractStateCommitmentInfo()
	classStateCommitmentInfo := getClassStateCommitmentInfo()
	osinput.ContractStateCommitmentInfo = contractStateCommitmentInfo
	osinput.ContractClassCommitmentInfo = classStateCommitmentInfo

	// populate transactions, deprecrated classes, and classes
	for _, tx := range block.Transactions {
		osinput.Transactions = append(osinput.Transactions, tx)

		switch decTxn := tx.(type) {
		case *core.DeclareTransaction:
			class, err := handler.Class(rpc.BlockID{Hash: block.Hash}, *decTxn.ClassHash)
			if err != nil {
				return nil, errors.New(err.Message) // Todo: hanlde error correctly
			}
			if class.Program != "" {
				depCompClass := AdaptClassToDeprecatedCompiledClass(class)
				osinput.DeprecatedCompiledClasses[*decTxn.ClassHash] = depCompClass
			} else {
				compClass := AdaptClassToCompiledClass(class)
				osinput.CompiledClasses[*decTxn.ClassHash] = compClass
			}
		}
	}

	// Todo: ClassHashToCompiledClassHash
	classHashToCompiledClassHash := getClassHashToCompiledClassHash(handler, contractAddresses)
	osinput.ClassHashToCompiledClassHash = classHashToCompiledClassHash

	// Todo: CompiledClassVisitedPcs??

	// Todo: contracts
	contracts := getContracts(contractAddresses)
	osinput.Contracts = contracts

	osinput.GeneralConfig = loadExampleStarknetOSConfig()

	osinput.BlockHash = *block.Hash
	return osinput, nil
}

func AdaptClassToDeprecatedCompiledClass(class *rpc.Class) core.Cairo0Class {
	panic("Todo")
}
func AdaptClassToCompiledClass(class *rpc.Class) core.CompiledClass {
	panic("Todo")
}

func loadExampleStarknetOSConfig() StarknetGeneralConfig {
	seqAddr, _ := new(felt.Felt).SetString("0x6c95526293b61fa708c6cba66fd015afee89309666246952456ab970e9650aa")
	feeTokenAddr, _ := new(felt.Felt).SetString("0x6c95526293b61fa708c6cba66fd015afee89309666246952456ab970e9650aa")
	chainID := new(felt.Felt).SetBytes([]byte("TODO"))
	return StarknetGeneralConfig{
		StarknetOsConfig: StarknetOsConfig{
			ChainID:         *chainID,
			FeeTokenAddress: *feeTokenAddr,
		},
		InvokeTxMaxNSteps: 3000000,
		ValidateMaxNSteps: 1000000,
		ConstantGasPrice:  false,
		SequencerAddress:  *seqAddr,
		CairoResourceFeeWeights: map[string]float64{
			"poseidon_builtin":    0.0,
			"n_steps":             1.0,
			"ecdsa_builtin":       0.0,
			"keccak_builtin":      0.0,
			"range_check_builtin": 0.0,
			"pedersen_builtin":    0.0,
			"output_builtin":      0.0,
			"bitwise_builtin":     0.0,
			"ec_op_builtin":       0.0,
		},
		EnforceL1HandlerFee: true,
		UseKzgDa:            false,
	}
}

// Todo: using handler as proxy for state access
func getContractDataThatTxsUse(handler *rpc.Handler, txs []core.Transaction) []felt.Felt {
	// Todo: Given a set of transactions, and acess to contract-trie,
	// return the set of contract addresses (and contract class hashes)
	// Note: cairol-lang uses get_state_selector() to return subset of Merkle Trie the txn affects
	panic("unimplemented")
}

// Todo: using handler as proxy for state access
func getClassHashToCompiledClassHash(handler *rpc.Handler, classHashes []felt.Felt) map[felt.Felt]felt.Felt {
	// Todo: Given access to the class trie, and a set of classhashes, return the associated compiled class hashes
	panic("unimplemented")
}

func getContractStateCommitmentInfo() CommitmentInfo {

	// Todo: Given the old and new contract Trie, collect all the
	// nodes that were modified

	panic("unimplemented")
}

func getClassStateCommitmentInfo() CommitmentInfo {

	// Todo: Given the old and new class Trie, collect all the
	// nodes that were modified

	panic("unimplemented")
}

func getContracts(contractAddresses []felt.Felt) map[felt.Felt]ContractState {
	// Todo: given a batch of transactions, collect the set of ContractState's
	// for every contract that the contract touched
	panic("unimplemented")
}

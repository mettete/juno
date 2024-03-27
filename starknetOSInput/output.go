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
				depCompClass := AdaptClassToDeprecatedCompileClass(class)
				osinput.DeprecatedCompiledClasses[*decTxn.ClassHash] = depCompClass
			} else {
				compClass := AdaptClassToCompileClass(class)
				osinput.CompiledClasses[*decTxn.ClassHash] = compClass
			}
		}
	}

	osinput.GeneralConfig = loadExampleStarknetOSConfig()

	osinput.BlockHash = *block.Hash
	return osinput, nil
}

func AdaptClassToDeprecatedCompileClass(class *rpc.Class) core.Cairo0Class {
	panic("Todo")
	return core.Cairo0Class{}
}
func AdaptClassToCompileClass(class *rpc.Class) core.CompiledClass {
	panic("Todo")
	return core.CompiledClass{}
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

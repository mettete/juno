package osinput

import (
	"errors"

	"github.com/NethermindEth/juno/core"
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

	// ClassHashToCompiledClassHash

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

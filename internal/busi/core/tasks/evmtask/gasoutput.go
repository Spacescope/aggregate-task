package evmtask

import (
	"aggregate-task/pkg/models/evmmodel"
	"context"

	log "github.com/sirupsen/logrus"
)

type GasOutput struct {
}

func (g *GasOutput) Name() string {
	return "evm_derived_gas_outputs"
}

func (g *GasOutput) Model() interface{} {
	return new(evmmodel.GasOutputs)
}

func (g *GasOutput) Run(ctx context.Context, height uint64 /*, version int*/) error {
	log.Info("abc")
	return nil
}

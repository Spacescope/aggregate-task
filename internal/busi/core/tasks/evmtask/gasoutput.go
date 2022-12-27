package evmtask

import (
	"context"

	"aggregate-task/internal/busi/core/tasks/evmtask/gasoutput"
	"aggregate-task/pkg/models/evmmodel"
	"aggregate-task/pkg/utils"

	"github.com/filecoin-project/go-state-types/big"

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

func (g *GasOutput) Run(ctx context.Context, height int64) error {
	btr := new(gasoutput.BlockTransactionReceipt)
	if err := btr.GetBlockTransactionReceipt(ctx, height); err != nil {
		log.Error("GetBlockTransactionReceipt[%v] error: %v", height, err)
		return err
	}

	gasOutputSlice := make([]*evmmodel.GasOutputs, 0)

	var version uint32

	for btr.HashNext() {
		t, r := btr.Next()

		baseFeePerGas, _ := big.FromString(btr.BlockHeader.BaseFeePerGas)
		gasFeeCap, _ := big.FromString(t.MaxFeePerGas)
		gasPremium, _ := big.FromString(t.MaxPriorityFeePerGas)

		gasOutputs := gasoutput.ComputeGasOutputs(r.GasUsed, int64(t.GasLimit), baseFeePerGas, gasFeeCap, gasPremium, true)

		gasOutput := evmmodel.GasOutputs{
			Height:        height,
			Version:       version,
			StateRoot:     btr.BlockHeader.StateRoot,
			ParentBaseFee: btr.BlockHeader.BaseFeePerGas,

			Cid:        t.Hash,
			From:       t.From,
			To:         t.To,
			Value:      t.Value,
			GasFeeCap:  t.MaxFeePerGas,
			GasPremium: t.MaxPriorityFeePerGas,
			GasLimit:   int64(t.GasLimit),
			Nonce:      t.Nonce,
			Method:     t.Type,

			Status:             r.Status,
			GasUsed:            r.GasUsed,
			BaseFeeBurn:        gasOutputs.BaseFeeBurn.String(),
			OverEstimationBurn: gasOutputs.OverEstimationBurn.String(),
			MinerPenalty:       gasOutputs.MinerPenalty.String(),
			MinerTip:           gasOutputs.MinerTip.String(),
			Refund:             gasOutputs.Refund.String(),
			GasRefund:          gasOutputs.GasRefund,
			GasBurned:          gasOutputs.GasBurned,
		}

		gasOutputSlice = append(gasOutputSlice, &gasOutput)
	}

	if len(gasOutputSlice) > 0 {
		if err := utils.BatchWrite(ctx, new(evmmodel.GasOutputs), height, &gasOutputSlice); err != nil {
			log.Errorf("BatchWrite[%v] error: %v", height, err)
			return err
		}
	}

	return nil
}

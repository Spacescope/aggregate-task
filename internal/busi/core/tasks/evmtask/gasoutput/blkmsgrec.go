package gasoutput

import (
	"aggregate-task/pkg/models/dependentmodel"
	"context"

	log "github.com/sirupsen/logrus"

	"aggregate-task/pkg/utils"
)

type BlockTransactionReceipt struct {
	BlockHeader dependentmodel.BlockHeader
	Transaction []dependentmodel.Transaction
	Receipt     []dependentmodel.Receipt
	idx         int
}

func (btr *BlockTransactionReceipt) getEVMBlockHeader(ctx context.Context, height int64) (*dependentmodel.BlockHeader, bool, error) {
	evmBlockHeader := new(dependentmodel.BlockHeader)
	b, err := utils.EngineGroup[utils.DBOBTask].Where("height = ?", height).Get(evmBlockHeader)
	if err != nil {
		log.Infof("[%v] Executed sql error: %v", height, err)
		return nil, false, err
	}

	if b {
		return evmBlockHeader, true, nil
	} else {
		return nil, false, nil
	}
}

func (btr *BlockTransactionReceipt) getEVMTransaction(ctx context.Context, height int64) ([]*dependentmodel.Transaction, error) {
	evmTransactions := make([]*dependentmodel.Transaction, 0)
	if err := utils.EngineGroup[utils.DBOBTask].Where("height = ?", height).Find(&evmTransactions); err != nil {
		log.Infof("[%v] Executed sql error: %v", height, err)
		return nil, err
	}
	return evmTransactions, nil
}

func (btr *BlockTransactionReceipt) getEVMReceipt(ctx context.Context, height int64, txnHash string) (*dependentmodel.Receipt, bool) {
	evmReceipt := new(dependentmodel.Receipt)
	b, err := utils.EngineGroup[utils.DBOBTask].Where("height = ? and transaction_hash = ?", height, txnHash).Get(evmReceipt)
	if err != nil {
		log.Infof("[%v] Executed sql error: %v", height, err)
		return nil, false
	}

	if b {
		return evmReceipt, true
	} else {
		return nil, false
	}
}

func (btr *BlockTransactionReceipt) GetBlockTransactionReceipt(ctx context.Context, height int64) error {
	bh, b, err := btr.getEVMBlockHeader(ctx, height)
	if err != nil {
		return err
	}
	if !b {
		return nil
	}
	btr.BlockHeader = *bh

	txns, err := btr.getEVMTransaction(ctx, height)
	if err != nil {
		return err
	}

	for _, txn := range txns {
		receipt, b := btr.getEVMReceipt(ctx, height, txn.Hash)
		if !b {
			continue
		}

		btr.Transaction = append(btr.Transaction, *txn)
		btr.Receipt = append(btr.Receipt, *receipt)
	}

	return nil
}

func (btr *BlockTransactionReceipt) HashNext() bool {
	if btr.idx < len(btr.Transaction) {
		return true
	}

	return false
}

func (btr *BlockTransactionReceipt) Next() (*dependentmodel.Transaction, *dependentmodel.Receipt) {
	if btr.HashNext() {
		txn := btr.Transaction[btr.idx]
		receipt := btr.Receipt[btr.idx]
		btr.idx++
		return &txn, &receipt
	}

	return nil, nil
}

// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package mtnode

import (
	"context"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionPublisher interface {
	PublishTransaction(ctx context.Context, tx *types.Transaction) error
	CheckHealth(ctx context.Context) error
	Initialize(context.Context) error
	Start(context.Context) error
	StopAndWait()
}

type MtInterface struct {
	txStreamer  *TransactionStreamer
	txPublisher TransactionPublisher
	mtNode      *Node
}

func NewMtInterface(txStreamer *TransactionStreamer, txPublisher TransactionPublisher) (*MtInterface, error) {
	return &MtInterface{
		txStreamer:  txStreamer,
		txPublisher: txPublisher,
	}, nil
}

func (a *MtInterface) Initialize(n *Node) {
	a.mtNode = n
}

func (a *MtInterface) PublishTransaction(ctx context.Context, tx *types.Transaction) error {
	return a.txPublisher.PublishTransaction(ctx, tx)
}

func (a *MtInterface) TransactionStreamer() *TransactionStreamer {
	return a.txStreamer
}

func (a *MtInterface) BlockChain() *core.BlockChain {
	return a.txStreamer.bc
}

func (a *MtInterface) MtNode() interface{} {
	return a.mtNode
}

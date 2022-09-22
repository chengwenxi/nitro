//
// Copyright 2022, Mantlenetwork, Inc. All rights reserved.
//

package precompiles

// This precompile represents MtOS's internal actions as calls it makes to itself
type MtosActs struct {
	Address addr // 0xa4b05

	CallerNotMtOSError func() error
}

func (con MtosActs) StartBlock(c ctx, evm mech, l1BaseFee huge, l1BlockNumber, l2BlockNumber, timeLastBlock uint64) error {
	return con.CallerNotMtOSError()
}

func (con MtosActs) BatchPostingReport(c ctx, evm mech, batchTimestamp huge, batchPosterAddress addr, batchNumber uint64, batchDataGas uint64, l1BaseFeeWei huge) error {
	return con.CallerNotMtOSError()
}

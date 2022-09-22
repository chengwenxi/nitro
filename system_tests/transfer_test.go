// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantlenetworkio/mantle/blob/main/LICENSE

package mttest

import (
	"context"
	"fmt"
	"math/big"
	"testing"
)

func TestTransfer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l2info, _, client, l2stack := CreateTestL2(t, ctx)
	defer requireClose(t, l2stack)

	l2info.GenerateAccount("User2")

	tx := l2info.PrepareTx("Owner", "User2", l2info.TransferGas, big.NewInt(1e12), nil)

	err := client.SendTransaction(ctx, tx)
	Require(t, err)

	_, err = EnsureTxSucceeded(ctx, client, tx)
	Require(t, err)

	bal, err := client.BalanceAt(ctx, l2info.GetAddress("Owner"), nil)
	Require(t, err)
	fmt.Println("Owner balance is: ", bal)
	bal2, err := client.BalanceAt(ctx, l2info.GetAddress("User2"), nil)
	Require(t, err)
	if bal2.Cmp(big.NewInt(1e12)) != 0 {
		Fail(t, "Unexpected recipient balance: ", bal2)
	}
}

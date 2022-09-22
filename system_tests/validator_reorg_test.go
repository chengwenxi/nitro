// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantlenetworkio/mantle/blob/main/LICENSE

//go:build validatorreorgtest
// +build validatorreorgtest

package mttest

import "testing"

func TestBlockValidatorReorg(t *testing.T) {
	testSequencerInboxReaderImpl(t, true)
}

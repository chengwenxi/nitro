// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

//go:build validatorreorgtest
// +build validatorreorgtest

package arbtest

import "testing"

func TestBlockValidatorReorg(t *testing.T) {
	testSequencerInboxReaderImpl(t, true)
}

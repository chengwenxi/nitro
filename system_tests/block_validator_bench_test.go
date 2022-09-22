// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantlenetworkio/mantle/blob/main/LICENSE

// race detection makes things slow and miss timeouts
//go:build block_validator_bench
// +build block_validator_bench

package mttest

import (
	"testing"

	"github.com/mantlenetworkio/mantle/das"
)

func TestBlockValidatorBenchmark(t *testing.T) {
	testBlockValidatorSimple(t, das.OnchainDataAvailabilityString, true)
}

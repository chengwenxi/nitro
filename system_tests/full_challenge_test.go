// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantlenetworkio/mantle/blob/main/LICENSE

//go:build fullchallengetest
// +build fullchallengetest

//
// Copyright 2021-2022, Mantlenetwork, Inc. All rights reserved.
//

package mttest

import (
	"testing"
)

func TestFullChallengeAsserterIncorrect(t *testing.T) {
	RunChallengeTest(t, false)
}

func TestFullChallengeAsserterCorrect(t *testing.T) {
	RunChallengeTest(t, true)
}

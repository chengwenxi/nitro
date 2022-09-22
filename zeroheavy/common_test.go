// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package zeroheavy

import (
	"testing"

	"github.com/mantlenetworkio/mantle/util/testhelpers"
)

func Require(t *testing.T, err error, printables ...interface{}) {
	t.Helper()
	testhelpers.RequireImpl(t, err, printables...)
}

func Fail(t *testing.T, printables ...interface{}) {
	t.Helper()
	testhelpers.FailImpl(t, printables...)
}

func ShowError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

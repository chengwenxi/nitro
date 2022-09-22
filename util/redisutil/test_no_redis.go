// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

//go:build !redistest
// +build !redistest

package redisutil

import "testing"

// t param is used to make sure this is only called in tests
func GetTestRedisURL(t *testing.T) string {
	return ""
}

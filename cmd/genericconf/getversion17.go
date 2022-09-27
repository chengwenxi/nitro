// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

//go:build !go1.18

package genericconf

func GetVersion() (string, string) {
	return "development", "development"
}

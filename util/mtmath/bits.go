// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mtmath

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

type bytes32 = common.Hash

// flips the nth bit in an ethereum word, starting from the left
func FlipBit(data bytes32, bit byte) bytes32 {
	data[bit/8] ^= 1 << (7 - bit%8)
	return data
}

// unrolls a series of slices into a singular, concatenated slice
func ConcatByteSlices(slices ...[]byte) []byte {
	unrolled := []byte{}
	for _, slice := range slices {
		unrolled = append(unrolled, slice...)
	}
	return unrolled
}

// the number of eth-words needed to store n bytes
func WordsForBytes(nbytes uint64) uint64 {
	return (nbytes + 31) / 32
}

// casts a uint64 to its big-endian representation
func UintToBytes(value uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, value)
	return result
}

// casts a uint32 to its big-endian representation
func Uint32ToBytes(value uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, value)
	return result
}

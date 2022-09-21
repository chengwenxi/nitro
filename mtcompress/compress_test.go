// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package mtcompress

import (
	"bytes"
	"testing"

	"github.com/mantlenetworkio/mantle/util/testhelpers"
)

func testDecompress(t *testing.T, compressed, decompressed []byte) {
	res, err := Decompress(compressed, len(decompressed)*2+64)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res, decompressed) {
		t.Fatal("results differ ", res, " vs. ", decompressed)
	}
}

func testCompressDecompress(t *testing.T, data []byte) {
	compressedWell, err := CompressWell(data)
	if err != nil {
		t.Fatal(err)
	}
	testDecompress(t, compressedWell, data)

	compressedFast, err := CompressFast(data)
	if err != nil {
		t.Fatal(err)
	}
	testDecompress(t, compressedFast, data)
}

func TestMtCompress(t *testing.T) {
	asciiData := []byte("This is a long and repetitive string. Yadda yadda yadda yadda yadda. The quick brown fox jumped over the lazy dog.")
	for i := 0; i < 8; i++ {
		asciiData = append(asciiData, asciiData...)
	}
	testCompressDecompress(t, asciiData)

	source := testhelpers.NewPseudoRandomDataSource(t, 0)
	randData := source.GetData(2500)
	testCompressDecompress(t, randData)

	// test empty data:
	testCompressDecompress(t, []byte{})
}

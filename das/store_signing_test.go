// Copyright 2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package das

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestStoreSigning(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	Require(t, err)

	publicKey := privateKey.Public()
	addr := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	weirdMessage := []byte("The quick brown fox jumped over the lazy dog.")
	timeout := uint64(time.Now().Unix())

	signer := DasSignerFromPrivateKey(privateKey)
	sig, err := applyDasSigner(signer, weirdMessage, timeout)
	Require(t, err)

	recoveredAddr, err := DasRecoverSigner(weirdMessage, timeout, sig)
	Require(t, err)

	if recoveredAddr != addr {
		t.Fatal()
	}
}

// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package validator

/*
#cgo CFLAGS: -g -Wall -I../target/include/
#include "mtitrator.h"

extern ResolvedPreimage preimageResolver(size_t context, const uint8_t* hash);

ResolvedPreimage preimageResolverC(size_t context, const uint8_t* hash) {
  return preimageResolver(context, hash);
}
*/
import "C"

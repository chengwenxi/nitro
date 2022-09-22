//
// Copyright 2021-2022, Mantlenetwork, Inc. All rights reserved.
//

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/mantlenetworkio/mantle/mtnode"
	"github.com/mantlenetworkio/mantle/mtutil"
	"github.com/mantlenetworkio/mantle/util/redisutil"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: seq-coordinator-invalidate [redis url] [signing key] [msg index]\n")
		os.Exit(1)
	}
	redisUrl := os.Args[1]
	signingKey := os.Args[2]
	msgIndex, err := strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		panic("Failed to parse msg index: " + err.Error())
	}
	redisClient, err := redisutil.RedisClientFromURL(redisUrl)
	if err != nil {
		panic(err)
	}
	if redisClient == nil {
		panic("redis url not defined")
	}
	err = mtnode.StandaloneSeqCoordinatorInvalidateMsgIndex(context.Background(), redisClient, signingKey, mtutil.MessageIndex(msgIndex))
	if err != nil {
		panic(err)
	}
}

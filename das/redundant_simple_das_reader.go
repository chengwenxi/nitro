// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package das

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mantlenetworkio/mantle/mtstate"
	"github.com/mantlenetworkio/mantle/util/pretty"
)

type RedundantSimpleDASReader struct {
	inners []mtstate.DataAvailabilityReader
}

func NewRedundantSimpleDASReader(inners []mtstate.DataAvailabilityReader) mtstate.DataAvailabilityReader {
	return &RedundantSimpleDASReader{inners}
}

type rsdrResponse struct {
	data []byte
	err  error
}

func (r RedundantSimpleDASReader) GetByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	log.Trace("das.RedundantSimpleDASReader.GetByHash", "key", pretty.PrettyHash(hash), "this", r)

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	numPending := len(r.inners)
	results := make(chan rsdrResponse, numPending)
	for _, inner := range r.inners {
		go func(inn mtstate.DataAvailabilityReader) {
			res, err := inn.GetByHash(subCtx, hash)
			results <- rsdrResponse{res, err}
		}(inner)
	}
	var anyError error
	for numPending > 0 {
		select {
		case res := <-results:
			if res.err != nil {
				anyError = res.err
				numPending--
			} else {
				return res.data, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, anyError
}

func (r RedundantSimpleDASReader) HealthCheck(ctx context.Context) error {
	for _, simpleDASReader := range r.inners {
		err := simpleDASReader.HealthCheck(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedundantSimpleDASReader) ExpirationPolicy(ctx context.Context) (mtstate.ExpirationPolicy, error) {
	// If at least one inner service has KeepForever,
	// then whole redundant service can serve after timeout.

	// If no inner service has KeepForever,
	// but at least one inner service has DiscardAfterArchiveTimeout,
	// then whole redundant service can serve till archive timeout.

	// If no inner service has KeepForever, DiscardAfterArchiveTimeout,
	// but at least one inner service has DiscardAfterDataTimeout,
	// then whole redundant service can serve till data timeout.
	var res mtstate.ExpirationPolicy = -1
	for _, serv := range r.inners {
		expirationPolicy, err := serv.ExpirationPolicy(ctx)
		if err != nil {
			return -1, err
		}
		switch expirationPolicy {
		case mtstate.KeepForever:
			return mtstate.KeepForever, nil
		case mtstate.DiscardAfterArchiveTimeout:
			res = mtstate.DiscardAfterArchiveTimeout
		case mtstate.DiscardAfterDataTimeout:
			if res != mtstate.DiscardAfterArchiveTimeout {
				res = mtstate.DiscardAfterDataTimeout
			}
		}
	}
	if res == -1 {
		return -1, errors.New("unknown expiration policy")
	}
	return res, nil
}

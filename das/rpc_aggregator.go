// Copyright 2021-2022, Mantlenetwork, Inc.
// For license information, see https://github.com/mantle/blob/master/LICENSE

package das

import (
	"context"
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/mantlenetworkio/mantle/mtutil"
	"github.com/mantlenetworkio/mantle/solgen/go/bridgegen"

	"github.com/ethereum/go-ethereum/common"
)

type BackendConfig struct {
	URL                 string `json:"url"`
	PubKeyBase64Encoded string `json:"pubkey"`
	SignerMask          uint64 `json:"signermask"`
}

func NewRPCAggregator(ctx context.Context, config DataAvailabilityConfig) (*Aggregator, error) {
	services, err := setUpServices(config)
	if err != nil {
		return nil, err
	}
	return NewAggregator(ctx, config, services)
}

func NewRPCAggregatorWithL1Info(config DataAvailabilityConfig, l1client mtutil.L1Interface, seqInboxAddress common.Address) (*Aggregator, error) {
	services, err := setUpServices(config)
	if err != nil {
		return nil, err
	}
	return NewAggregatorWithL1Info(config, services, l1client, seqInboxAddress)
}

func NewRPCAggregatorWithSeqInboxCaller(config DataAvailabilityConfig, seqInboxCaller *bridgegen.SequencerInboxCaller) (*Aggregator, error) {
	services, err := setUpServices(config)
	if err != nil {
		return nil, err
	}
	return NewAggregatorWithSeqInboxCaller(config, services, seqInboxCaller)
}

func setUpServices(config DataAvailabilityConfig) ([]ServiceDetails, error) {
	var cs []BackendConfig
	err := json.Unmarshal([]byte(config.AggregatorConfig.Backends), &cs)
	if err != nil {
		return nil, err
	}

	var services []ServiceDetails

	for _, b := range cs {
		url, err := url.Parse(b.URL)
		if err != nil {
			return nil, err
		}
		// Prometheus metric names must contain only chars [a-zA-Z0-9:_]
		invalidPromCharRegex := regexp.MustCompile(`[^a-zA-Z0-9:_]+`)
		metricName := invalidPromCharRegex.ReplaceAllString(url.Hostname(), "_")

		service, err := NewDASRPCClient(b.URL)
		if err != nil {
			return nil, err
		}

		pubKey, err := DecodeBase64BLSPublicKey([]byte(b.PubKeyBase64Encoded))
		if err != nil {
			return nil, err
		}

		d, err := NewServiceDetails(service, *pubKey, uint64(b.SignerMask), metricName)
		if err != nil {
			return nil, err
		}

		services = append(services, *d)
	}

	return services, nil
}

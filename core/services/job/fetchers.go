package job

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Fetcher interface {
	PipelineStage
	Fetch() (interface{}, error)
}

type FetcherType string

var (
	FetcherTypeBridge FetcherType = "bridge"
	FetcherTypeHttp   FetcherType = "http"
	FetcherTypeMedian FetcherType = "median"
)

type BaseFetcher struct {
	// ID           uint64        `json:"-" gorm:"primary_key"`

	notifiee Notifiee `gorm:"-"`

	// The following fields are mutually exclusive.  This is enforced by a DB constraint.
	// OffchainReportingJobID models.ID `json:"-"`
	// FluxMonitorJobID       models.ID `json:"-"`
	// DirectRequestJobID     models.ID `json:"-"`
}

// func (f BaseFetcher) GetID() uint64 {
// 	return f.ID
// }

// func (f *BaseFetcher) SetNotifiee(n Notifiee) {
// 	f.notifiee = n
// 	for _, transformer := range f.Transformers {
// 		transformer.SetNotifiee(n)
// 	}
// }

type Fetchers []Fetcher

func (f Fetchers) SetNotifiee(n Notifiee) {
	for _, fetcher := range f {
		fetcher.SetNotifiee(n)
	}
}

func (f *Fetchers) UnmarshalJSON(bs []byte) (err error) {
	var spec []json.RawMessage
	err = json.Unmarshal(bs, &spec)
	if err != nil {
		return err
	}

	for _, fetcherBytes := range spec {
		fetcher, err := UnmarshalFetcherJSON([]byte(fetcherBytes))
		if err != nil {
			return err
		}
		*f = append(*f, fetcher)
	}
	return nil
}

func UnmarshalFetcherJSON(bs []byte) (_ Fetcher, err error) {
	var header struct {
		Type FetcherType `json:"type"`
	}
	err = json.Unmarshal(bs, &header)
	if err != nil {
		return nil, err
	}

	var fetcher Fetcher
	switch header.Type {
	case FetcherTypeBridge:
		bridgeFetcher := BridgeFetcher{}
		err = json.Unmarshal(bs, &bridgeFetcher)
		if err != nil {
			return nil, err
		}
		fetcher = &bridgeFetcher

	case FetcherTypeHttp:
		httpFetcher := HttpFetcher{}
		err = json.Unmarshal(bs, &httpFetcher)
		if err != nil {
			return nil, err
		}
		fetcher = &httpFetcher

	case FetcherTypeMedian:
		medianFetcher := MedianFetcher{}
		err = json.Unmarshal(bs, &medianFetcher)
		if err != nil {
			return nil, err
		}
		fetcher = &medianFetcher

	default:
		return nil, errors.New("unknown fetcher type")
	}

	return fetcher, nil
}
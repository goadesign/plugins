package fetcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/v3/http"
	archiversvc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/archiver"
	archiverkc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/archiver/kitclient"
	archiverks "goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/archiver/server"
	fetchersvc "goa.design/plugins/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// fetcher service example implementation.
// The example methods log the requests and return zero values.
type fetchersvcsvc struct {
	logger  log.Logger
	archive endpoint.Endpoint
}

// NewFetcher returns the fetcher service implementation.
func NewFetcher(logger log.Logger, archiverHost string) fetchersvc.Service {
	u := url.URL{
		Scheme: "http",
		Host:   archiverHost,
		Path:   archiverks.ArchiveArchiverPath(),
	}
	var (
		dec = goahttp.ResponseDecoder
		enc = goahttp.RequestEncoder
	)
	arc := kithttp.NewClient(
		"POST",
		&u,
		archiverkc.EncodeArchiveRequest(enc),
		archiverkc.DecodeArchiveResponse(dec),
	)
	return &fetchersvcsvc{logger: logger, archive: arc.Endpoint()}
}

// Fetch makes a GET request to the given URL and stores the results in the
// archiver service which must be running or the request fails
func (s *fetchersvcsvc) Fetch(ctx context.Context, p *fetchersvc.FetchPayload) (*fetchersvc.FetchMedia, error) {
	// Make request to external endpoint
	resp, err := http.Get(p.URL)
	if err != nil {
		return nil, fetchersvc.MakeBadRequest(fmt.Errorf("bad request URL: %s", err))
	}

	// Read response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fetchersvc.MakeBadRequest(fmt.Errorf("failed to decode response: %s", err))
	}

	// Archive response using archiver service
	res, err := s.archive(ctx, &archiversvc.ArchivePayload{Status: resp.StatusCode, Body: string(body)})
	if err != nil {
		return nil, fetchersvc.MakeBadRequest(fmt.Errorf("failed to decode response: %s", err))
	}

	// Return response
	media := res.(*archiversvc.ArchiveMedia)
	return &fetchersvc.FetchMedia{
		ArchiveHref: media.Href,
		Status:      media.Status,
	}, nil
}

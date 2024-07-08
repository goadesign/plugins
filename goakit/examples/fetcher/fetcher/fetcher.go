package fetcherapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	goahttp "goa.design/goa/v3/http"

	genarchiver "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/archiver"
	genarchiverclient "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/http/archiver/kitclient"
	genarchiverserver "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/http/archiver/server"
	genfetcher "goa.design/plugins/v3/goakit/examples/fetcher/fetcher/gen/fetcher"
)

// fetcher service example implementation.
// The example methods log the requests and return zero values.
type fetchersvc struct {
	logger  log.Logger
	archive endpoint.Endpoint
}

// NewFetcher returns the fetcher service implementation.
func NewFetcher(logger log.Logger, archiverHost string) genfetcher.Service {
	u := url.URL{
		Scheme: "http",
		Host:   archiverHost,
		Path:   genarchiverserver.ArchiveArchiverPath(),
	}
	var (
		dec = goahttp.ResponseDecoder
		enc = goahttp.RequestEncoder
	)
	arc := kithttp.NewClient(
		"POST",
		&u,
		genarchiverclient.EncodeArchiveRequest(enc),
		genarchiverclient.DecodeArchiveResponse(dec),
	)
	return &fetchersvc{logger: logger, archive: arc.Endpoint()}
}

// Fetch makes a GET request to the given URL and stores the results in the
// archiver service which must be running or the request fails
func (s *fetchersvc) Fetch(ctx context.Context, p *genfetcher.FetchPayload) (*genfetcher.FetchMedia, error) {
	// Make request to external endpoint
	resp, err := http.Get(p.URL)
	if err != nil {
		return nil, genfetcher.MakeBadRequest(fmt.Errorf("bad request URL: %s", err))
	}

	// Read response
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, genfetcher.MakeBadRequest(fmt.Errorf("failed to decode response: %s", err))
	}

	// Archive response using archiver service
	res, err := s.archive(ctx, &genarchiver.ArchivePayload{Status: resp.StatusCode, Body: string(body)})
	if err != nil {
		return nil, genfetcher.MakeBadRequest(fmt.Errorf("failed to decode response: %s", err))
	}

	// Return response
	media := res.(*genarchiver.ArchiveMedia)
	return &genfetcher.FetchMedia{
		ArchiveHref: media.Href,
		Status:      media.Status,
	}, nil
}

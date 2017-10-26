package fetcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	goahttp "goa.design/goa/http"
	archiversvc "goa.design/plugins/goakit/examples/client/archiver/gen/archiver"
	archiverkc "goa.design/plugins/goakit/examples/client/archiver/gen/http/archiver/kitclient"
	archiverks "goa.design/plugins/goakit/examples/client/archiver/gen/http/archiver/server"
	"goa.design/plugins/goakit/examples/client/fetcher/gen/fetcher"
)

// fetcher service implementation.
type fetchersvc struct {
	logger  log.Logger
	archive endpoint.Endpoint
}

// NewFetcher returns the fetcher service implementation.
func NewFetcher(logger log.Logger, archiverHost string) fetcher.Service {
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
	return &fetchersvc{logger: logger, archive: arc.Endpoint()}
}

// Fetch makes a GET request to the given URL and stores the results in the
// archiver service which must be running or the request fails
func (s *fetchersvc) Fetch(ctx context.Context, p *fetcher.FetchPayload) (*fetcher.FetchMedia, error) {
	// Make request to external endpoint
	resp, err := http.Get(p.URL)
	if err != nil {
		return nil, &fetcher.Error{}
	}

	// Read response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	now := int(time.Now().Unix())
	if err != nil {
		return nil, &fetcher.Error{
			ID:     strconv.Itoa(now),
			Code:   "bad_request",
			Detail: fmt.Sprintf("failed to decode response: %s", err),
		}
	}

	// Archive response using archiver service
	res, err := s.archive(ctx, &archiversvc.ArchivePayload{Status: resp.StatusCode, Body: string(body)})
	if err != nil {
		return nil, &fetcher.Error{
			ID:     strconv.Itoa(now),
			Code:   "bad_request",
			Detail: fmt.Sprintf("failed to decode response: %s", err),
		}
	}

	// Return response
	media := res.(*archiversvc.ArchiveMedia)
	return &fetcher.FetchMedia{
		ArchiveHref: media.Href,
		Status:      media.Status,
	}, nil
}

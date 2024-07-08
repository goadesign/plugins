package archiverapi

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kit/log"

	genarchiver "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/archiver"
	genserver "goa.design/plugins/v3/goakit/examples/fetcher/archiver/gen/http/archiver/server"
)

// archiver service example implementation.
// The example methods log the requests and return zero values.
type archiversvc struct {
	logger log.Logger
	db     *Archive
}

// NewArchiver returns the archiver service implementation.
func NewArchiver(logger log.Logger) genarchiver.Service {
	return &archiversvc{
		logger: logger,
		db:     &Archive{RWMutex: &sync.RWMutex{}},
	}
}

// Archive HTTP response
func (s *archiversvc) Archive(ctx context.Context, p *genarchiver.ArchivePayload) (*genarchiver.ArchiveMedia, error) {
	doc := &Document{Status: p.Status, Body: p.Body}
	s.db.Store(doc)
	return archiveMedia(doc), nil
}

// Read HTTP response from archive
func (s *archiversvc) Read(ctx context.Context, p *genarchiver.ReadPayload) (*genarchiver.ArchiveMedia, error) {
	doc := s.db.Read(p.ID)
	if doc == nil {
		return nil, genarchiver.MakeNotFound(fmt.Errorf("could not find document with ID %q", p.ID))
	}
	return archiveMedia(doc), nil
}

// archiveMedia converts a Document into a app.ArchiveMedia
func archiveMedia(doc *Document) *genarchiver.ArchiveMedia {
	return &genarchiver.ArchiveMedia{
		Href:   genserver.ReadArchiverPath(doc.ID),
		Status: doc.Status,
		Body:   doc.Body,
	}
}

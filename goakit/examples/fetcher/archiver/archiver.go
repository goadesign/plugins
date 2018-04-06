package archiver

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kit/kit/log"
	archiversvc "goa.design/plugins/goakit/examples/fetcher/archiver/gen/archiver"
	"goa.design/plugins/goakit/examples/fetcher/archiver/gen/http/archiver/server"
)

// archiver service example implementation.
type archiversvcsvc struct {
	logger log.Logger
	db     *Archive
}

// NewArchiver returns the archiver service implementation.
func NewArchiver(logger log.Logger) archiversvc.Service {
	return &archiversvcsvc{
		logger: logger,
		db:     &Archive{RWMutex: &sync.RWMutex{}},
	}
}

// Archive HTTP response
func (s *archiversvcsvc) Archive(ctx context.Context, p *archiversvc.ArchivePayload) (*archiversvc.ArchiveMedia, error) {
	doc := &Document{Status: p.Status, Body: p.Body}
	s.db.Store(doc)
	return archiveMedia(doc), nil
}

// Read HTTP response from archive
func (s *archiversvcsvc) Read(ctx context.Context, p *archiversvc.ReadPayload) (*archiversvc.ArchiveMedia, error) {
	doc := s.db.Read(p.ID)
	if doc == nil {
		return nil, archiversvc.MakeBadRequest(fmt.Errorf("could not find document with ID %q", p.ID))
	}
	return archiveMedia(doc), nil
}

// archiveMedia converts a Document into a app.ArchiveMedia
func archiveMedia(doc *Document) *archiversvc.ArchiveMedia {
	return &archiversvc.ArchiveMedia{
		Href:   server.ReadArchiverPath(doc.ID),
		Status: doc.Status,
		Body:   doc.Body,
	}
}

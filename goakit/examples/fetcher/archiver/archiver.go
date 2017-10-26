package archiver

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"goa.design/plugins/goakit/examples/client/archiver/gen/archiver"
	"goa.design/plugins/goakit/examples/client/archiver/gen/http/archiver/server"
)

type (
	// archiver service implementation.
	archiversvc struct {
		logger log.Logger
		db     *Archive
	}
)

// NewArchiver returns the archiver service implementation.
func NewArchiver(logger log.Logger) archiver.Service {
	return &archiversvc{
		logger: logger,
		db:     &Archive{RWMutex: &sync.RWMutex{}},
	}
}

// Archive HTTP response
func (s *archiversvc) Archive(ctx context.Context, p *archiver.ArchivePayload) (*archiver.ArchiveMedia, error) {
	doc := &Document{Status: p.Status, Body: p.Body}
	s.db.Store(doc)
	return archiveMedia(doc), nil
}

// Read HTTP response from archive
func (s *archiversvc) Read(ctx context.Context, p *archiver.ReadPayload) (*archiver.ArchiveMedia, error) {
	doc := s.db.Read(p.ID)
	if doc == nil {
		return nil, &archiver.Error{
			ID:     strconv.Itoa(int(time.Now().Unix())),
			Code:   "bad_request",
			Detail: fmt.Sprintf("could not find document with ID %q", p.ID),
		}
	}
	return archiveMedia(doc), nil
}

// archiveMedia converts a Document into a app.ArchiveMedia
func archiveMedia(doc *Document) *archiver.ArchiveMedia {
	return &archiver.ArchiveMedia{
		Href:   server.ReadArchiverPath(doc.ID),
		Status: doc.Status,
		Body:   doc.Body,
	}
}

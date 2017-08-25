package cellar

import (
	"context"

	"github.com/go-kit/kit/log"
	"goa.design/plugins/goakit/examples/cellar/gen/storage"
)

// storage service example implementation.
// The example methods log the requests and return zero values.
type storagesvc struct {
	logger log.Logger
}

// NewStorage returns the storage service implementation.
func NewStorage(logger log.Logger) storage.Service {
	return &storagesvc{logger}
}

// List all stored bottles
func (s *storagesvc) List(ctx context.Context) (storage.StoredBottleCollection, error) {
	var res storage.StoredBottleCollection
	s.logger.Log("msg", "storage.list")
	return res, nil
}

// Show bottle by ID
func (s *storagesvc) Show(ctx context.Context, p *storage.ShowPayload) (*storage.StoredBottle, error) {
	var res *storage.StoredBottle
	s.logger.Log("msg", "storage.show")
	return res, nil
}

// Add new bottle and return its ID.
func (s *storagesvc) Add(ctx context.Context, p *storage.Bottle) (string, error) {
	var res string
	s.logger.Log("msg", "storage.add")
	return res, nil
}

// Remove bottle from storage
func (s *storagesvc) Remove(ctx context.Context, p *storage.RemovePayload) error {
	s.logger.Log("msg", "storage.remove")
	return nil
}

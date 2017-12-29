package cellar

import (
	"context"
	"log"

	storage "goa.design/plugins/cors/examples/cellar/gen/storage"
)

// storage service example implementation.
// The example methods log the requests and return zero values.
type storageSvc struct {
	logger *log.Logger
}

// NewStorage returns the storage service implementation.
func NewStorage(logger *log.Logger) storage.Service {
	return &storageSvc{logger}
}

// List all stored bottles
func (s *storageSvc) List(ctx context.Context) (storage.StoredBottleCollection, error) {
	var res storage.StoredBottleCollection
	s.logger.Print("storage.list")
	return res, nil
}

// Show bottle by ID
func (s *storageSvc) Show(ctx context.Context, p *storage.ShowPayload) (*storage.StoredBottle, error) {
	res := &storage.StoredBottle{}
	s.logger.Print("storage.show")
	return res, nil
}

// Add new bottle and return its ID.
func (s *storageSvc) Add(ctx context.Context, p *storage.Bottle) (string, error) {
	var res string
	s.logger.Print("storage.add")
	return res, nil
}

// Remove bottle from storage
func (s *storageSvc) Remove(ctx context.Context, p *storage.RemovePayload) error {
	s.logger.Print("storage.remove")
	return nil
}

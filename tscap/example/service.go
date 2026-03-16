package example

import (
	"context"
	"log"

	tscapsvc "goa.design/plugins/v3/tscap/example/gen/tscap"
)

type tscapService struct {
	logger *log.Logger
}

// NewTscap returns the tscap service implementation.
func NewTscap(logger *log.Logger) tscapsvc.Service {
	return &tscapService{logger: logger}
}

// List items - requires read capability.
func (s *tscapService) List(ctx context.Context) ([]string, error) {
	s.logger.Print("list called")
	return []string{"item1", "item2", "item3"}, nil
}

// Create an item - requires write capability.
func (s *tscapService) Create(ctx context.Context, p *tscapsvc.CreatePayload) (string, error) {
	s.logger.Printf("create called with name: %s", p.Name)
	return "created: " + p.Name, nil
}

// Admin action - requires admin capability.
func (s *tscapService) Admin(ctx context.Context, p *tscapsvc.AdminPayload) error {
	s.logger.Printf("admin called with id: %s", p.ID)
	return nil
}

// Health check - no capability required.
func (s *tscapService) Health(ctx context.Context) (string, error) {
	return "ok", nil
}

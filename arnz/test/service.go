package test

import "context"

type Service struct{}

func (s *Service) Authenticated(ctx context.Context) error {
	return nil
}

func (s *Service) Authorized(ctx context.Context) error {
	return nil
}

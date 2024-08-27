package example

import "context"

type Service struct{}

func (s *Service) Create(ctx context.Context) error {
	return nil
}

func (s *Service) Read(ctx context.Context) error {
	return nil
}

func (s *Service) Update(ctx context.Context) error {
	return nil
}

func (s *Service) Delete(ctx context.Context) error {
	return nil
}

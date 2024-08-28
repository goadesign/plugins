package example

import (
	"context"

	genarnz "goa.design/plugins/v3/arnz/example/gen/arnz"
)

type Service struct{}

func (s *Service) Create(ctx context.Context) (res *genarnz.ResponseBody, err error) {
	return &genarnz.ResponseBody{Action: "created!"}, nil
}

func (s *Service) Read(ctx context.Context) (res *genarnz.ResponseBody, err error) {
	return &genarnz.ResponseBody{Action: "read!"}, nil
}

func (s *Service) Update(ctx context.Context) (res *genarnz.ResponseBody, err error) {
	return &genarnz.ResponseBody{Action: "updated!"}, nil
}

func (s *Service) Delete(ctx context.Context) (res *genarnz.ResponseBody, err error) {
	return &genarnz.ResponseBody{Action: "deleted!"}, nil
}

func (s *Service) Health(ctx context.Context) (res *genarnz.ResponseBody, err error) {
	return &genarnz.ResponseBody{Action: "healthy!"}, nil
}

package example

import (
	"context"

	genlike "goa.design/plugins/v3/arnz/example/gen/like"
	genmatch "goa.design/plugins/v3/arnz/example/gen/match"
)

type LikeService struct{}

func (s *LikeService) Create(ctx context.Context) (res *genlike.ResponseBody, err error) {
	return &genlike.ResponseBody{Action: "created!"}, nil
}

func (s *LikeService) Read(ctx context.Context) (res *genlike.ResponseBody, err error) {
	return &genlike.ResponseBody{Action: "read!"}, nil
}

func (s *LikeService) Update(ctx context.Context) (res *genlike.ResponseBody, err error) {
	return &genlike.ResponseBody{Action: "updated!"}, nil
}

func (s *LikeService) Delete(ctx context.Context) (res *genlike.ResponseBody, err error) {
	return &genlike.ResponseBody{Action: "deleted!"}, nil
}

type MatchService struct{}

func (s *MatchService) Create(ctx context.Context) (res *genmatch.ResponseBody, err error) {
	return &genmatch.ResponseBody{Action: "created!"}, nil
}

func (s *MatchService) Read(ctx context.Context) (res *genmatch.ResponseBody, err error) {
	return &genmatch.ResponseBody{Action: "read!"}, nil
}

func (s *MatchService) Update(ctx context.Context) (res *genmatch.ResponseBody, err error) {
	return &genmatch.ResponseBody{Action: "updated!"}, nil
}

func (s *MatchService) Delete(ctx context.Context) (res *genmatch.ResponseBody, err error) {
	return &genmatch.ResponseBody{Action: "deleted!"}, nil
}

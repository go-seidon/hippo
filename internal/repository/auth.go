package repository

import "context"

type AuthRepository interface {
	FindClient(ctx context.Context, p FindClientParam) (*FindClientResult, error)
}

type FindClientParam struct {
	ClientId string
}

type FindClientResult struct {
	ClientId     string
	ClientSecret string
}

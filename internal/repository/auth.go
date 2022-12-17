package repository

import (
	"context"
	"time"
)

type Auth interface {
	FindClient(ctx context.Context, p FindClientParam) (*FindClientResult, error)
}

type FindClientParam struct {
	Id       string
	ClientId string
}

type FindClientResult struct {
	Id           string
	ClientId     string
	ClientSecret string
	Name         string
	Type         string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

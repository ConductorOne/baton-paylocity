package client

import (
	"context"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

type ClientService interface {
	ListEmployees(ctx context.Context, options PageOptions) ([]*User, string, *v2.RateLimitDescription, error)
	ListPositionCodes(ctx context.Context, options PageOptions) ([]*Position, string, *v2.RateLimitDescription, error)
	GetUserById(ctx context.Context, userId string) (*User, *v2.RateLimitDescription, error)
}

type ClientServiceImpl struct {
	client *PaylocityClient
}

func (m *ClientServiceImpl) ListEmployees(ctx context.Context, options PageOptions) ([]*User, string, *v2.RateLimitDescription, error) {
	return m.client.ListEmployees(ctx, options)
}

func (m *ClientServiceImpl) ListPositionCodes(ctx context.Context, options PageOptions) ([]*Position, string, *v2.RateLimitDescription, error) {
	return m.client.ListPositionCodes(ctx, options)
}

func (m *ClientServiceImpl) GetUserById(ctx context.Context, userId string) (*User, *v2.RateLimitDescription, error) {
	return m.client.GetUserById(ctx, userId)
}

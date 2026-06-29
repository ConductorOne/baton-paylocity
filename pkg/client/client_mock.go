package client

import (
	"context"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

type MockService struct {
	ListEmployeesFunc     func(ctx context.Context, options PageOptions) ([]*User, string, *v2.RateLimitDescription, error)
	ListPositionCodesFunc func(ctx context.Context, options PageOptions) ([]*Position, string, *v2.RateLimitDescription, error)
}

func (m *MockService) ListEmployees(ctx context.Context, options PageOptions) ([]*User, string, *v2.RateLimitDescription, error) {
	return m.ListEmployeesFunc(ctx, options)
}

func (m *MockService) ListPositionCodes(ctx context.Context, options PageOptions) ([]*Position, string, *v2.RateLimitDescription, error) {
	return m.ListPositionCodesFunc(ctx, options)
}

package client

import (
	"context"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

type ClientService interface {
	ListEmployees(ctx context.Context, options PageOptions) ([]*User, string, *v2.RateLimitDescription, error)
	ListPositionCodes(ctx context.Context, options PageOptions) ([]*Position, string, *v2.RateLimitDescription, error)
}

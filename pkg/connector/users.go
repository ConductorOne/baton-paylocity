package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-paylocity/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

type userBuilder struct {
	c *client.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var annot *annotations.Annotations
	var nextPage string

	offset, limit, err := parsePaginationToken(pToken, 20)
	if err != nil {
		return nil, "", nil, fmt.Errorf("cannot parse pagination token, error: %w", err)
	}

	// FIXME(shackra): change ListEmployees to use the `offset` and `limit` args
	resp, rl, err := o.c.ListEmployees(ctx, offset, limit)
	annot = annot.WithRateLimiting(rl)
	if err != nil {
		return nil, "", *annot, fmt.Errorf("cannot list employees, error: %w", err)
	}
	nextPage = makeNextPage(offset, limit, resp.TotalCount)

	users, err := employees2users(resp.Employees, parentResourceID)
	if err != nil {
		return nil, "", nil, err
	}

	return users, nextPage, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder() *userBuilder {
	return &userBuilder{}
}

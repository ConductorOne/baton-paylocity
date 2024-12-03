package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-paylocity/pkg/client"
	"github.com/conductorone/baton-paylocity/pkg/models"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	pageSize = 100
)

type roleBuilder struct {
	c *client.Client
}

func (o *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

func newRoleBuilder(client *client.Client) *roleBuilder {
	return &roleBuilder{c: client}
}

func (o *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var annot *annotations.Annotations
	var nextPage string

	offset, limit, err := parsePaginationToken(pToken, pageSize)
	if err != nil {
		return nil, "", nil, fmt.Errorf("cannot parse pagination token, error: %w", err)
	}
	resp, rl, err := o.c.ListPositionCodes(ctx, offset, limit)
	annot = annot.WithRateLimiting(rl)
	if err != nil {
		return nil, "", *annot, fmt.Errorf("cannot list position codes, error: %w", err)
	}

	nextPage = makeNextPageStr(offset, limit, resp.Total)

	roles, err := positionCodes2roles(resp.Data, parentResourceID)
	if err != nil {
		return nil, "", nil, err
	}
	return roles, nextPage, nil, nil
}

func (o *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	options := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Role member", resource.DisplayName)),
	}
	return []*v2.Entitlement{ent.NewAssignmentEntitlement(resource, "member", options...)}, "", nil, nil
}

func (o *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var employees []models.Employee
	{ // scope listing employees
		page := 0
		for page != -1 {
			resp, _, err := o.c.ListEmployees(ctx, page, pageSize)
			if err != nil {
				return nil, "", nil, fmt.Errorf("cannot list all employees for grants, error: %w", err)
			}
			employees = append(employees, resp.Employees...)

			page = makeNextPage(page, pageSize, resp.TotalCount)
		}
	}
	var grants []*v2.Grant
	for _, e := range employees {
		positionIDOfEmployee, err := resourceSdk.NewResourceID(roleResourceType, e.Position.PositionCode)
		if err != nil {
			return nil, "", nil, fmt.Errorf("cannot make resource's ID, error: %w", err)
		}
		if resource.Id == positionIDOfEmployee {
			grant := grant.NewGrant(
				resource,
				"member",
				positionIDOfEmployee,
			)

			grants = append(grants, grant)
			break
		}
	}

	return grants, "", nil, nil
}

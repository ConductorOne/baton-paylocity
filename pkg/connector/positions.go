package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/conductorone/baton-paylocity/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	PermissionMember = "member"
)

type PositionBuilder struct {
	client client.ClientService
}

func (o *PositionBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return positionResourceType
}

func (o *PositionBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	outputAnnotations := annotations.New()
	options := getPageOptions(pToken, client.ItemsPerPage)

	positions, nextPage, rateLimitDesc, err := o.client.ListPositionCodes(ctx, options)
	if rateLimitDesc != nil {
		outputAnnotations.WithRateLimiting(rateLimitDesc)
	}
	if err != nil {
		return nil, "", outputAnnotations, fmt.Errorf("listing position codes failed: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(positions))
	for _, position := range positions {
		positionResource, err := parsePositionToResource(position, parentResourceID)
		if err != nil {
			return nil, "", outputAnnotations, err
		}
		resources = append(resources, positionResource)
	}

	return resources, nextPage, outputAnnotations, nil
}

func (o *PositionBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	opts := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(fmt.Sprintf("%s Member", resource.DisplayName)),
		entitlement.WithDescription(fmt.Sprintf("Is a member of the %s position in Paylocity", resource.DisplayName)),
	}

	ent := entitlement.NewAssignmentEntitlement(resource, PermissionMember, opts...)

	return []*v2.Entitlement{ent}, "", nil, nil
}

// Grants are intentionally not implemented here for performance reasons.
// Therefore, implementing grants here would require fetching all users for every single position and filtering them in memory.
// Instead, grants are created on the userBuilder, as the employee data already
// contains the necessary 'positionCode' to build the relationship.
func (o *PositionBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parsePositionToResource(position *client.Position, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"title":          position.Title,
		"code":           position.Code,
		"effective_date": position.EffectiveDate.Format(time.RFC3339),
		"closed_date":    position.ClosedDate.Format(time.RFC3339),
	}

	positionTraitOptions := []resource.RoleTraitOption{
		resource.WithRoleProfile(profile),
	}

	positionResource, err := resource.NewRoleResource(
		position.Title,
		positionResourceType,
		position.Code,
		positionTraitOptions,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return positionResource, nil
}

func newPositionsBuilder(service client.ClientService) *PositionBuilder {
	return &PositionBuilder{
		client: service,
	}
}

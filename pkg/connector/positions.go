package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/conductorone/baton-paylocity/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

var _ connectorbuilder.ResourceSyncerV2 = (*PositionBuilder)(nil)

const (
	PermissionMember = "member"
)

type PositionBuilder struct {
	client client.ClientService
}

func (o *PositionBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return positionResourceType
}

func (o *PositionBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	outputAnnotations := annotations.New()
	options := getPageOptions(&pToken.PageToken, client.ItemsPerPage)

	positions, nextPage, rateLimitDesc, err := o.client.ListPositionCodes(ctx, options)
	if rateLimitDesc != nil {
		outputAnnotations.WithRateLimiting(rateLimitDesc)
	}
	if err != nil {
		return nil, &resource.SyncOpResults{Annotations: outputAnnotations}, fmt.Errorf("listing position codes failed: %w", err)
	}

	resources := make([]*v2.Resource, 0, len(positions))
	for _, position := range positions {
		positionResource, err := parsePositionToResource(position, parentResourceID)
		if err != nil {
			return nil, &resource.SyncOpResults{Annotations: outputAnnotations}, err
		}
		resources = append(resources, positionResource)
	}

	return resources, &resource.SyncOpResults{NextPageToken: nextPage, Annotations: outputAnnotations}, nil
}

func (o *PositionBuilder) Entitlements(_ context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	opts := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(fmt.Sprintf("%s Member", res.DisplayName)),
		entitlement.WithDescription(fmt.Sprintf("Is a member of the %s position in Paylocity", res.DisplayName)),
	}

	ent := entitlement.NewAssignmentEntitlement(res, PermissionMember, opts...)

	return []*v2.Entitlement{ent}, &resource.SyncOpResults{}, nil
}

// Grants are intentionally not implemented here for performance reasons.
// Therefore, implementing grants here would require fetching all users for every single position and filtering them in memory.
// Instead, grants are created on the userBuilder, as the employee data already
// contains the necessary 'positionCode' to build the relationship.
func (o *PositionBuilder) Grants(_ context.Context, _ *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	return nil, nil, nil
}

func parsePositionToResource(position *client.Position, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"title":          position.Title,
		"code":           position.Code,
		"effective_date": position.EffectiveDate.Format(time.RFC3339),
		"closed_date":    position.ClosedDate.Format(time.RFC3339),
	}

	positionTraitOptions := []resource.RoleTraitOption{}

	positionResource, err := resource.NewRoleResource(
		position.Title,
		positionResourceType,
		position.Code,
		positionTraitOptions,
		resource.WithParentResourceID(parentResourceID),
		resource.WithResourceProfile(profile),
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

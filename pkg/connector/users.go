package connector

import (
	"context"
	"fmt"
	"strings"

	"github.com/conductorone/baton-paylocity/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	service client.ClientService
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	options := getPageOptions(pToken, client.ItemsPerPage)

	outputAnnotations := annotations.New()
	employees, nextPageToken, rateLimitDesc, err := o.service.ListEmployees(ctx, options)
	if rateLimitDesc != nil {
		outputAnnotations.WithRateLimiting(rateLimitDesc)
	}
	if err != nil {
		return nil, "", outputAnnotations, fmt.Errorf("listing employees failed: %w", err)
	}

	var resources []*v2.Resource
	for _, employee := range employees {
		userResource, err := parseUserToResource(employee, parentResourceID)
		if err != nil {
			return nil, "", outputAnnotations, err
		}
		resources = append(resources, userResource)
	}

	return resources, nextPageToken, outputAnnotations, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// The Grants function for positions is implemented here on the userBuilder for performance reasons.
// This allows assigning grants directly while iterating through users, as the Paylocity API
// already includes the 'positionCode' in the employee list, avoiding the need to iterate
// over all positions to resolve each user's assignment individually.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	userId := resource.Id.Resource
	outputAnnotations := annotations.New()

	user, rateLimitDesc, err := o.service.GetUserById(ctx, userId)
	if rateLimitDesc != nil {
		outputAnnotations.WithRateLimiting(rateLimitDesc)
	}
	if err != nil {
		return nil, "", outputAnnotations, fmt.Errorf("failed to get user %s for grants: %w", userId, err)
	}

	var grants []*v2.Grant
	if user.Position.PositionCode != "" {
		positionResource := &v2.Resource{
			Id: &v2.ResourceId{
				ResourceType: positionResourceType.Id,
				Resource:     user.Position.PositionCode,
			},
		}
		grants = append(grants, grant.NewGrant(positionResource, PermissionMember, resource.Id))
	}

	return grants, "", outputAnnotations, nil
}

func parseUserToResource(user *client.User, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"user_id":       user.ID,
		"display_name":  user.DisplayName,
		"first_name":    user.Info.FirstName,
		"last_name":     user.Info.LastName,
		"job_title":     user.Info.JobTitle,
		"employee_type": user.Position.EmployeeType,
		"position_code": user.Position.PositionCode,
		"department":    user.Position.Department,
	}

	var status v2.UserTrait_Status_Status
	normalizedStatus := strings.ToLower(user.Status)
	switch normalizedStatus {
	case "active":
		status = v2.UserTrait_Status_STATUS_ENABLED
	default:
		status = v2.UserTrait_Status_STATUS_DISABLED
	}

	userTraitOptions := []resource.UserTraitOption{
		resource.WithEmail(user.Info.Email, true),
		resource.WithUserLogin(user.Info.Email),
		resource.WithStatus(status),
		resource.WithUserProfile(profile),
	}

	userResource, err := resource.NewUserResource(
		user.DisplayName,
		userResourceType,
		user.ID,
		userTraitOptions,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return userResource, nil
}

func newUserBuilder(service client.ClientService) *userBuilder {
	return &userBuilder{
		service: service,
	}
}

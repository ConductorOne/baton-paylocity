package connector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/conductorone/baton-paylocity/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/test"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newTestUserBuilder() (*userBuilder, *client.MockService) {
	mockService := &client.MockService{}

	builder := newUserBuilder(mockService)

	return builder, mockService
}
func TestUserList(t *testing.T) {
	ctx := context.Background()

	mockUser1 := &client.User{
		ID:          "101",
		DisplayName: "Ana Gomez",
		Status:      "Active",
		Info:        client.InfoPayload{Email: "ana.gomez@example.com"},
	}

	t.Run("should list users successfully", func(t *testing.T) {
		userBuilder, mockClientService := newTestUserBuilder()
		mockClientService.ListEmployeesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.User, string, *v2.RateLimitDescription, error) {
			return []*client.User{mockUser1}, "", nil, nil
		}

		resources, opResult, err := userBuilder.List(ctx, nil, rs.SyncOpAttrs{})

		require.NoError(t, err)
		require.Len(t, resources, 1)
		require.Empty(t, opResult.NextPageToken)
		test.AssertNoRatelimitAnnotations(t, opResult.Annotations)
		require.Equal(t, "Ana Gomez", resources[0].DisplayName)
	})

	t.Run("should paginate correctly", func(t *testing.T) {
		// Arrange
		userBuilder, mockClientService := newTestUserBuilder()
		mockUser2 := &client.User{ID: "102", DisplayName: "Juan Perez"}

		callCount := 0
		mockClientService.ListEmployeesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.User, string, *v2.RateLimitDescription, error) {
			callCount++
			if options.PageToken == "" {
				return []*client.User{mockUser1}, "token_pagina_2", nil, nil
			}
			if options.PageToken == "token_pagina_2" {
				return []*client.User{mockUser2}, "", nil, nil
			}
			return nil, "", nil, fmt.Errorf("unexpected page token")
		}

		resources1, opResult, err1 := userBuilder.List(ctx, nil, rs.SyncOpAttrs{})
		require.NoError(t, err1)
		require.Len(t, resources1, 1)
		require.Equal(t, "token_pagina_2", opResult.NextPageToken)

		resources2, opResult2, err2 := userBuilder.List(ctx, nil, rs.SyncOpAttrs{
			PageToken: pagination.Token{Token: opResult.NextPageToken},
		})
		require.NoError(t, err2)
		require.Len(t, resources2, 1)
		require.Empty(t, opResult2.NextPageToken)

		// Assert
		require.Equal(t, 2, callCount, "ListEmployees should have been called twice")
		require.Equal(t, "Juan Perez", resources2[0].DisplayName)
	})

	t.Run("should handle rate limit errors", func(t *testing.T) {
		// Arrange
		userBuilder, mockClientService := newTestUserBuilder()
		expectedReset := time.Now().Add(10 * time.Second)
		mockClientService.ListEmployeesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.User, string, *v2.RateLimitDescription, error) {
			return nil, "", &v2.RateLimitDescription{ResetAt: timestamppb.New(expectedReset)}, fmt.Errorf("rate limited")
		}

		_, optResult, err := userBuilder.List(ctx, nil, rs.SyncOpAttrs{})

		require.Error(t, err)
		require.NotNil(t, optResult.Annotations)

		rateLimitAnn := &v2.RateLimitDescription{}
		ok, pickErr := optResult.Annotations.Pick(rateLimitAnn)
		require.NoError(t, pickErr)
		require.True(t, ok, "rate limit annotation should be present")
		require.Equal(t, expectedReset.Unix(), rateLimitAnn.ResetAt.Seconds)
	})
}

func TestUserGrants(t *testing.T) {
	ctx := context.Background()
	userResource := &v2.Resource{
		Id: &v2.ResourceId{ResourceType: userResourceType.Id, Resource: "101"},
	}
	mockUser := &client.User{
		ID:       "101",
		Position: client.PositionPayload{PositionCode: "DEV-01"},
	}

	t.Run("should return grant for user with position", func(t *testing.T) {
		// Arrange
		userBuilder, mockClientService := newTestUserBuilder()
		mockClientService.GetUserByIdFunc = func(ctx context.Context, employeeID string) (*client.User, *v2.RateLimitDescription, error) {
			return mockUser, nil, nil
		}

		// Act
		grants, _, err := userBuilder.Grants(ctx, userResource, rs.SyncOpAttrs{})

		// Assert
		require.NoError(t, err)
		require.Len(t, grants, 1)
		require.Equal(t, "DEV-01", grants[0].Entitlement.Resource.Id.Resource)
		require.Equal(t, "101", grants[0].Principal.Id.Resource)
	})

	t.Run("should return error if GetUserById fails", func(t *testing.T) {
		// Arrange
		userBuilder, mockClientService := newTestUserBuilder()
		mockClientService.GetUserByIdFunc = func(ctx context.Context, employeeID string) (*client.User, *v2.RateLimitDescription, error) {
			return nil, nil, fmt.Errorf("user not found")
		}

		_, _, err := userBuilder.Grants(ctx, userResource, rs.SyncOpAttrs{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get user 101 for grants: user not found")
	})

	t.Run("should handle rate limit on GetUserById", func(t *testing.T) {
		// Arrange
		userBuilder, mockClientService := newTestUserBuilder()
		expectedReset := time.Now().Add(15 * time.Second)
		mockClientService.GetUserByIdFunc = func(ctx context.Context, employeeID string) (*client.User, *v2.RateLimitDescription, error) {
			return nil, &v2.RateLimitDescription{ResetAt: timestamppb.New(expectedReset)}, fmt.Errorf("rate limited")
		}

		_, opResult, err := userBuilder.Grants(ctx, userResource, rs.SyncOpAttrs{})

		require.Error(t, err)
		require.NotNil(t, opResult.Annotations)

		rateLimitAnn := &v2.RateLimitDescription{}
		ok, pickErr := opResult.Annotations.Pick(rateLimitAnn)
		require.NoError(t, pickErr)
		require.True(t, ok, "rate limit annotation should be present")
		require.Equal(t, expectedReset.Unix(), rateLimitAnn.ResetAt.Seconds)
	})
}

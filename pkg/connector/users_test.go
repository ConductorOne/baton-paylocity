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
		ID:            "101",
		Info:          client.InfoPayload{DisplayName: "Ana Gomez", Email: "ana.gomez@example.com"},
		CurrentStatus: client.StatusPayload{Status: "Active"},
	}

	t.Run("should list users successfully", func(t *testing.T) {
		userBuilder, mockClientService := newTestUserBuilder()
		mockClientService.ListEmployeesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.User, string, *v2.RateLimitDescription, error) {
			return []*client.User{mockUser1}, "", nil, nil
		}

		resources, nextPageToken, annotations, err := userBuilder.List(ctx, nil, &pagination.Token{})

		require.NoError(t, err)
		require.Len(t, resources, 1)
		require.Empty(t, nextPageToken)
		test.AssertNoRatelimitAnnotations(t, annotations)
		require.Equal(t, "Ana Gomez", resources[0].DisplayName)
	})

	t.Run("should paginate correctly", func(t *testing.T) {
		// Arrange
		userBuilder, mockClientService := newTestUserBuilder()
		mockUser2 := &client.User{ID: "102", Info: client.InfoPayload{DisplayName: "Juan Perez"}}

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

		resources1, nextPageToken1, _, err1 := userBuilder.List(ctx, nil, &pagination.Token{})
		require.NoError(t, err1)
		require.Len(t, resources1, 1)
		require.Equal(t, "token_pagina_2", nextPageToken1)

		resources2, nextPageToken2, _, err2 := userBuilder.List(ctx, nil, &pagination.Token{Token: nextPageToken1})
		require.NoError(t, err2)
		require.Len(t, resources2, 1)
		require.Empty(t, nextPageToken2)

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

		_, _, annotations, err := userBuilder.List(ctx, nil, &pagination.Token{})

		require.Error(t, err)
		require.NotNil(t, annotations)

		rateLimitAnn := &v2.RateLimitDescription{}
		ok, pickErr := annotations.Pick(rateLimitAnn)
		require.NoError(t, pickErr)
		require.True(t, ok, "rate limit annotation should be present")
		require.Equal(t, expectedReset.Unix(), rateLimitAnn.ResetAt.Seconds)
	})
}

func TestUserGrants(t *testing.T) {
	ctx := context.Background()

	t.Run("should return grant for user with position", func(t *testing.T) {
		userBuilder, _ := newTestUserBuilder()
		userResource, err := parseUserToResource(&client.User{
			ID:       "101",
			Info:     client.InfoPayload{DisplayName: "Ana Gomez"},
			Position: client.PositionPayload{PositionCode: "DEV-01"},
		}, nil)
		require.NoError(t, err)

		grants, _, _, err := userBuilder.Grants(ctx, userResource, &pagination.Token{})

		require.NoError(t, err)
		require.Len(t, grants, 1)
		require.Equal(t, "DEV-01", grants[0].Entitlement.Resource.Id.Resource)
		require.Equal(t, "101", grants[0].Principal.Id.Resource)
	})

	t.Run("should return no grants for user without position", func(t *testing.T) {
		userBuilder, _ := newTestUserBuilder()
		userResource, err := parseUserToResource(&client.User{
			ID:   "P1",
			Info: client.InfoPayload{DisplayName: "Administrator"},
		}, nil)
		require.NoError(t, err)

		grants, _, _, err := userBuilder.Grants(ctx, userResource, &pagination.Token{})

		require.NoError(t, err)
		require.Empty(t, grants)
	})
}

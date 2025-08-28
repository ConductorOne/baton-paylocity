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

func newTestPositionBuilder() (*PositionBuilder, *client.MockService) {
	mockService := &client.MockService{}
	builder := newPositionsBuilder(mockService)
	return builder, mockService
}

func TestPositionList(t *testing.T) {
	ctx := context.Background()
	mockPosition1 := &client.Position{Code: "DEV-01", Title: "Software Developer"}

	t.Run("should list positions successfully", func(t *testing.T) {
		positionBuilder, mockClientService := newTestPositionBuilder()
		mockClientService.ListPositionCodesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.Position, string, *v2.RateLimitDescription, error) {
			return []*client.Position{mockPosition1}, "", nil, nil
		}

		resources, nextPageToken, annotations, err := positionBuilder.List(ctx, nil, &pagination.Token{})

		require.NoError(t, err)
		require.Len(t, resources, 1)
		require.Empty(t, nextPageToken)
		test.AssertNoRatelimitAnnotations(t, annotations)
		require.Equal(t, "Software Developer", resources[0].DisplayName)
	})

	t.Run("should paginate correctly", func(t *testing.T) {
		positionBuilder, mockClientService := newTestPositionBuilder()
		mockPosition2 := &client.Position{Code: "MGR-01", Title: "Engineering Manager"}

		callCount := 0
		mockClientService.ListPositionCodesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.Position, string, *v2.RateLimitDescription, error) {
			callCount++
			if options.PageToken == "" {
				return []*client.Position{mockPosition1}, "50", nil, nil
			}
			if options.PageToken == "50" {
				return []*client.Position{mockPosition2}, "", nil, nil
			}
			return nil, "", nil, fmt.Errorf("unexpected page token (offset)")
		}

		resources1, nextPageToken1, _, err1 := positionBuilder.List(ctx, nil, &pagination.Token{Size: 50})
		require.NoError(t, err1)
		require.Len(t, resources1, 1)
		require.Equal(t, "50", nextPageToken1)

		resources2, nextPageToken2, _, err2 := positionBuilder.List(ctx, nil, &pagination.Token{Token: nextPageToken1, Size: 50})
		require.NoError(t, err2)
		require.Len(t, resources2, 1)
		require.Empty(t, nextPageToken2)

		require.Equal(t, 2, callCount, "ListPositionCodes should have been called twice")
	})

	t.Run("should handle rate limit errors", func(t *testing.T) {
		positionBuilder, mockClientService := newTestPositionBuilder()
		expectedReset := time.Now().Add(10 * time.Second)
		mockClientService.ListPositionCodesFunc = func(ctx context.Context, options client.PageOptions) ([]*client.Position, string, *v2.RateLimitDescription, error) {
			return nil, "", &v2.RateLimitDescription{ResetAt: timestamppb.New(expectedReset)}, fmt.Errorf("rate limited")
		}

		_, _, annotations, err := positionBuilder.List(ctx, nil, &pagination.Token{})

		require.Error(t, err)
		require.NotNil(t, annotations)
	})
}

func TestPositionEntitlements(t *testing.T) {
	ctx := context.Background()

	t.Run("should return member entitlement for a position", func(t *testing.T) {
		positionBuilder, _ := newTestPositionBuilder()
		positionResource := &v2.Resource{
			Id:          &v2.ResourceId{ResourceType: positionResourceType.Id, Resource: "DEV-01"},
			DisplayName: "Software Developer",
		}

		entitlements, _, _, err := positionBuilder.Entitlements(ctx, positionResource, &pagination.Token{})

		require.NoError(t, err)
		require.Len(t, entitlements, 1)
		ent := entitlements[0]
		require.Equal(t, "member", ent.Slug)
		require.Equal(t, "Software Developer Member", ent.DisplayName)
		require.NotNil(t, ent.Resource.Id)
		require.Equal(t, positionResourceType.Id, ent.Resource.Id.ResourceType)
		require.Equal(t, "DEV-01", ent.Resource.Id.Resource)
	})
}

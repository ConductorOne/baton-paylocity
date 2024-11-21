package client

import (
	"errors"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func withBearerToken(token string) uhttp.RequestOption {
	return uhttp.WithHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}

func shouldRefresh(err error) (bool, error) {
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return true, nil
		}
	}
	return false, err
}

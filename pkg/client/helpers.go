package client

import (
	"errors"
	"fmt"
	"strconv"

	liburl "net/url"

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

func urlAddQuery(url string, params map[string]interface{}) (string, error) {
	p := liburl.Values{}
	for k, v := range params {
		switch value := v.(type) {
		case string:
			p.Add(k, value)
		case int:
			p.Add(k, strconv.Itoa(value))
		case bool:
			p.Add(k, strconv.FormatBool(value))
		default:
			continue
		}
	}

	parsed, err := liburl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("cannot parse URL, error: %w", err)
	}

	parsed.RawQuery = p.Encode()

	return parsed.String(), nil
}

package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-paylocity/pkg/models"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

type Client struct {
	baseURL      string
	companyID    string
	bearerToken  string
	clientID     string
	clientSecret string
	httpClient   uhttp.BaseHttpClient
}

var ErrUnauthorized = errors.New("Unauthorized")

// getBearerToken fetch a bearer token from the server.
func (c *Client) getBearerToken(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Add("client_id", c.clientID)
	form.Add("client_secret", c.clientSecret)
	form.Add("grant_type", "client_credentials")

	endpoint, err := url.JoinPath(c.baseURL, "/public/security/v1/token")
	if err != nil {
		return "", fmt.Errorf("cannot make endpoint URL, error: %w", err)
	}
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("cannot parse endpoint URL, error: %w", err)
	}

	req, err := c.httpClient.NewRequest(ctx, http.MethodPost, endpointURL, uhttp.WithFormBody(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to make request, error: %w", err)
	}

	var bodyResponse models.AuthResponse
	resp, err := c.httpClient.Do(req,
		uhttp.WithJSONResponse(&bodyResponse),
	)
	if err != nil {
		return "", fmt.Errorf("request failed to complete, error: %w", err)
	}
	defer resp.Body.Close()

	// NOTE(shackra): ensure that the server responds with 403 and not 401
	if resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("Server responded with 403 Forbidden, check your credentials are valid and try again")
	}

	return bodyResponse.AccessToken, nil
}

func (c *Client) get(ctx context.Context, endpoint string, target any, options ...uhttp.RequestOption) error {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("cannot parse endpoint URL, error: %w", err)
	}
	request, err := c.httpClient.NewRequest(ctx, http.MethodGet, parsedURL, options...)
	if err != nil {
		return fmt.Errorf("cannot create request, error: %w", err)
	}

	var ratelimitData v2.RateLimitDescription
	doOptions := []uhttp.DoOption{
		uhttp.WithRatelimitData(&ratelimitData),
		uhttp.WithJSONResponse(target),
	}

	resp, err := c.httpClient.Do(request, doOptions...)
	if err != nil {
		return fmt.Errorf("request failed, error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	return nil
}

// executePreparedRequest calls a function that is in charge of
// configuring the request and execute it, if for some reason the
// token needs to be refreshed this function takes care of that and
// tries a second time.
func (c *Client) executePreparedRequest(ctx context.Context, f func(string) (bool, error)) error {
	shouldRefresh, err := f(c.bearerToken)
	if !shouldRefresh {
		if err != nil {
			return err
		}
		return nil
	}

	token, err := c.getBearerToken(ctx)
	if err != nil {
		return fmt.Errorf("cannot refresh bearer token, error: %w", err)
	}

	c.bearerToken = token
	// call the same function again, this function should use the
	// new token
	shouldRefresh, err = f(c.bearerToken)

	if !shouldRefresh {
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("request failed due to authorization, check your credentials and try again")
}

func (c *Client) ListEmployees(ctx context.Context) ([]models.Employee, error) {
	joinedURL, err := url.JoinPath(c.baseURL, "/coreHr/v1/companies/", c.companyID, "/employees")
	if err != nil {
		return nil, fmt.Errorf("cannot make endpoint URL, error: %w", err)
	}

	var target models.EmployeesResponse
	err = c.executePreparedRequest(ctx, func(t string) (bool, error) {
		inerr := c.get(ctx, joinedURL, &target, uhttp.WithAcceptJSONHeader(), withBearerToken(t))
		return shouldRefresh(inerr)
	})
	if err != nil {
		return nil, fmt.Errorf("request failed, error: %w", err)
	}

	return target.Employees, nil
}

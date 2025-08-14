package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type PaylocityClient struct {
	httpClient *uhttp.BaseHttpClient
	baseURL    string
	companyID  string
}

func New(ctx context.Context, clientID, clientSecret, baseURL, companyID string) (*PaylocityClient, error) {
	cfg := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("%s/public/security/v1/token", baseURL),
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	tokenSource := cfg.TokenSource(ctx)

	oauthClient := oauth2.NewClient(ctx, tokenSource)
	cli, err := uhttp.NewBaseHttpClientWithContext(ctx, oauthClient)
	if err != nil {
		return nil, err
	}

	client := &PaylocityClient{
		httpClient: cli,
		baseURL:    baseURL,
		companyID:  companyID,
	}

	return client, nil
}

func (c *PaylocityClient) ListPositionCodes(ctx context.Context, options PageOptions) ([]*Position, string, *v2.RateLimitDescription, error) {
	var res []*Position

	safeLimit := options.Limit
	if safeLimit <= 0 {
		safeLimit = ItemsPerPage
	}

	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, "", nil, fmt.Errorf("invalid base URL: %w", err)
	}
	endpoint := baseURL.JoinPath("/apiHub/positionManagement/v1/companies", c.companyID, "positions")

	opts := []ReqOpt{
		withLimitParam(safeLimit),
		withOffset(options.PageToken),
		withTotalCount(),
	}
	for _, opt := range opts {
		opt(endpoint)
	}

	header, rl, err := c.doRequest(ctx, http.MethodGet, endpoint, &res, nil)
	if err != nil {
		return nil, "", rl, err
	}

	total, _ := strconv.Atoi(header.Get("X-Pcty-Total-Count"))
	offset, _ := strconv.Atoi(options.PageToken)
	nextOffsetStr := getNextPageToken(offset, safeLimit, total)

	return res, nextOffsetStr, rl, nil
}

func (c *PaylocityClient) ListEmployees(ctx context.Context, options PageOptions) ([]*User, string, *v2.RateLimitDescription, error) {
	var res EmployeesResponse

	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, "", nil, fmt.Errorf("invalid base URL: %w", err)
	}
	endpoint := baseURL.JoinPath("/coreHr/v1/companies", c.companyID, "employees")

	opts := []ReqOpt{
		withLimitParam(options.Limit),
		withNextToken(options.PageToken),
	}
	for _, opt := range opts {
		opt(endpoint)
	}

	_, rl, err := c.doRequest(ctx, http.MethodGet, endpoint, &res, nil)
	if err != nil {
		return nil, "", rl, err
	}

	return res.Employees, res.NextToken, rl, nil
}

func (c *PaylocityClient) GetUserById(ctx context.Context, userId string) (*User, *v2.RateLimitDescription, error) {
	var user User

	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid base URL: %w", err)
	}

	endpoint := baseURL.JoinPath("/coreHr/v1/companies", c.companyID, "employees", userId)

	opts := []ReqOpt{
		withIncludes("info", "position", "status"),
	}
	for _, opt := range opts {
		opt(endpoint)
	}

	_, rl, err := c.doRequest(ctx, http.MethodGet, endpoint, &user, nil)
	if err != nil {
		return nil, rl, err
	}

	return &user, rl, nil
}

func (c *PaylocityClient) doRequest(ctx context.Context, method string, url *url.URL, target interface{}, body interface{}) (*http.Header, *v2.RateLimitDescription, error) {
	requestOptions := []uhttp.RequestOption{
		uhttp.WithAcceptJSONHeader(),
	}
	if body != nil {
		requestOptions = append(requestOptions, uhttp.WithContentTypeJSONHeader(), uhttp.WithJSONBody(body))
	}

	request, err := c.httpClient.NewRequest(ctx, method, url, requestOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var rateLimitData v2.RateLimitDescription
	var pError PaylocityErrorResponse
	doOptions := []uhttp.DoOption{
		uhttp.WithRatelimitData(&rateLimitData),
		uhttp.WithErrorResponse(&pError),
	}
	if target != nil {
		doOptions = append(doOptions, uhttp.WithJSONResponse(target))
	}

	response, err := c.httpClient.Do(request, doOptions...)
	if err != nil {
		if len(pError) > 0 {
			return nil, nil, fmt.Errorf("paylocity API error: %s", pError.Message())
		}
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}()

	return &response.Header, &rateLimitData, nil
}

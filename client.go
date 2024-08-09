package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	utils "github.com/adamchol/go-anthropic-sdk/internal"
)

// Anthropic API Client for making requests
type Client struct {
	config ClientConfig

	requestBuilder utils.RequestBuilder
}

// NewClientWithConfig creates an Anthropic API client with specified configuration
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config:         config,
		requestBuilder: utils.NewRequestBuilder(),
	}
}

// NewClient creates an Anthropic API client with API key
func NewClient(apiKey string) *Client {
	return NewClientWithConfig(DefaultConfig(apiKey))
}

type requestOptions struct {
	body   any
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (*http.Request, error) {
	args := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(args)
	}
	req, err := c.requestBuilder.Build(ctx, method, url, args.body, args.header)
	if err != nil {
		return nil, err
	}
	c.setCommonHeaders(req)
	return req, nil
}

func (c *Client) sendRequest(request *http.Request, v any) error {
	resp, err := c.config.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if isFailureStatusCode(resp) {
		// TODO: handle errors properly
		return handleErrorResponse(resp)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", string(c.config.APIVersion))
	req.Header.Set("x-api-key", c.config.authToken)
}

func (c *Client) fullURL(suffix string) string {
	return fmt.Sprintf("%s%s", c.config.BaseUrl, suffix)
}

func isFailureStatusCode(response *http.Response) bool {
	return response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest
}

type ErrorResponse struct {
	Type  string    `json:"type"`
	Error *APIError `json:"error"`
}

type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func handleErrorResponse(resp *http.Response) error {
	var errResp ErrorResponse

	err := json.NewDecoder(resp.Body).Decode(&errResp)
	if err != nil || errResp.Error == nil {
		return err
	}

	return errors.New(errResp.Error.Message)
}

package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	const mockKey = "mock-key"
	client := NewClient(mockKey)
	if client.config.authToken != mockKey {
		t.Errorf("Client does not contain proper auth key")
	}
}

func TestNewRequest(t *testing.T) {
	client := NewClient("mock-key")
	mockBody := struct {
		Text string   `json:"text"`
		List []string `json:"list"`
	}{
		"some text",
		[]string{
			"mock",
			"text",
		},
	}
	req, err := client.newRequest(context.Background(), http.MethodPost, "https://mockurl.com", withBody(mockBody))
	assert.NoError(t, err)

	assert.Equal(t, req.Header.Get("content-type"), "application/json")

	assert.NotEqual(t, req.Header.Get("anthropic-version"), "")

	assert.NotEqual(t, req.Header.Get("x-api-key"), "")

	assert.Equal(t, req.URL.Scheme, "https")
	assert.Equal(t, req.URL.Host, "mockurl.com")

	assert.Equal(t, req.Method, "POST")

	var jsonBody struct {
		Text string   `json:"text"`
		List []string `json:"list"`
	}
	err = json.NewDecoder(req.Body).Decode(&jsonBody)
	assert.NoError(t, err)
	defer req.Body.Close()

	if jsonBody.Text != mockBody.Text || jsonBody.List[0] != mockBody.List[0] || jsonBody.List[1] != mockBody.List[1] {
		t.Errorf("JSON body of the new request incorrect")
	}

	assert.Equal(t, jsonBody.Text, mockBody.Text)
	assert.Equal(t, jsonBody.List[0], mockBody.List[0])
	assert.Equal(t, jsonBody.List[1], mockBody.List[1])
}

type MockRoundTripper struct {
	roundTripFunc func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req), nil
}

func MockHTTPClient(mockRoundTripper *MockRoundTripper) *http.Client {
	return &http.Client{Transport: mockRoundTripper}
}

func TestSendRequest(t *testing.T) {
	mockRoundTripper := &MockRoundTripper{
		roundTripFunc: func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"key":"value"}`)),
				Header:     make(http.Header),
			}
		},
	}

	mockClient := MockHTTPClient(mockRoundTripper)
	client := NewClientWithConfig(ClientConfig{
		HTTPClient: mockClient,
	})

	request, err := client.newRequest(context.Background(), http.MethodPost, "")
	assert.NoError(t, err)

	var result map[string]string
	err = client.sendRequest(request, &result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])

	mockRoundTripper.roundTripFunc = func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":"bad request"}`)),
			Header:     make(http.Header),
		}
	}

	err = client.sendRequest(request, &result)
	assert.Error(t, err)
}

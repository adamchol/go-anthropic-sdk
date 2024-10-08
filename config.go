package anthropic

import "net/http"

const (
	anthropicAPIURLv1 = "https://api.anthropic.com/v1"
)

type APIVersion string

const (
	latest      APIVersion = "2023-06-01"
	v2023_06_01 APIVersion = "2023-06-01"
	v2023_01_01 APIVersion = "2023-01-01"
	initial     APIVersion = "2023-01-01"
)

// ClientConfig is a configuration of client
type ClientConfig struct {
	authToken string

	BaseUrl    string
	APIVersion APIVersion

	HTTPClient *http.Client
}

// DefaultConfig creates a standard configuration with api key.
// This method is called when creating new client with [NewClient]
func DefaultConfig(apiKey string) ClientConfig {
	return ClientConfig{
		authToken:  apiKey,
		BaseUrl:    anthropicAPIURLv1,
		APIVersion: latest,
		HTTPClient: &http.Client{},
	}
}

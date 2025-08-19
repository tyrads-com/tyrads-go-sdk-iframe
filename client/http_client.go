package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tyrads-com/tyrads-go-sdk-iframe/config"
)

type HttpClient struct {
	client *http.Client
	config *config.Config
}

type HttpError struct {
	Message string `json:"message"`
}

func NewHttpClient(cfg *config.Config) *HttpClient {
	return &HttpClient{
		client: &http.Client{},
		config: cfg,
	}
}

// DoRequest sends an HTTP request and returns the parsed JSON response body and error.
// On success (status code 2xx), it parses the response body into the provided result interface.
// On error, it attempts to extract an error message from the response body.
//
// Parameters:
//   - method: HTTP method (e.g., "GET", "POST").
//   - path: API endpoint path.
//   - body: Request payload to be marshaled to JSON (can be nil).
//
// Returns:
//   - parsed: The parsed JSON response body (nil if error or result is nil).
//   - err: Error if the request fails or the response status is not 2xx.
func (hc *HttpClient) DoRequest(method, path string, body interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s/%s%s", hc.config.SdkApiBaseURL, hc.config.SdkApiVersion, path)

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", hc.config.ApiKey)
	req.Header.Set("X-API-Secret", hc.config.ApiSecret)
	req.Header.Set("X-SDK-Version", hc.config.SdkApiVersion)
	req.Header.Set("X-SDK-Platform", hc.config.SdkPlatform)
	q := req.URL.Query()
	q.Add("lang", hc.config.Language)
	req.URL.RawQuery = q.Encode()

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("no response received from the server: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errMsg HttpError
		json.Unmarshal(bodyBytes, &errMsg)
		if errMsg.Message == "" {
			errMsg.Message = "Unknown error"
		}
		return nil, fmt.Errorf("%s", errMsg.Message)
	}

	var parsed any
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}
	return parsed, nil
}

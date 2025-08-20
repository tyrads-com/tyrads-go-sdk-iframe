package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tyrads-com/tyrads-go-sdk-iframe/config"
)

func TestNewHttpClient(t *testing.T) {
	cfg := config.NewConfig("test-key", "test-secret")
	client := NewHttpClient(cfg)

	if client == nil {
		t.Fatal("expected HTTP client instance, got nil")
	}

	if client.config != cfg {
		t.Error("expected config to be set")
	}

	if client.client == nil {
		t.Error("expected HTTP client to be initialized")
	}
}

func TestDoRequest_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"token": "test-token-123",
		},
		"success": true,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("expected X-API-Key test-key, got %s", r.Header.Get("X-API-Key"))
		}
		if r.Header.Get("X-API-Secret") != "test-secret" {
			t.Errorf("expected X-API-Secret test-secret, got %s", r.Header.Get("X-API-Secret"))
		}
		if r.Header.Get("X-SDK-Version") != "v3.0" {
			t.Errorf("expected X-SDK-Version v3.0, got %s", r.Header.Get("X-SDK-Version"))
		}
		if r.Header.Get("X-SDK-Platform") != "Web" {
			t.Errorf("expected X-SDK-Platform Web, got %s", r.Header.Get("X-SDK-Platform"))
		}

		// Verify query parameter
		if r.URL.Query().Get("lang") != "en" {
			t.Errorf("expected lang parameter en, got %s", r.URL.Query().Get("lang"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cfg := config.NewConfig("test-key", "test-secret", func(c *config.Config) {
		c.SdkApiBaseURL = server.URL
	})
	client := NewHttpClient(cfg)

	requestBody := map[string]interface{}{
		"publisherUserId": "user123",
		"age":             25,
		"gender":          1,
	}

	result, err := client.DoRequest("POST", "/auth", requestBody)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be map[string]interface{}")
	}

	if !resultMap["success"].(bool) {
		t.Error("expected success to be true")
	}
}

func TestDoRequest_HTTPError(t *testing.T) {
	errorResponse := HttpError{
		Message: "Invalid API key",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	cfg := config.NewConfig("invalid-key", "invalid-secret", func(c *config.Config) {
		c.SdkApiBaseURL = server.URL
	})
	client := NewHttpClient(cfg)

	_, err := client.DoRequest("POST", "/auth", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "Invalid API key" {
		t.Errorf("expected error 'Invalid API key', got '%s'", err.Error())
	}
}

func TestDoRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	cfg := config.NewConfig("test-key", "test-secret", func(c *config.Config) {
		c.SdkApiBaseURL = server.URL
	})
	client := NewHttpClient(cfg)

	_, err := client.DoRequest("GET", "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "failed to parse response body: invalid character 'i' looking for beginning of value" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestDoRequest_NetworkError(t *testing.T) {
	cfg := config.NewConfig("test-key", "test-secret", func(c *config.Config) {
		c.SdkApiBaseURL = "http://nonexistent-server-12345"
	})
	client := NewHttpClient(cfg)

	_, err := client.DoRequest("GET", "/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should contain "no response received from the server"
	if err.Error()[:38] != "no response received from the server: " {
		t.Errorf("unexpected error message format: %s", err.Error())
	}
}
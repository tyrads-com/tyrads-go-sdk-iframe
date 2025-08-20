package tyrads

import (
	"os"
	"testing"

	"github.com/tyrads-com/tyrads-go-sdk-iframe/contract"
)

func TestNewTyrAdsSdk(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		apiSecret string
		lang      string
		envKey    string
		envSecret string
		wantLang  string
	}{
		{
			name:      "with all parameters",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			lang:      "es",
			wantLang:  "es",
		},
		{
			name:      "with default language",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			lang:      "",
			wantLang:  "en",
		},
		{
			name:      "with env variables",
			apiKey:    "",
			apiSecret: "",
			lang:      "fr",
			envKey:    "env-key",
			envSecret: "env-secret",
			wantLang:  "fr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" {
				os.Setenv("TYRADS_API_KEY", tt.envKey)
				defer os.Unsetenv("TYRADS_API_KEY")
			}
			if tt.envSecret != "" {
				os.Setenv("TYRADS_API_SECRET", tt.envSecret)
				defer os.Unsetenv("TYRADS_API_SECRET")
			}

			sdk := NewTyrAdsSdk(tt.apiKey, tt.apiSecret, tt.lang)

			if sdk == nil {
				t.Fatal("expected SDK instance, got nil")
			}

			if sdk.config.Language != tt.wantLang {
				t.Errorf("expected language %s, got %s", tt.wantLang, sdk.config.Language)
			}

			expectedKey := tt.apiKey
			if expectedKey == "" {
				expectedKey = tt.envKey
			}
			if sdk.config.ApiKey != expectedKey {
				t.Errorf("expected API key %s, got %s", expectedKey, sdk.config.ApiKey)
			}
		})
	}
}

func TestIframeUrl(t *testing.T) {
	sdk := NewTyrAdsSdk("test-key", "test-secret", "en")

	tests := []struct {
		name             string
		authSignOrToken  interface{}
		deeplinkTo       *string
		expectedURL      string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:            "with string token",
			authSignOrToken: "test-token",
			expectedURL:     "https://sdk.tyrads.com?token=test-token",
		},
		{
			name:            "with AuthenticationSign",
			authSignOrToken: contract.NewAuthenticationSign("auth-token", "user123", 25, 1),
			expectedURL:     "https://sdk.tyrads.com?token=auth-token",
		},
		{
			name:            "with deeplink",
			authSignOrToken: "test-token",
			deeplinkTo:      stringPtr("offers"),
			expectedURL:     "https://sdk.tyrads.com?token=test-token&to=offers",
		},
		{
			name:             "with empty deeplink",
			authSignOrToken:  "test-token",
			deeplinkTo:       stringPtr(""),
			expectError:      true,
			expectedErrorMsg: "invalid deeplinkTo argument: must be a non-empty string or nil",
		},
		{
			name:             "with invalid auth type",
			authSignOrToken:  123,
			expectError:      true,
			expectedErrorMsg: "invalid argument: must be an AuthenticationSign or a string token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := sdk.IframeUrl(tt.authSignOrToken, tt.deeplinkTo)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if err.Error() != tt.expectedErrorMsg {
					t.Errorf("expected error '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if url != tt.expectedURL {
				t.Errorf("expected URL '%s', got '%s'", tt.expectedURL, url)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

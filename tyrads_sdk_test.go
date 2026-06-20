package tyrads

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
			expectedURL:     "https://v4.sdk.tyrads.com?token=test-token",
		},
		{
			name:            "with AuthenticationSign",
			authSignOrToken: contract.NewAuthenticationSign("auth-token", "user123"),
			expectedURL:     "https://v4.sdk.tyrads.com?token=auth-token",
		},
		{
			name:            "with deeplink",
			authSignOrToken: "test-token",
			deeplinkTo:      stringPtr("offers"),
			expectedURL:     "https://v4.sdk.tyrads.com?token=test-token&to=offers",
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

func TestNewTyrAdsSdk_DefaultsToLatestApiVersion(t *testing.T) {
	sdk := NewTyrAdsSdk("k", "s", "en")
	if sdk.config.SdkApiVersion != "v4.0" {
		t.Errorf("expected default SdkApiVersion v4.0, got %s", sdk.config.SdkApiVersion)
	}
}

func TestNewTyrAdsSdk_WithApiVersionOverride(t *testing.T) {
	sdk := NewTyrAdsSdk("k", "s", "en", WithApiVersion("v3.0"))
	if sdk.config.SdkApiVersion != "v3.0" {
		t.Errorf("expected overridden SdkApiVersion v3.0, got %s", sdk.config.SdkApiVersion)
	}
}

func TestAuthenticate(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		sdk := NewTyrAdsSdk("test-key", "test-secret", "en")
		request := contract.NewAuthenticationRequest("")

		result, err := sdk.Authenticate(*request)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "validation error") {
			t.Errorf("expected error to contain 'validation error', got '%s'", err.Error())
		}
		if result != nil {
			t.Error("expected nil result when error occurs")
		}
	})

	t.Run("v4.0 hits /initialize/auth and returns token", func(t *testing.T) {
		var gotPath string
		var gotBody map[string]any
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			b, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(b, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"token": "tok-v4"},
			})
		}))
		defer server.Close()

		sdk := NewTyrAdsSdk("k", "s", "en")
		sdk.config.SdkApiBaseURL = server.URL

		req := contract.NewAuthenticationRequest("u-v4",
			contract.WithSub1("s1"),
			contract.WithUserGroup(map[string]any{"tier": "vip"}),
		)
		sign, err := sdk.Authenticate(*req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sign == nil || sign.Token != "tok-v4" {
			t.Fatalf("expected token 'tok-v4', got %+v", sign)
		}
		if gotPath != "/v4.0/initialize/auth" {
			t.Errorf("expected request path '/v4.0/initialize/auth', got %q", gotPath)
		}
		if gotBody["userGroup"] != `{"tier":"vip"}` {
			t.Errorf("expected userGroup to arrive as JSON-encoded string, got %v", gotBody["userGroup"])
		}
	})

	t.Run("v3.0 override hits /auth", func(t *testing.T) {
		var gotPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"token": "tok-v3"},
			})
		}))
		defer server.Close()

		sdk := NewTyrAdsSdk("k", "s", "en", WithApiVersion("v3.0"))
		sdk.config.SdkApiBaseURL = server.URL

		req := contract.NewAuthenticationRequest("u-v3")
		sign, err := sdk.Authenticate(*req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sign == nil || sign.Token != "tok-v3" {
			t.Fatalf("expected token 'tok-v3', got %+v", sign)
		}
		if gotPath != "/v3.0/auth" {
			t.Errorf("expected request path '/v3.0/auth', got %q", gotPath)
		}
	})
}

func TestIframeUrl_VersionAwareHost(t *testing.T) {
	t.Run("v3 override uses legacy sdk.tyrads.com host", func(t *testing.T) {
		sdk := NewTyrAdsSdk("k", "s", "en", WithApiVersion("v3.0"))

		got, err := sdk.IframeUrl("tok", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "https://sdk.tyrads.com?token=tok"
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}

		gotWidget, err := sdk.IframePremiumWidget("tok", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		wantWidget := "https://sdk.tyrads.com/widget?token=tok"
		if gotWidget != wantWidget {
			t.Errorf("expected %s, got %s", wantWidget, gotWidget)
		}
	})

	t.Run("explicit IFrameBaseURL override wins over version derivation", func(t *testing.T) {
		sdk := NewTyrAdsSdk("k", "s", "en")
		sdk.config.IFrameBaseURL = "https://staging.example.com"

		got, err := sdk.IframeUrl("tok", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "https://staging.example.com?token=tok"
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}

		gotWidget, err := sdk.IframePremiumWidget("tok", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		wantWidget := "https://staging.example.com/widget?token=tok"
		if gotWidget != wantWidget {
			t.Errorf("expected %s, got %s", wantWidget, gotWidget)
		}
	})
}

func TestIframePremiumWidget(t *testing.T) {
	sdk := NewTyrAdsSdk("test-key", "test-secret", "en")

	tests := []struct {
		name             string
		authSignOrToken  interface{}
		name_param       *string
		expectedURL      string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:            "with string token",
			authSignOrToken: "test-token",
			expectedURL:     "https://v4.sdk.tyrads.com/widget?token=test-token",
		},
		{
			name:            "with AuthenticationSign",
			authSignOrToken: contract.NewAuthenticationSign("auth-token", "user123"),
			expectedURL:     "https://v4.sdk.tyrads.com/widget?token=auth-token",
		},
		{
			name:            "with name parameter",
			authSignOrToken: "test-token",
			name_param:      stringPtr("premium-offers"),
			expectedURL:     "https://v4.sdk.tyrads.com/widget?token=test-token&name=premium-offers",
		},
		{
			name:             "with empty name",
			authSignOrToken:  "test-token",
			name_param:       stringPtr(""),
			expectError:      true,
			expectedErrorMsg: "invalid name argument: must be a non-empty string or nil",
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
			url, err := sdk.IframePremiumWidget(tt.authSignOrToken, tt.name_param)

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

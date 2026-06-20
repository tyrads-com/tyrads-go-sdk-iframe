package config

import "testing"

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		apiSecret string
		opts      []ConfigOptions
		expected  *Config
	}{
		{
			name:      "default config",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			expected: &Config{
				IFrameBaseURL: "",
				SdkApiBaseURL: "https://api.tyrads.com",
				SdkApiVersion: "v4.0",
				SdkPlatform:   "Web",
				ApiKey:        "test-key",
				ApiSecret:     "test-secret",
				Language:      "en",
			},
		},
		{
			name:      "with custom language",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			opts: []ConfigOptions{
				func(c *Config) { c.Language = "es" },
			},
			expected: &Config{
				IFrameBaseURL: "",
				SdkApiBaseURL: "https://api.tyrads.com",
				SdkApiVersion: "v4.0",
				SdkPlatform:   "Web",
				ApiKey:        "test-key",
				ApiSecret:     "test-secret",
				Language:      "es",
			},
		},
		{
			name:      "with custom iframe URL",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			opts: []ConfigOptions{
				func(c *Config) { c.IFrameBaseURL = "https://custom.domain.com" },
			},
			expected: &Config{
				IFrameBaseURL: "https://custom.domain.com",
				SdkApiBaseURL: "https://api.tyrads.com",
				SdkApiVersion: "v4.0",
				SdkPlatform:   "Web",
				ApiKey:        "test-key",
				ApiSecret:     "test-secret",
				Language:      "en",
			},
		},
		{
			name:      "with overridden api version",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			opts: []ConfigOptions{
				func(c *Config) { c.SdkApiVersion = "v3.0" },
			},
			expected: &Config{
				IFrameBaseURL: "",
				SdkApiBaseURL: "https://api.tyrads.com",
				SdkApiVersion: "v3.0",
				SdkPlatform:   "Web",
				ApiKey:        "test-key",
				ApiSecret:     "test-secret",
				Language:      "en",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.apiKey, tt.apiSecret, tt.opts...)

			if config.IFrameBaseURL != tt.expected.IFrameBaseURL {
				t.Errorf("expected IFrameBaseURL %s, got %s", tt.expected.IFrameBaseURL, config.IFrameBaseURL)
			}
			if config.SdkApiBaseURL != tt.expected.SdkApiBaseURL {
				t.Errorf("expected SdkApiBaseURL %s, got %s", tt.expected.SdkApiBaseURL, config.SdkApiBaseURL)
			}
			if config.SdkApiVersion != tt.expected.SdkApiVersion {
				t.Errorf("expected SdkApiVersion %s, got %s", tt.expected.SdkApiVersion, config.SdkApiVersion)
			}
			if config.SdkPlatform != tt.expected.SdkPlatform {
				t.Errorf("expected SdkPlatform %s, got %s", tt.expected.SdkPlatform, config.SdkPlatform)
			}
			if config.ApiKey != tt.expected.ApiKey {
				t.Errorf("expected ApiKey %s, got %s", tt.expected.ApiKey, config.ApiKey)
			}
			if config.ApiSecret != tt.expected.ApiSecret {
				t.Errorf("expected ApiSecret %s, got %s", tt.expected.ApiSecret, config.ApiSecret)
			}
			if config.Language != tt.expected.Language {
				t.Errorf("expected Language %s, got %s", tt.expected.Language, config.Language)
			}
		})
	}
}

func TestResolveIFrameBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name:     "default v4 resolves to versioned subdomain",
			config:   &Config{SdkApiVersion: "v4.0"},
			expected: "https://v4.sdk.tyrads.com",
		},
		{
			name:     "v3 resolves to legacy host",
			config:   &Config{SdkApiVersion: "v3.0"},
			expected: "https://sdk.tyrads.com",
		},
		{
			name:     "future v5 resolves to versioned subdomain",
			config:   &Config{SdkApiVersion: "v5.0"},
			expected: "https://v5.sdk.tyrads.com",
		},
		{
			name:     "explicit override wins over version derivation",
			config:   &Config{SdkApiVersion: "v4.0", IFrameBaseURL: "https://custom.domain.com"},
			expected: "https://custom.domain.com",
		},
		{
			name:     "empty version falls back to v4 host",
			config:   &Config{SdkApiVersion: ""},
			expected: "https://v4.sdk.tyrads.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.ResolveIFrameBaseURL()
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

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
				IFrameBaseURL: "https://sdk.tyrads.com",
				SdkApiBaseURL: "https://api.tyrads.com",
				SdkApiVersion: "v3.0",
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
				IFrameBaseURL: "https://sdk.tyrads.com",
				SdkApiBaseURL: "https://api.tyrads.com",
				SdkApiVersion: "v3.0",
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

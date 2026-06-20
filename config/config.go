package config

import (
	"fmt"
	"strings"
)

type Config struct {
	IFrameBaseURL string
	SdkApiBaseURL string
	SdkApiVersion string
	SdkPlatform   string
	ApiKey        string
	ApiSecret     string
	Language      string
}

type ConfigOptions func(*Config)

func NewConfig(apiKey, apiSecret string, opts ...ConfigOptions) *Config {
	c := new(Config)
	c.IFrameBaseURL = ""
	c.SdkApiBaseURL = "https://api.tyrads.com"
	c.SdkApiVersion = "v4.0"
	c.SdkPlatform = "Web"
	c.ApiKey = apiKey
	c.ApiSecret = apiSecret
	c.Language = "en"

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ResolveIFrameBaseURL returns the iframe host to use for the current configuration.
// If IFrameBaseURL is explicitly set, it is returned as-is. Otherwise the host is
// derived from SdkApiVersion: v3.x uses the legacy sdk.tyrads.com host, and v4+
// uses the versioned subdomain (v4.sdk.tyrads.com, v5.sdk.tyrads.com, ...).
func (c *Config) ResolveIFrameBaseURL() string {
	if c.IFrameBaseURL != "" {
		return c.IFrameBaseURL
	}
	return iframeBaseURLForVersion(c.SdkApiVersion)
}

func iframeBaseURLForVersion(version string) string {
	major := strings.SplitN(version, ".", 2)[0]
	if major == "v3" {
		return "https://sdk.tyrads.com"
	}
	if major == "" {
		major = "v4"
	}
	return fmt.Sprintf("https://%s.sdk.tyrads.com", major)
}

package config

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
	c.IFrameBaseURL = "https://sdk.tyrads.com"
	c.SdkApiBaseURL = "https://api.tyrads.com"
	c.SdkApiVersion = "v3.0"
	c.SdkPlatform = "Web"
	c.ApiKey = apiKey
	c.ApiSecret = apiSecret
	c.Language = "en"

	for _, opt := range opts {
		opt(c)
	}

	return c
}

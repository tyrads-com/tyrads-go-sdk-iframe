package tyrads

import "github.com/tyrads-com/tyrads-go-sdk-iframe/config"

// SdkOption configures optional SDK behavior at construction time.
type SdkOption func(*config.Config)

// WithApiVersion overrides the API version used for backend requests.
// If not provided, the SDK uses the latest version baked into config.
func WithApiVersion(v string) SdkOption {
	return func(c *config.Config) { c.SdkApiVersion = v }
}

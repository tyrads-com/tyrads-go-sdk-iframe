package tyrads

import (
	"fmt"
	"net/url"
	"os"

	"github.com/tyrads-com/tyrads-go-sdk-iframe/client"
	"github.com/tyrads-com/tyrads-go-sdk-iframe/config"
	"github.com/tyrads-com/tyrads-go-sdk-iframe/contract"
	"github.com/tyrads-com/tyrads-go-sdk-iframe/enum"
)

type AuthenticationRequest = contract.AuthenticationRequest
type AuthenticationSign = contract.AuthenticationSign

type TyrAdsSdk struct {
	config     *config.Config
	httpClient *client.HttpClient
}

// NewTyrAdsSdk creates and returns a new instance of TyrAdsSdk with the specified configuration.
// It initializes the SDK with API credentials and language settings.
//
// Parameters:
//   - apiKey: The API key for authentication. If empty, it will be retrieved from the TYRADS_API_KEY environment variable.
//   - apiSecret: The API secret for authentication. If empty, it will be retrieved from the TYRADS_API_SECRET environment variable.
//   - lang: The language code for SDK responses. Defaults to "en" if not specified or empty.
//   - opts: Optional SdkOption values, e.g. WithApiVersion("v3.0"). If omitted, the SDK uses the latest API version.
//
// Returns:
//   - *TyrAdsSdk: A pointer to the newly created TyrAdsSdk instance configured with the provided parameters.
func NewTyrAdsSdk(apiKey, apiSecret, lang string, opts ...SdkOption) *TyrAdsSdk {
	if apiKey == "" {
		apiKey = os.Getenv(string(enum.TYRADS_API_KEY))
	}
	if apiSecret == "" {
		apiSecret = os.Getenv(string(enum.TYRADS_API_SECRET))
	}
	if lang == "" {
		lang = "en"
	}
	cfg := config.NewConfig(apiKey, apiSecret, func(c *config.Config) {
		c.Language = lang
	})
	for _, opt := range opts {
		opt(cfg)
	}
	return &TyrAdsSdk{
		config:     cfg,
		httpClient: client.NewHttpClient(cfg),
	}
}

// Authenticate performs authentication using the provided request and returns an AuthenticationSign.
// It validates the authentication request, makes a POST request to the authentication endpoint,
// and processes the response to create an AuthenticationSign containing the authentication token
// and user information.
//
// Parameters:
//   - request: AuthenticationRequest containing the authentication details
//
// Returns:
//   - *AuthenticationSign: Contains the authentication token and user information
//   - error: Returns an error if validation fails, request fails, or response parsing fails
func (sdk *TyrAdsSdk) Authenticate(request AuthenticationRequest) (*AuthenticationSign, error) {
	if err := request.ValidateAuthenticationRequest(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	data := request.GetParsedAuthenticationRequestData()
	path := "/initialize/auth"
	if sdk.config.SdkApiVersion == "v3.0" {
		path = "/auth"
	}
	resp, err := sdk.httpClient.DoRequest("POST", path, data)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	respMap, ok := resp.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	dataMap, ok := respMap["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	token, ok := dataMap["token"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token format")
	}

	return contract.NewAuthenticationSign(token, request.PublisherUserID), nil
}

// IframeOption configures optional iframe URL parameters such as placementId.
type IframeOption func(*iframeOptions)

type iframeOptions struct {
	placementID *int
}

// WithPlacementID attaches a placementId query parameter to the generated iframe URL.
// Only supported on iframe v4 and above; passing it to a v3-configured SDK returns an error.
// The value must be a positive integer.
func WithPlacementID(v int) IframeOption {
	return func(o *iframeOptions) { o.placementID = &v }
}

// IframeUrl generates a URL for an iframe integration with authentication.
// It accepts either a string token or an AuthenticationSign struct pointer as the first parameter,
// and an optional deeplinkTo string pointer for specifying a target destination.
//
// Parameters:
//   - authSignOrToken: Either a string token or *AuthenticationSign for authentication
//   - deeplinkTo: Optional pointer to a string specifying the target destination
//   - opts: Optional IframeOption values, e.g. WithPlacementID(123). v4+ only.
//
// Returns:
//   - string: The generated iframe URL with authentication and optional deeplink/placement parameters
//   - error: An error if invalid arguments are provided
func (sdk *TyrAdsSdk) IframeUrl(authSignOrToken interface{}, deeplinkTo *string, opts ...IframeOption) (string, error) {
	var token string

	switch v := authSignOrToken.(type) {
	case string:
		token = v
	case *AuthenticationSign:
		token = v.Token
	default:
		return "", fmt.Errorf("invalid argument: must be an AuthenticationSign or a string token")
	}

	if deeplinkTo != nil && *deeplinkTo == "" {
		return "", fmt.Errorf("invalid deeplinkTo argument: must be a non-empty string or nil")
	}

	options := &iframeOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if err := sdk.validateIframeOptions(options); err != nil {
		return "", err
	}

	iframeUrl := fmt.Sprintf("%s?token=%s", sdk.config.ResolveIFrameBaseURL(), url.QueryEscape(token))
	if deeplinkTo != nil {
		iframeUrl += fmt.Sprintf("&to=%s", url.QueryEscape(*deeplinkTo))
	}
	if options.placementID != nil {
		iframeUrl += fmt.Sprintf("&placementId=%d", *options.placementID)
	}

	return iframeUrl, nil
}

// IframePremiumWidget generates a URL for embedding a premium widget iframe.
// It accepts either an authentication sign or a token string as the first parameter,
// and an optional name parameter.
//
// Parameters:
//   - authSignOrToken: Can be either an *AuthenticationSign or a string token
//   - name: Optional pointer to a string for naming the widget. If provided, must be non-empty
//   - opts: Optional IframeOption values, e.g. WithPlacementID(123). v4+ only.
//
// Returns:
//   - string: The generated iframe URL
//   - error: An error if invalid parameters are provided
//
// The function will return an error if:
//   - authSignOrToken is neither an AuthenticationSign nor a string
//   - name pointer is provided but points to an empty string
//   - WithPlacementID is supplied with a non-positive value or on a v3 SDK
func (sdk *TyrAdsSdk) IframePremiumWidget(authSignOrToken interface{}, name *string, opts ...IframeOption) (string, error) {
	var token string

	switch v := authSignOrToken.(type) {
	case string:
		token = v
	case *AuthenticationSign:
		token = v.Token
	default:
		return "", fmt.Errorf("invalid argument: must be an AuthenticationSign or a string token")
	}

	if name != nil && *name == "" {
		return "", fmt.Errorf("invalid name argument: must be a non-empty string or nil")
	}

	options := &iframeOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if err := sdk.validateIframeOptions(options); err != nil {
		return "", err
	}

	iframeUrl := fmt.Sprintf("%s/widget?token=%s", sdk.config.ResolveIFrameBaseURL(), url.QueryEscape(token))
	if name != nil {
		iframeUrl += fmt.Sprintf("&name=%s", url.QueryEscape(*name))
	}
	if options.placementID != nil {
		iframeUrl += fmt.Sprintf("&placementId=%d", *options.placementID)
	}

	return iframeUrl, nil
}

func (sdk *TyrAdsSdk) validateIframeOptions(o *iframeOptions) error {
	if o.placementID == nil {
		return nil
	}
	if *o.placementID <= 0 {
		return fmt.Errorf("invalid placementId argument: must be a positive integer")
	}
	if !config.IsV4OrAbove(sdk.config.SdkApiVersion) {
		return fmt.Errorf("placementId is only supported on iframe v4 and above")
	}
	return nil
}

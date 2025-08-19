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
//
// Returns:
//   - *TyrAdsSdk: A pointer to the newly created TyrAdsSdk instance configured with the provided parameters.
func NewTyrAdsSdk(apiKey, apiSecret, lang string) *TyrAdsSdk {
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
	resp, err := sdk.httpClient.DoRequest("POST", "/auth", data)
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

	return contract.NewAuthenticationSign(token, request.PublisherUserID, request.Age, request.Gender), nil
}

// IframeUrl generates a URL for an iframe integration with authentication.
// It accepts either a string token or an AuthenticationSign struct pointer as the first parameter,
// and an optional deeplinkTo string pointer for specifying a target destination.
//
// Parameters:
//   - authSignOrToken: Either a string token or *AuthenticationSign for authentication
//   - deeplinkTo: Optional pointer to a string specifying the target destination
//
// Returns:
//   - string: The generated iframe URL with authentication and optional deeplink parameters
//   - error: An error if invalid arguments are provided
func (sdk *TyrAdsSdk) IframeUrl(authSignOrToken interface{}, deeplinkTo *string) (string, error) {
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

	iframeUrl := fmt.Sprintf("%s?token=%s", sdk.config.IFrameBaseURL, url.QueryEscape(token))
	if deeplinkTo != nil {
		iframeUrl += fmt.Sprintf("&to=%s", url.QueryEscape(*deeplinkTo))
	}

	return iframeUrl, nil
}

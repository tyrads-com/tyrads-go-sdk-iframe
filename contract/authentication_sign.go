package contract

// AuthenticationSign represents authentication information.
type AuthenticationSign struct {
	Token           string
	PublisherUserID string
}

// NewAuthenticationSign creates a new AuthenticationSign instance.
func NewAuthenticationSign(token, publisherUserID string) *AuthenticationSign {
	return &AuthenticationSign{
		Token:           token,
		PublisherUserID: publisherUserID,
	}
}

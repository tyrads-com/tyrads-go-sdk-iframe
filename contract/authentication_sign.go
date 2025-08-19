package contract

// AuthenticationSign represents authentication information.
type AuthenticationSign struct {
	Token           string
	PublisherUserID string
	Age             int
	Gender          int // 1 = male, 2 = female
}

// NewAuthenticationSign creates a new AuthenticationSign instance.
func NewAuthenticationSign(token, publisherUserID string, age, gender int) *AuthenticationSign {
	return &AuthenticationSign{
		Token:           token,
		PublisherUserID: publisherUserID,
		Age:             age,
		Gender:          gender,
	}
}

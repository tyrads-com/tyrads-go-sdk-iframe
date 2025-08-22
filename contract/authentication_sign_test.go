package contract

import "testing"

func TestNewAuthenticationSign(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		publisherUserID string
		expected        *AuthenticationSign
	}{
		{
			name:            "valid authentication sign",
			token:           "test-token-123",
			publisherUserID: "user456",
			expected: &AuthenticationSign{
				Token:           "test-token-123",
				PublisherUserID: "user456",
			},
		},
		{
			name:            "another valid sign",
			token:           "another-token",
			publisherUserID: "user789",
			expected: &AuthenticationSign{
				Token:           "another-token",
				PublisherUserID: "user789",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAuthenticationSign(tt.token, tt.publisherUserID)

			if result.Token != tt.expected.Token {
				t.Errorf("expected Token %s, got %s", tt.expected.Token, result.Token)
			}
			if result.PublisherUserID != tt.expected.PublisherUserID {
				t.Errorf("expected PublisherUserID %s, got %s", tt.expected.PublisherUserID, result.PublisherUserID)
			}
		})
	}
}

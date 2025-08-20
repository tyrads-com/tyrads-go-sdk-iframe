package contract

import "testing"

func TestNewAuthenticationSign(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		publisherUserID string
		age             int
		gender          int
		expected        *AuthenticationSign
	}{
		{
			name:            "valid authentication sign",
			token:           "test-token-123",
			publisherUserID: "user456",
			age:             30,
			gender:          2,
			expected: &AuthenticationSign{
				Token:           "test-token-123",
				PublisherUserID: "user456",
				Age:             30,
				Gender:          2,
			},
		},
		{
			name:            "with male gender",
			token:           "another-token",
			publisherUserID: "user789",
			age:             25,
			gender:          1,
			expected: &AuthenticationSign{
				Token:           "another-token",
				PublisherUserID: "user789",
				Age:             25,
				Gender:          1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAuthenticationSign(tt.token, tt.publisherUserID, tt.age, tt.gender)

			if result.Token != tt.expected.Token {
				t.Errorf("expected Token %s, got %s", tt.expected.Token, result.Token)
			}
			if result.PublisherUserID != tt.expected.PublisherUserID {
				t.Errorf("expected PublisherUserID %s, got %s", tt.expected.PublisherUserID, result.PublisherUserID)
			}
			if result.Age != tt.expected.Age {
				t.Errorf("expected Age %d, got %d", tt.expected.Age, result.Age)
			}
			if result.Gender != tt.expected.Gender {
				t.Errorf("expected Gender %d, got %d", tt.expected.Gender, result.Gender)
			}
		})
	}
}

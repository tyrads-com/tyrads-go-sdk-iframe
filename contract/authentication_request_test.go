package contract

import (
	"reflect"
	"testing"
)

func TestNewAuthenticationRequest(t *testing.T) {
	email := "test@example.com"
	phone := "+1234567890"
	age25 := 25
	age30 := 30
	gender1 := 1
	gender2 := 2

	tests := []struct {
		name            string
		publisherUserID string
		opts            []AuthenticationRequestOptions
		expected        *AuthenticationRequest
	}{
		{
			name:            "basic request with age and gender",
			publisherUserID: "user123",
			opts: []AuthenticationRequestOptions{
				func(ar *AuthenticationRequest) { ar.Age = &age25 },
				func(ar *AuthenticationRequest) { ar.Gender = &gender1 },
			},
			expected: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             &age25,
				Gender:          &gender1,
			},
		},
		{
			name:            "with email and phone",
			publisherUserID: "user456",
			opts: []AuthenticationRequestOptions{
				func(ar *AuthenticationRequest) { ar.Age = &age30 },
				func(ar *AuthenticationRequest) { ar.Gender = &gender2 },
				func(ar *AuthenticationRequest) { ar.Email = &email },
				func(ar *AuthenticationRequest) { ar.PhoneNumber = &phone },
			},
			expected: &AuthenticationRequest{
				PublisherUserID: "user456",
				Age:             &age30,
				Gender:          &gender2,
				Email:           &email,
				PhoneNumber:     &phone,
			},
		},
		{
			name:            "without age and gender",
			publisherUserID: "user789",
			expected: &AuthenticationRequest{
				PublisherUserID: "user789",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewAuthenticationRequest(tt.publisherUserID, tt.opts...)

			if !reflect.DeepEqual(req, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, req)
			}
		})
	}
}

func TestValidateAuthenticationRequest(t *testing.T) {
	validEmail := "test@example.com"
	invalidEmail := "invalid-email"
	validPhone := "+1234567890"
	invalidPhone := "invalid-phone"

	tests := []struct {
		name    string
		request *AuthenticationRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with age and gender",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "valid request without age and gender",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
			},
			wantErr: false,
		},
		{
			name: "empty publisher user ID",
			request: &AuthenticationRequest{
				PublisherUserID: "",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
			},
			wantErr: true,
			errMsg:  "publisher user ID cannot be empty and must be a string",
		},
		{
			name: "negative age",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := -1; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
			},
			wantErr: true,
			errMsg:  "age must be a non-negative integer",
		},
		{
			name: "invalid gender",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 3; return &v }(),
			},
			wantErr: true,
			errMsg:  "gender must be either 1 (male) or 2 (female)",
		},
		{
			name: "valid email",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
				Email:           &validEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
				Email:           &invalidEmail,
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "valid phone",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
				PhoneNumber:     &validPhone,
			},
			wantErr: false,
		},
		{
			name: "invalid phone",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Age:             func() *int { v := 25; return &v }(),
				Gender:          func() *int { v := 1; return &v }(),
				PhoneNumber:     &invalidPhone,
			},
			wantErr: true,
			errMsg:  "invalid phone number format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.ValidateAuthenticationRequest()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetParsedAuthenticationRequestData(t *testing.T) {
	email := "test@example.com"
	phone := "+1234567890"
	sub1 := "sub1-value"
	incentivized := true

	request := &AuthenticationRequest{
		PublisherUserID: "user123",
		Age:             func() *int { v := 25; return &v }(),
		Gender:          func() *int { v := 1; return &v }(),
		Email:           &email,
		PhoneNumber:     &phone,
		Sub1:            &sub1,
		Incentivized:    &incentivized,
	}

	expected := map[string]interface{}{
		"publisherUserId": "user123",
		"age":             25,
		"gender":          1,
		"email":           "test@example.com",
		"phoneNumber":     "+1234567890",
		"sub1":            "sub1-value",
		"incentivized":    true,
	}

	result := request.GetParsedAuthenticationRequestData()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestGetParsedAuthenticationRequestData_SkipsNilAndEmptyFields(t *testing.T) {
	emptyString := ""
	sub1 := "sub1-value"

	request := &AuthenticationRequest{
		PublisherUserID: "user123",
		Age:             func() *int { v := 25; return &v }(),
		Gender:          func() *int { v := 1; return &v }(),
		Email:           &emptyString, // should be skipped
		Sub1:            &sub1,        // should be included
		Sub2:            nil,          // should be skipped
	}

	expected := map[string]interface{}{
		"publisherUserId": "user123",
		"age":             25,
		"gender":          1,
		"sub1":            "sub1-value",
	}

	result := request.GetParsedAuthenticationRequestData()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

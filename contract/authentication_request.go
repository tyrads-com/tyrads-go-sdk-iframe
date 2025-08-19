package contract

import (
	"errors"
	"fmt"
	"regexp"
)

// AuthenticationRequest represents a request for user authentication.
type AuthenticationRequest struct {
	PublisherUserID   string  `json:"publisherUserId"`
	Age               int     `json:"age"`
	Gender            int     `json:"gender"`
	Email             *string `json:"email,omitempty"`
	PhoneNumber       *string `json:"phoneNumber,omitempty"`
	Sub1              *string `json:"sub1,omitempty"`
	Sub2              *string `json:"sub2,omitempty"`
	Sub3              *string `json:"sub3,omitempty"`
	Sub4              *string `json:"sub4,omitempty"`
	Sub5              *string `json:"sub5,omitempty"`
	UserGroup         *string `json:"userGroup,omitempty"`
	MediaSourceName   *string `json:"mediaSourceName,omitempty"`
	MediaSourceID     *string `json:"mediaSourceId,omitempty"`
	MediaSubSourceID  *string `json:"mediaSubSourceId,omitempty"`
	Incentivized      *bool   `json:"incentivized,omitempty"`
	MediaAdsetName    *string `json:"mediaAdsetName,omitempty"`
	MediaAdsetID      *string `json:"mediaAdsetId,omitempty"`
	MediaCreativeName *string `json:"mediaCreativeName,omitempty"`
	MediaCreativeID   *string `json:"mediaCreativeId,omitempty"`
	MediaCampaignName *string `json:"mediaCampaignName,omitempty"`
}

type AuthenticationRequestOptions func(*AuthenticationRequest)

// NewAuthenticationRequest creates a new AuthenticationRequest instance.
func NewAuthenticationRequest(publisherUserID string, age, gender int, opts ...AuthenticationRequestOptions) *AuthenticationRequest {
	req := &AuthenticationRequest{
		PublisherUserID: publisherUserID,
		Age:             age,
		Gender:          gender,
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// ValidateAuthenticationRequest validates an AuthenticationRequest.
// Returns error if validation fails.
func (ar *AuthenticationRequest) ValidateAuthenticationRequest() error {
	if ar.PublisherUserID == "" {
		return errors.New("publisher user ID cannot be empty and must be a string")
	}
	if ar.Age < 0 {
		return errors.New("age must be a non-negative integer")
	}
	if ar.Gender != 1 && ar.Gender != 2 {
		return errors.New("gender must be either 1 (male) or 2 (female)")
	}
	if ar.Email != nil {
		emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
		if !emailRegex.MatchString(*ar.Email) {
			return errors.New("invalid email format")
		}
	}
	if ar.PhoneNumber != nil {
		phoneRegex := regexp.MustCompile(`^\+?[0-9\- ]{7,20}$`)
		if !phoneRegex.MatchString(*ar.PhoneNumber) {
			return errors.New("invalid phone number format")
		}
	}
	stringFields := map[string]*string{
		"sub1":              ar.Sub1,
		"sub2":              ar.Sub2,
		"sub3":              ar.Sub3,
		"sub4":              ar.Sub4,
		"sub5":              ar.Sub5,
		"userGroup":         ar.UserGroup,
		"mediaSourceName":   ar.MediaSourceName,
		"mediaSourceId":     ar.MediaSourceID,
		"mediaSubSourceId":  ar.MediaSubSourceID,
		"mediaAdsetName":    ar.MediaAdsetName,
		"mediaAdsetId":      ar.MediaAdsetID,
		"mediaCreativeName": ar.MediaCreativeName,
		"mediaCreativeId":   ar.MediaCreativeID,
		"mediaCampaignName": ar.MediaCampaignName,
	}
	for field, value := range stringFields {
		if value != nil && fmt.Sprintf("%T", value) != "*string" {
			return fmt.Errorf("%s must be a string", field)
		}
	}
	if ar.Incentivized != nil && fmt.Sprintf("%T", ar.Incentivized) != "*bool" {
		return errors.New("incentivized must be a boolean")
	}
	return nil
}

// GetParsedAuthenticationRequestData returns a map containing the authentication request data.
// Only includes fields that are defined and non-empty.
func (ar *AuthenticationRequest) GetParsedAuthenticationRequestData() map[string]interface{} {
	data := map[string]interface{}{
		"age":             ar.Age,
		"gender":          ar.Gender,
		"publisherUserId": ar.PublisherUserID,
	}
	optionalFields := map[string]interface{}{
		"email":             ar.Email,
		"phoneNumber":       ar.PhoneNumber,
		"sub1":              ar.Sub1,
		"sub2":              ar.Sub2,
		"sub3":              ar.Sub3,
		"sub4":              ar.Sub4,
		"sub5":              ar.Sub5,
		"userGroup":         ar.UserGroup,
		"mediaSourceName":   ar.MediaSourceName,
		"mediaSourceId":     ar.MediaSourceID,
		"mediaSubSourceId":  ar.MediaSubSourceID,
		"incentivized":      ar.Incentivized,
		"mediaAdsetName":    ar.MediaAdsetName,
		"mediaAdsetId":      ar.MediaAdsetID,
		"mediaCreativeName": ar.MediaCreativeName,
		"mediaCreativeId":   ar.MediaCreativeID,
		"mediaCampaignName": ar.MediaCampaignName,
	}
	for key, value := range optionalFields {
		switch v := value.(type) {
		case *string:
			if v != nil && *v != "" {
				data[key] = *v
			}
		case *bool:
			if v != nil {
				data[key] = *v
			}
		}
	}
	return data
}

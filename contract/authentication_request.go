package contract

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

const (
	maxOptionalStringLength  = 255
	maxUserGroupStringLength = 4096
)

// AuthenticationRequest represents a request for user authentication.
type AuthenticationRequest struct {
	PublisherUserID   string      `json:"publisherUserId"`
	Age               *int        `json:"age,omitempty"`
	Gender            *int        `json:"gender,omitempty"`
	Email             *string     `json:"email,omitempty"`
	PhoneNumber       *string     `json:"phoneNumber,omitempty"`
	Sub1              *string     `json:"sub1,omitempty"`
	Sub2              *string     `json:"sub2,omitempty"`
	Sub3              *string     `json:"sub3,omitempty"`
	Sub4              *string     `json:"sub4,omitempty"`
	Sub5              *string     `json:"sub5,omitempty"`
	UserGroup         interface{} `json:"userGroup,omitempty"`
	MediaSourceName   *string     `json:"mediaSourceName,omitempty"`
	MediaSourceID     *string     `json:"mediaSourceId,omitempty"`
	MediaSubSourceID  *string     `json:"mediaSubSourceId,omitempty"`
	Incentivized      *bool       `json:"incentivized,omitempty"`
	MediaAdsetName    *string     `json:"mediaAdsetName,omitempty"`
	MediaAdsetID      *string     `json:"mediaAdsetId,omitempty"`
	MediaCreativeName *string     `json:"mediaCreativeName,omitempty"`
	MediaCreativeID   *string     `json:"mediaCreativeId,omitempty"`
	MediaCampaignName *string     `json:"mediaCampaignName,omitempty"`
}

type AuthenticationRequestOptions func(*AuthenticationRequest)

// NewAuthenticationRequest creates a new AuthenticationRequest instance.
func NewAuthenticationRequest(publisherUserID string, opts ...AuthenticationRequestOptions) *AuthenticationRequest {
	req := &AuthenticationRequest{
		PublisherUserID: publisherUserID,
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

// WithAge sets the Age field on an AuthenticationRequest.
func WithAge(v int) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Age = &v }
}

// WithGender sets the Gender field on an AuthenticationRequest.
func WithGender(v int) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Gender = &v }
}

// WithEmail sets the Email field on an AuthenticationRequest.
func WithEmail(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Email = &v }
}

// WithPhoneNumber sets the PhoneNumber field on an AuthenticationRequest.
func WithPhoneNumber(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.PhoneNumber = &v }
}

// WithSub1 sets the Sub1 tracking field on an AuthenticationRequest.
func WithSub1(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Sub1 = &v }
}

// WithSub2 sets the Sub2 tracking field on an AuthenticationRequest.
func WithSub2(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Sub2 = &v }
}

// WithSub3 sets the Sub3 tracking field on an AuthenticationRequest.
func WithSub3(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Sub3 = &v }
}

// WithSub4 sets the Sub4 tracking field on an AuthenticationRequest.
func WithSub4(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Sub4 = &v }
}

// WithSub5 sets the Sub5 tracking field on an AuthenticationRequest.
func WithSub5(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Sub5 = &v }
}

// WithUserGroup sets the UserGroup field. The value can be a string (sent as-is)
// or a struct/map (JSON-encoded to a string before being sent to the backend).
func WithUserGroup(v interface{}) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.UserGroup = v }
}

// WithMediaSourceName sets the MediaSourceName field on an AuthenticationRequest.
func WithMediaSourceName(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaSourceName = &v }
}

// WithMediaSourceID sets the MediaSourceID field on an AuthenticationRequest.
func WithMediaSourceID(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaSourceID = &v }
}

// WithMediaSubSourceID sets the MediaSubSourceID field on an AuthenticationRequest.
func WithMediaSubSourceID(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaSubSourceID = &v }
}

// WithIncentivized sets the Incentivized field on an AuthenticationRequest.
func WithIncentivized(v bool) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.Incentivized = &v }
}

// WithMediaAdsetName sets the MediaAdsetName field on an AuthenticationRequest.
func WithMediaAdsetName(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaAdsetName = &v }
}

// WithMediaAdsetID sets the MediaAdsetID field on an AuthenticationRequest.
func WithMediaAdsetID(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaAdsetID = &v }
}

// WithMediaCreativeName sets the MediaCreativeName field on an AuthenticationRequest.
func WithMediaCreativeName(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaCreativeName = &v }
}

// WithMediaCreativeID sets the MediaCreativeID field on an AuthenticationRequest.
func WithMediaCreativeID(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaCreativeID = &v }
}

// WithMediaCampaignName sets the MediaCampaignName field on an AuthenticationRequest.
func WithMediaCampaignName(v string) AuthenticationRequestOptions {
	return func(ar *AuthenticationRequest) { ar.MediaCampaignName = &v }
}

// ValidateAuthenticationRequest validates an AuthenticationRequest.
// Returns error if validation fails.
func (ar *AuthenticationRequest) ValidateAuthenticationRequest() error {
	if ar.PublisherUserID == "" {
		return errors.New("publisher user ID cannot be empty and must be a string")
	}
	if ar.Age != nil && *ar.Age < 0 {
		return errors.New("age must be a non-negative integer")
	}
	if ar.Gender != nil && (*ar.Gender != 1 && *ar.Gender != 2) {
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
		if value != nil && len(*value) > maxOptionalStringLength {
			return fmt.Errorf("%s must not exceed %d characters", field, maxOptionalStringLength)
		}
	}
	if err := validateUserGroup(ar.UserGroup); err != nil {
		return err
	}
	return nil
}

// validateUserGroup checks that userGroup is either a string, a struct,
// or a map. Primitives (int, bool, float), slices, channels, and functions
// are rejected. nil is treated as "not set".
func validateUserGroup(v interface{}) error {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		if len(s) > maxUserGroupStringLength {
			return fmt.Errorf("userGroup must not exceed %d characters", maxUserGroupStringLength)
		}
		return nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Struct:
		// json.Marshal will handle these as JSON objects.
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("userGroup could not be JSON-encoded: %w", err)
		}
		if len(b) > maxUserGroupStringLength {
			return fmt.Errorf("userGroup must not exceed %d characters when JSON-encoded", maxUserGroupStringLength)
		}
		return nil
	default:
		return fmt.Errorf("userGroup must be a string, struct, or map, got %s", rv.Kind())
	}
}

// GetParsedAuthenticationRequestData returns a map containing the authentication request data.
// Only includes fields that are defined and non-empty.
func (ar *AuthenticationRequest) GetParsedAuthenticationRequestData() map[string]interface{} {
	data := map[string]interface{}{
		"publisherUserId": ar.PublisherUserID,
	}
	optionalFields := map[string]interface{}{
		"age":               ar.Age,
		"gender":            ar.Gender,
		"email":             ar.Email,
		"phoneNumber":       ar.PhoneNumber,
		"sub1":              ar.Sub1,
		"sub2":              ar.Sub2,
		"sub3":              ar.Sub3,
		"sub4":              ar.Sub4,
		"sub5":              ar.Sub5,
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
		case *int:
			if v != nil {
				data[key] = *v
			}
		}
	}
	if encoded, ok := encodeUserGroup(ar.UserGroup); ok {
		data["userGroup"] = encoded
	}
	return data
}

// encodeUserGroup returns the wire form of userGroup. Strings are passed
// through as-is; structs and maps are JSON-encoded to a string. Returns
// ok=false to indicate the field should be omitted (nil or empty string).
func encodeUserGroup(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return "", false
		}
		return s, true
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", false
	}
	return string(b), true
}

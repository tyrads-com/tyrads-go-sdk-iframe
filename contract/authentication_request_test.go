package contract

import (
	"reflect"
	"strings"
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

func TestNewAuthenticationRequest_WithHelpers(t *testing.T) {
	userGroupMap := map[string]any{"tier": "vip", "region": "us"}

	req := NewAuthenticationRequest("user-helpers",
		WithAge(30),
		WithGender(2),
		WithEmail("a@b.co"),
		WithPhoneNumber("+15551234567"),
		WithSub1("s1"),
		WithSub2("s2"),
		WithSub3("s3"),
		WithSub4("s4"),
		WithSub5("s5"),
		WithUserGroup("group-x"),
		WithMediaSourceName("Facebook"),
		WithMediaSourceID("fb-1"),
		WithMediaSubSourceID("fb-sub-1"),
		WithIncentivized(true),
		WithMediaAdsetName("adset-name"),
		WithMediaAdsetID("adset-1"),
		WithMediaCreativeName("creative-name"),
		WithMediaCreativeID("creative-1"),
		WithMediaCampaignName("campaign-name"),
	)

	if req.PublisherUserID != "user-helpers" {
		t.Errorf("expected PublisherUserID 'user-helpers', got %q", req.PublisherUserID)
	}
	if req.Age == nil || *req.Age != 30 {
		t.Error("expected Age=30")
	}
	if req.Gender == nil || *req.Gender != 2 {
		t.Error("expected Gender=2")
	}
	if req.Email == nil || *req.Email != "a@b.co" {
		t.Error("expected Email=a@b.co")
	}
	if req.PhoneNumber == nil || *req.PhoneNumber != "+15551234567" {
		t.Error("expected PhoneNumber=+15551234567")
	}
	subPairs := map[string]*string{"s1": req.Sub1, "s2": req.Sub2, "s3": req.Sub3, "s4": req.Sub4, "s5": req.Sub5}
	for want, got := range subPairs {
		if got == nil || *got != want {
			t.Errorf("expected Sub%s pointer to %q, got %v", want[1:], want, got)
		}
	}
	if req.UserGroup != "group-x" {
		t.Errorf("expected UserGroup='group-x', got %v", req.UserGroup)
	}
	if req.MediaSourceName == nil || *req.MediaSourceName != "Facebook" {
		t.Error("expected MediaSourceName=Facebook")
	}
	if req.MediaSourceID == nil || *req.MediaSourceID != "fb-1" {
		t.Error("expected MediaSourceID=fb-1")
	}
	if req.MediaSubSourceID == nil || *req.MediaSubSourceID != "fb-sub-1" {
		t.Error("expected MediaSubSourceID=fb-sub-1")
	}
	if req.Incentivized == nil || *req.Incentivized != true {
		t.Error("expected Incentivized=true")
	}
	if req.MediaAdsetName == nil || *req.MediaAdsetName != "adset-name" {
		t.Error("expected MediaAdsetName=adset-name")
	}
	if req.MediaAdsetID == nil || *req.MediaAdsetID != "adset-1" {
		t.Error("expected MediaAdsetID=adset-1")
	}
	if req.MediaCreativeName == nil || *req.MediaCreativeName != "creative-name" {
		t.Error("expected MediaCreativeName=creative-name")
	}
	if req.MediaCreativeID == nil || *req.MediaCreativeID != "creative-1" {
		t.Error("expected MediaCreativeID=creative-1")
	}
	if req.MediaCampaignName == nil || *req.MediaCampaignName != "campaign-name" {
		t.Error("expected MediaCampaignName=campaign-name")
	}

	// userGroup as a map
	reqMap := NewAuthenticationRequest("u", WithUserGroup(userGroupMap))
	gotMap, ok := reqMap.UserGroup.(map[string]any)
	if !ok {
		t.Fatalf("expected UserGroup to be map[string]any, got %T", reqMap.UserGroup)
	}
	if gotMap["tier"] != "vip" || gotMap["region"] != "us" {
		t.Errorf("unexpected userGroup map content: %v", gotMap)
	}
}

func TestValidateAuthenticationRequest(t *testing.T) {
	validEmail := "test@example.com"
	invalidEmail := "invalid-email"
	validPhone := "+1234567890"
	invalidPhone := "invalid-phone"
	tooLong := strings.Repeat("a", maxOptionalStringLength+1)
	maxLen := strings.Repeat("a", maxOptionalStringLength)

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
		{
			name: "sub1 at max length is allowed",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Sub1:            &maxLen,
			},
			wantErr: false,
		},
		{
			name: "sub1 over max length rejected",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				Sub1:            &tooLong,
			},
			wantErr: true,
			errMsg:  "sub1 must not exceed 255 characters",
		},
		{
			name: "mediaCampaignName over max length rejected",
			request: &AuthenticationRequest{
				PublisherUserID:   "user123",
				MediaCampaignName: &tooLong,
			},
			wantErr: true,
			errMsg:  "mediaCampaignName must not exceed 255 characters",
		},
		{
			name: "userGroup as string ok",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       "group-1",
			},
			wantErr: false,
		},
		{
			name: "userGroup as map ok",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       map[string]any{"tier": "vip"},
			},
			wantErr: false,
		},
		{
			name: "userGroup as int rejected",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       42,
			},
			wantErr: true,
			errMsg:  "userGroup must be a string, struct, or map, got int",
		},
		{
			name: "userGroup as bool rejected",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       true,
			},
			wantErr: true,
			errMsg:  "userGroup must be a string, struct, or map, got bool",
		},
		{
			name: "userGroup as slice rejected",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       []string{"a", "b"},
			},
			wantErr: true,
			errMsg:  "userGroup must be a string, struct, or map, got slice",
		},
		{
			name: "userGroup string over max length rejected",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       strings.Repeat("x", maxUserGroupStringLength+1),
			},
			wantErr: true,
			errMsg:  "userGroup must not exceed 4096 characters",
		},
		{
			name: "userGroup map exceeding encoded length rejected",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup: map[string]any{
					"big": strings.Repeat("y", maxUserGroupStringLength),
				},
			},
			wantErr: true,
			errMsg:  "userGroup must not exceed 4096 characters when JSON-encoded",
		},
		{
			name: "userGroup as struct ok",
			request: &AuthenticationRequest{
				PublisherUserID: "user123",
				UserGroup:       struct{ Tier string }{Tier: "vip"},
			},
			wantErr: false,
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

func TestGetParsedAuthenticationRequestData_AllFields(t *testing.T) {
	req := NewAuthenticationRequest("user-all",
		WithAge(40),
		WithGender(1),
		WithEmail("a@b.co"),
		WithPhoneNumber("+10000000000"),
		WithSub1("s1"),
		WithSub2("s2"),
		WithSub3("s3"),
		WithSub4("s4"),
		WithSub5("s5"),
		WithUserGroup("plain-string"),
		WithMediaSourceName("MSN"),
		WithMediaSourceID("MSID"),
		WithMediaSubSourceID("MSSID"),
		WithIncentivized(false),
		WithMediaAdsetName("MAN"),
		WithMediaAdsetID("MAID"),
		WithMediaCreativeName("MCN"),
		WithMediaCreativeID("MCID"),
		WithMediaCampaignName("MCAN"),
	)

	expected := map[string]interface{}{
		"publisherUserId":   "user-all",
		"age":               40,
		"gender":            1,
		"email":             "a@b.co",
		"phoneNumber":       "+10000000000",
		"sub1":              "s1",
		"sub2":              "s2",
		"sub3":              "s3",
		"sub4":              "s4",
		"sub5":              "s5",
		"userGroup":         "plain-string",
		"mediaSourceName":   "MSN",
		"mediaSourceId":     "MSID",
		"mediaSubSourceId":  "MSSID",
		"incentivized":      false,
		"mediaAdsetName":    "MAN",
		"mediaAdsetId":      "MAID",
		"mediaCreativeName": "MCN",
		"mediaCreativeId":   "MCID",
		"mediaCampaignName": "MCAN",
	}

	got := req.GetParsedAuthenticationRequestData()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}

func TestGetParsedAuthenticationRequestData_UserGroupJSONEncoded(t *testing.T) {
	req := NewAuthenticationRequest("user-map",
		WithUserGroup(map[string]any{"tier": "vip"}),
	)

	data := req.GetParsedAuthenticationRequestData()

	got, ok := data["userGroup"].(string)
	if !ok {
		t.Fatalf("expected userGroup to be a string in parsed data, got %T", data["userGroup"])
	}
	if got != `{"tier":"vip"}` {
		t.Errorf("expected userGroup JSON-encoded as %q, got %q", `{"tier":"vip"}`, got)
	}
}

func TestGetParsedAuthenticationRequestData_UserGroupStruct(t *testing.T) {
	type group struct {
		Tier   string `json:"tier"`
		Region string `json:"region"`
	}
	req := NewAuthenticationRequest("user-struct",
		WithUserGroup(group{Tier: "vip", Region: "us"}),
	)

	data := req.GetParsedAuthenticationRequestData()

	got, ok := data["userGroup"].(string)
	if !ok {
		t.Fatalf("expected userGroup to be a string in parsed data, got %T", data["userGroup"])
	}
	if got != `{"tier":"vip","region":"us"}` {
		t.Errorf("expected userGroup JSON-encoded, got %q", got)
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
		UserGroup:       "",           // should be skipped
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

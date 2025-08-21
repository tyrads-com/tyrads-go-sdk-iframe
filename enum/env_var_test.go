package enum

import "testing"

func TestEnvVarConstants(t *testing.T) {
	tests := []struct {
		name     string
		envVar   EnvVar
		expected string
	}{
		{
			name:     "TYRADS_API_KEY constant",
			envVar:   TYRADS_API_KEY,
			expected: "TYRADS_API_KEY",
		},
		{
			name:     "TYRADS_API_SECRET constant",
			envVar:   TYRADS_API_SECRET,
			expected: "TYRADS_API_SECRET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.envVar) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(tt.envVar))
			}
		})
	}
}

package friendlycaptcha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAPIEndpoint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "global shorthand",
			input:    "global",
			expected: "https://global.frcapi.com",
		},
		{
			name:     "eu shorthand",
			input:    "eu",
			expected: "https://eu.frcapi.com",
		},
		{
			name:     "full domain https",
			input:    "https://custom.example.com",
			expected: "https://custom.example.com",
		},
		{
			name:     "full domain http",
			input:    "http://localhost:1090",
			expected: "http://localhost:1090",
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(
				WithAPIKey("test-key"),
				WithAPIEndpoint(tt.input),
			)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, client.APIEndpoint)
		})
	}
}

func TestWithSiteverifyEndpoint_Deprecated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:     "global shorthand",
			input:    "global",
			expected: "https://global.frcapi.com",
		},
		{
			name:     "eu shorthand",
			input:    "eu",
			expected: "https://eu.frcapi.com",
		},
		{
			name:     "full URL with path - strips path",
			input:    "https://global.frcapi.com/api/v2/captcha/siteverify",
			expected: "https://global.frcapi.com",
		},
		{
			name:     "full URL without path",
			input:    "https://custom.example.com",
			expected: "https://custom.example.com",
		},
		{
			name:     "localhost with path - strips path",
			input:    "http://localhost:1090/api/v2/captcha/siteverify",
			expected: "http://localhost:1090",
		},
		{
			name:     "localhost without path",
			input:    "http://localhost:1090",
			expected: "http://localhost:1090",
		},
		{
			name:     "https with port and path",
			input:    "https://example.com:8080/some/path",
			expected: "https://example.com:8080",
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(
				WithAPIKey("test-key"),
				WithSiteverifyEndpoint(tt.input),
			)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, client.APIEndpoint)
		})
	}
}

func TestWithSiteverifyEndpoint_BackwardCompatibility(t *testing.T) {
	t.Parallel()

	// Test that the old usage still works
	client, err := NewClient(
		WithAPIKey("test-key"),
		WithSiteverifyEndpoint("https://eu.frcapi.com/api/v2/captcha/siteverify"),
	)

	assert.NoError(t, err)
	assert.Equal(t, "https://eu.frcapi.com", client.APIEndpoint)
}

func TestWithAPIEndpoint_DefaultValue(t *testing.T) {
	// Test that the default is set correctly
	client, err := NewClient(
		WithAPIKey("test-key"),
	)

	assert.NoError(t, err)
	assert.Equal(t, "https://global.frcapi.com", client.APIEndpoint)
}

func TestWithAPIEndpoint_OverridesWithSiteverifyEndpoint(t *testing.T) {
	t.Parallel()

	// Test that when both are provided, the last one wins
	client, err := NewClient(
		WithAPIKey("test-key"),
		WithSiteverifyEndpoint("https://global.frcapi.com/api/v2/captcha/siteverify"),
		WithAPIEndpoint("eu"),
	)

	assert.NoError(t, err)
	assert.Equal(t, "https://eu.frcapi.com", client.APIEndpoint)
}

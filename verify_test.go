package friendlycaptcha

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyCaptchaResponse_TooLongResponseIsTruncatedAndSent(t *testing.T) {
	t.Parallel()

	var sent VerifyRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&sent))
		_, _ = w.Write([]byte(`{"success":false,"error":{"error_code":"response_invalid"}}`))
	}))
	defer server.Close()

	client, err := NewClient(WithAPIKey("some-key"), WithAPIEndpoint(server.URL))
	require.NoError(t, err)

	result := client.VerifyCaptchaResponse(context.Background(), strings.Repeat("A", 3*1024*1024))
	assert.Len(t, sent.Response, 2*1024*1024)
	assert.False(t, result.ShouldAccept())
}

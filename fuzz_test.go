package loops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// FuzzClientResponse feeds random response bodies to the client. Ensures we never panic
// on malformed, empty, or arbitrary server responses (HTML errors, truncated JSON, etc.).
func FuzzClientResponse(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("null"))
	f.Add([]byte("{}"))
	f.Add([]byte("not json"))
	f.Add([]byte("<html>error</html>"))
	f.Add([]byte(`{"success":true,"teamName":""}`))
	f.Add([]byte(`{"success":false,"message":"bad"}`))
	f.Fuzz(func(t *testing.T, body []byte) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}))
		defer server.Close()
		client := NewClient("key", WithBaseURL(server.URL))
		ctx := context.Background()
		// Must not panic; we only care about robustness to arbitrary response bodies
		_, _ = client.GetAPIKey(ctx)
	})
}

// FuzzClientErrorResponse feeds random bodies for 4xx/5xx responses. Ensures error
// parsing never panics and always returns an APIError with correct status.
func FuzzClientErrorResponse(f *testing.F) {
	f.Add([]byte(""), 400)
	f.Add([]byte("not json"), 401)
	f.Add([]byte("<html>500</html>"), 500)
	f.Add([]byte(`{"success":false,"message":"x"}`), 409)
	f.Fuzz(func(t *testing.T, body []byte, status int) {
		if status < 400 || status > 599 {
			return
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			w.Write(body)
		}))
		defer server.Close()
		client := NewClient("key", WithBaseURL(server.URL))
		ctx := context.Background()
		_, err := client.GetAPIKey(ctx)
		if err == nil {
			return
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Errorf("expected *APIError for status %d, got %T", status, err)
			return
		}
		if apiErr.StatusCode != status {
			t.Errorf("StatusCode: got %d, want %d", apiErr.StatusCode, status)
		}
	})
}

// FuzzMergeBody feeds random JSON bytes and ensures mergeBody (via public API paths
// that use it) never panics. We test by building request structs with fuzz-derived
// extra data and calling mergeBody.
func FuzzMergeBody(f *testing.F) {
	f.Add([]byte("{}"))
	f.Add([]byte(`{"x":1}`))
	f.Add([]byte(`{"key":"value"}`))
	f.Add([]byte("null"))
	f.Add([]byte("not json"))
	f.Fuzz(func(t *testing.T, extraJSON []byte) {
		var extra map[string]interface{}
		_ = json.Unmarshal(extraJSON, &extra)
		req := &ContactRequest{Email: "a@b.com"}
		_, _ = mergeBody(req, extra)
	})
}

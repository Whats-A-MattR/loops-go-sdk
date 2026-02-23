package loops

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Negative / anti-tests: malformed responses, error statuses, edge cases.
// These ensure we don't assume happy path and handle failures safely.

func TestClient_200_MalformedJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.GetAPIKey(ctx)
	if err == nil {
		t.Fatal("expected error for malformed JSON 200 response")
	}
}

func TestClient_200_EmptyBody_DoesNotPanic_ResultZeroValued(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.GetAPIKey(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Current behavior: empty body yields zero-valued struct (no error)
	if got.Success != false || got.TeamName != "" {
		t.Errorf("expected zero-valued result, got %+v", got)
	}
}

func TestClient_200_HTMLBody_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html><html><body>Error</body></html>`))
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.GetAPIKey(ctx)
	if err == nil {
		t.Fatal("expected error for HTML body")
	}
}

func TestClient_400_EmptyBody_ReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(400)
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.GetAPIKey(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("StatusCode: got %d, want 400", apiErr.StatusCode)
	}
}

func TestClient_401_MalformedJSON_ReturnsAPIErrorWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte("<html>Unauthorized</html>"))
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.GetAPIKey(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode: got %d", apiErr.StatusCode)
	}
	if len(apiErr.Body) == 0 {
		t.Error("Body should be preserved for non-JSON error response")
	}
}

func TestClient_404_ReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"success":false,"message":"Contact not found."}`))
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.CreateContact(ctx, &ContactRequest{Email: "x@y.com"})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 || apiErr.Message != "Contact not found." {
		t.Errorf("got StatusCode=%d Message=%q", apiErr.StatusCode, apiErr.Message)
	}
}

func TestClient_500_ServerError_ReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"success":false,"message":"Internal server error"}`))
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.GetAPIKey(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode: got %d", apiErr.StatusCode)
	}
}

func TestClient_409_IdempotencyKeyUsed_ReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(409)
		w.Write([]byte(`{"success":false,"message":"Idempotency key has been used."}`))
	}))
	defer server.Close()
	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.SendEvent(ctx, &EventRequest{EventName: "e", Email: "a@b.com"}, "dup-key")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 409 {
		t.Errorf("StatusCode: got %d", apiErr.StatusCode)
	}
}

func TestClient_DeleteContact_BothEmailAndUserIDSet_ReturnsErrorBeforeRequest(t *testing.T) {
	// Validation should run before any HTTP call
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.DeleteContact(ctx, &ContactDeleteRequest{Email: "a@b.com", UserID: "u1"})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if apiErr, ok := err.(*APIError); !ok || apiErr.StatusCode != 400 {
		t.Errorf("expected 400 APIError, got %v", err)
	}
}

func TestClient_FindContact_BothParamsEmpty_ReturnsErrorBeforeRequest(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.FindContact(ctx, "", "")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestClient_SendEvent_EventNameMissing_ReturnsErrorBeforeRequest(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.SendEvent(ctx, &EventRequest{Email: "a@b.com"}, "")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestClient_SendTransactional_EmptyTransactionalID_ReturnsErrorBeforeRequest(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.SendTransactional(ctx, &TransactionalRequest{Email: "a@b.com"}, "")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestClient_CreateContactProperty_EmptyName_ReturnsErrorBeforeRequest(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.CreateContactProperty(ctx, &ContactPropertyCreateRequest{Type: "string"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestAPIError_ErrorString_NoPanicWhenBodyNil(t *testing.T) {
	e := &APIError{StatusCode: 500, Body: nil}
	_ = e.Error()
}

func TestAPIError_ErrorString_NoPanicWhenBodyLarge(t *testing.T) {
	body := make([]byte, 1e6)
	for i := range body {
		body[i] = 'x'
	}
	e := &APIError{StatusCode: 500, Body: body}
	s := e.Error()
	if len(s) == 0 {
		t.Error("Error() should return non-empty string")
	}
}

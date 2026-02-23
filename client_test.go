package loops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// specServer returns an httptest.Server that records the last request and responds with the given status and body.
// Used to assert request shape (method, path, headers, body) per OpenAPI spec.
func specServer(t *testing.T, status int, body []byte) (*httptest.Server, *http.Request) {
	var lastReq *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastReq = r
		w.WriteHeader(status)
		if body != nil {
			w.Write(body)
		}
	}))
	t.Cleanup(server.Close)
	return server, lastReq
}

func TestClient_GetAPIKey_SpecCompliant(t *testing.T) {
	resp := APIKeyResponse{Success: true, TeamName: "Acme"}
	body, _ := json.Marshal(resp)
	server, _ := specServer(t, 200, body)

	client := NewClient("test-key", WithBaseURL(server.URL))
	ctx := context.Background()

	got, err := client.GetAPIKey(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Success || got.TeamName != "Acme" {
		t.Errorf("got %+v", got)
	}

	var captured *http.Request
	s2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(s2.Close)
	client2 := NewClient("secret", WithBaseURL(s2.URL))
	_, _ = client2.GetAPIKey(ctx)
	if captured == nil {
		t.Fatal("request not captured")
	}
	if captured.Method != http.MethodGet {
		t.Errorf("method: got %s, want GET (per OpenAPI)", captured.Method)
	}
	if captured.URL.Path != "/api-key" {
		t.Errorf("path: got %s, want /api-key (per OpenAPI)", captured.URL.Path)
	}
	if auth := captured.Header.Get("Authorization"); auth != "Bearer secret" {
		t.Errorf("Authorization: got %q, want Bearer secret (OpenAPI securitySchemes.apiKey bearer)", auth)
	}
}

func TestClient_CreateContact_SpecCompliant(t *testing.T) {
	resp := ContactSuccessResponse{Success: true, ID: "contact_123"}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	var reqBodyBytes []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		reqBodyBytes, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.CreateContact(ctx, &ContactRequest{Email: "u@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Success || got.ID != "contact_123" {
		t.Errorf("got %+v", got)
	}
	if captured.Method != http.MethodPost {
		t.Errorf("method: got %s, want POST", captured.Method)
	}
	if captured.URL.Path != "/contacts/create" {
		t.Errorf("path: got %s, want /contacts/create", captured.URL.Path)
	}
	var reqBody struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(reqBodyBytes, &reqBody); err != nil {
		t.Fatal(err)
	}
	if reqBody.Email != "u@example.com" {
		t.Errorf("body email: got %q", reqBody.Email)
	}
}

func TestClient_CreateContact_EmailRequired(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.CreateContact(ctx, nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	_, err = client.CreateContact(ctx, &ContactRequest{})
	if err == nil {
		t.Fatal("expected error for missing email")
	}
	if apiErr, ok := err.(*APIError); !ok || apiErr.StatusCode != 400 {
		t.Errorf("expected APIError 400, got %v", err)
	}
}

func TestClient_UpdateContact_SpecCompliant(t *testing.T) {
	resp := ContactSuccessResponse{Success: true, ID: "contact_123"}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.UpdateContact(ctx, &ContactUpdateRequest{Email: "u@example.com", FirstName: "Jane"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "contact_123" {
		t.Errorf("got %+v", got)
	}
	if captured.URL.Path != "/contacts/update" || captured.Method != http.MethodPut {
		t.Errorf("path=%s method=%s (OpenAPI: PUT /contacts/update)", captured.URL.Path, captured.Method)
	}
}

func TestClient_UpdateContact_EmailOrUserIDRequired(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.UpdateContact(ctx, &ContactUpdateRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if apiErr, ok := err.(*APIError); !ok || apiErr.StatusCode != 400 {
		t.Errorf("got %v", err)
	}
}

func TestClient_FindContact_SpecCompliant(t *testing.T) {
	contacts := []Contact{{ID: "c1", Email: "u@example.com"}}
	body, _ := json.Marshal(contacts)
	var captured *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.FindContact(ctx, "u@example.com", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Email != "u@example.com" {
		t.Errorf("got %+v", got)
	}
	if captured.URL.Path != "/contacts/find" {
		t.Errorf("path: %s", captured.URL.Path)
	}
	if captured.URL.Query().Get("email") != "u@example.com" || captured.URL.Query().Get("userId") != "" {
		t.Errorf("query: email=%q userId=%q (only one allowed per OpenAPI)", captured.URL.Query().Get("email"), captured.URL.Query().Get("userId"))
	}
}

func TestClient_FindContact_ExactlyOneParam(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.FindContact(ctx, "", "")
	if err == nil {
		t.Fatal("expected error when both empty")
	}
	_, err = client.FindContact(ctx, "a@b.com", "user1")
	if err == nil {
		t.Fatal("expected error when both set")
	}
}

func TestClient_DeleteContact_SpecCompliant(t *testing.T) {
	resp := ContactDeleteResponse{Success: true, Message: "Contact deleted."}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	var reqBodyBytes []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		reqBodyBytes, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.DeleteContact(ctx, &ContactDeleteRequest{Email: "u@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Success || got.Message != "Contact deleted." {
		t.Errorf("got %+v", got)
	}
	if captured.Method != http.MethodPost || captured.URL.Path != "/contacts/delete" {
		t.Errorf("method=%s path=%s", captured.Method, captured.URL.Path)
	}
	var reqBody ContactDeleteRequest
	json.Unmarshal(reqBodyBytes, &reqBody)
	if reqBody.Email != "u@example.com" {
		t.Errorf("body email: got %q", reqBody.Email)
	}
}

func TestClient_DeleteContact_ExactlyOneRequired(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.DeleteContact(ctx, &ContactDeleteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = client.DeleteContact(ctx, &ContactDeleteRequest{Email: "a@b.com", UserID: "u1"})
	if err == nil {
		t.Fatal("expected error when both set")
	}
}

func TestClient_SendEvent_SpecCompliant(t *testing.T) {
	resp := EventSuccessResponse{Success: true}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	var reqBodyBytes []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		reqBodyBytes, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.SendEvent(ctx, &EventRequest{EventName: "signup", Email: "u@example.com"}, "idem-123")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Success {
		t.Errorf("got %+v", got)
	}
	if captured.URL.Path != "/events/send" || captured.Method != http.MethodPost {
		t.Errorf("path=%s method=%s", captured.URL.Path, captured.Method)
	}
	if key := captured.Header.Get("Idempotency-Key"); key != "idem-123" {
		t.Errorf("Idempotency-Key: got %q (OpenAPI optional header)", key)
	}
	var reqBody struct {
		EventName string `json:"eventName"`
		Email     string `json:"email"`
	}
	json.Unmarshal(reqBodyBytes, &reqBody)
	if reqBody.EventName != "signup" || reqBody.Email != "u@example.com" {
		t.Errorf("body: %+v", reqBody)
	}
}

func TestClient_SendEvent_EventNameAndIdentifierRequired(t *testing.T) {
	client := NewClient("key")
	ctx := context.Background()
	_, err := client.SendEvent(ctx, &EventRequest{Email: "a@b.com"}, "")
	if err == nil {
		t.Fatal("expected error when eventName missing")
	}
	_, err = client.SendEvent(ctx, &EventRequest{EventName: "e"}, "")
	if err == nil {
		t.Fatal("expected error when email and userId missing")
	}
}

func TestClient_SendTransactional_SpecCompliant(t *testing.T) {
	resp := TransactionalSuccessResponse{Success: true}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	var reqBodyBytes []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		reqBodyBytes, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.SendTransactional(ctx, &TransactionalRequest{
		Email:           "u@example.com",
		TransactionalID: "tx_abc",
	}, "idem-456")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Success {
		t.Errorf("got %+v", got)
	}
	if captured.URL.Path != "/transactional" || captured.Method != http.MethodPost {
		t.Errorf("path=%s method=%s", captured.URL.Path, captured.Method)
	}
	if key := captured.Header.Get("Idempotency-Key"); key != "idem-456" {
		t.Errorf("Idempotency-Key: got %q", key)
	}
	var reqBody TransactionalRequest
	json.Unmarshal(reqBodyBytes, &reqBody)
	if reqBody.Email != "u@example.com" || reqBody.TransactionalID != "tx_abc" {
		t.Errorf("body: %+v", reqBody)
	}
}

func TestClient_ListTransactionals_SpecCompliant(t *testing.T) {
	resp := ListTransactionalsResponse{
		Pagination: ListTransactionalsPagination{TotalResults: 1, ReturnedResults: 1, PerPage: 20, TotalPages: 1},
		Data:       []TransactionalEmail{{ID: "tx1", Name: "Welcome", LastUpdated: "2025-01-01", DataVariables: []string{"name"}}},
	}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.ListTransactionals(ctx, 25, "cursor_xyz")
	if err != nil {
		t.Fatal(err)
	}
	if got.Pagination.PerPage != 20 || len(got.Data) != 1 || got.Data[0].ID != "tx1" {
		t.Errorf("got %+v", got)
	}
	if captured.Method != http.MethodGet || captured.URL.Path != "/transactional" {
		t.Errorf("method=%s path=%s", captured.Method, captured.URL.Path)
	}
	q := captured.URL.Query()
	if q.Get("perPage") != "25" || q.Get("cursor") != "cursor_xyz" {
		t.Errorf("query: %s", captured.URL.RawQuery)
	}
}

func TestClient_GetLists_SpecCompliant(t *testing.T) {
	resp := []MailingList{{ID: "list_1", Name: "Main", Description: "Desc", IsPublic: true}}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.GetLists(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "list_1" {
		t.Errorf("got %+v", got)
	}
	if captured.Method != http.MethodGet || captured.URL.Path != "/lists" {
		t.Errorf("path=%s", captured.URL.Path)
	}
}

func TestClient_GetDedicatedSendingIPs_SpecCompliant(t *testing.T) {
	resp := []string{"1.2.3.4", "5.6.7.8"}
	body, _ := json.Marshal(resp)
	var captured *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(200)
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	client := NewClient("key", WithBaseURL(server.URL))
	ctx := context.Background()
	got, err := client.GetDedicatedSendingIPs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "1.2.3.4" {
		t.Errorf("got %+v", got)
	}
	if captured.Method != http.MethodGet || captured.URL.Path != "/dedicated-sending-ips" {
		t.Errorf("path=%s", captured.URL.Path)
	}
}

func TestClient_ContactProperties_SpecCompliant(t *testing.T) {
	createBody, _ := json.Marshal(ContactPropertySuccessResponse{Success: true})
	var createReq *http.Request
	var createBodyBytes []byte
	s1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createReq = r
		createBodyBytes, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write(createBody)
	}))
	t.Cleanup(s1.Close)
	client := NewClient("key", WithBaseURL(s1.URL))
	ctx := context.Background()
	_, err := client.CreateContactProperty(ctx, &ContactPropertyCreateRequest{Name: "planName", Type: "string"})
	if err != nil {
		t.Fatal(err)
	}
	if createReq.URL.Path != "/contacts/properties" || createReq.Method != http.MethodPost {
		t.Errorf("create: path=%s method=%s", createReq.URL.Path, createReq.Method)
	}
	var createPayload ContactPropertyCreateRequest
	json.Unmarshal(createBodyBytes, &createPayload)
	if createPayload.Name != "planName" || createPayload.Type != "string" {
		t.Errorf("create body: %+v", createPayload)
	}

	listBody, _ := json.Marshal([]ContactProperty{{Key: "planName", Label: "Plan", Type: "string"}})
	var listReq *http.Request
	s2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listReq = r
		w.WriteHeader(200)
		w.Write(listBody)
	}))
	t.Cleanup(s2.Close)
	client2 := NewClient("key", WithBaseURL(s2.URL))
	got, err := client2.ListContactProperties(ctx, "custom")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Key != "planName" {
		t.Errorf("list: %+v", got)
	}
	if listReq.URL.Path != "/contacts/properties" {
		t.Errorf("list path: %s", listReq.URL.Path)
	}
	if listReq.URL.Query().Get("list") != "custom" {
		t.Errorf("list query: %s", listReq.URL.RawQuery)
	}
}

func TestAPIError_ResponseParsing(t *testing.T) {
	body := []byte(`{"success":false,"message":"Invalid API key"}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write(body)
	}))
	t.Cleanup(server.Close)
	client := NewClient("bad-key", WithBaseURL(server.URL))
	ctx := context.Background()
	_, err := client.GetAPIKey(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 || apiErr.Message != "Invalid API key" || apiErr.Success != false {
		t.Errorf("got %+v", apiErr)
	}
	if !strings.Contains(apiErr.Error(), "401") || !strings.Contains(apiErr.Error(), "Invalid API key") {
		t.Errorf("Error() = %q", apiErr.Error())
	}
}

func TestNewClient_BaseURL(t *testing.T) {
	c := NewClient("k", WithBaseURL("https://custom.example.com/api/v1"))
	if c.baseURL != "https://custom.example.com/api/v1" {
		t.Errorf("baseURL: got %s", c.baseURL)
	}
	c2 := NewClient("k", WithBaseURL("https://example.com/"))
	if c2.baseURL != "https://example.com" {
		t.Errorf("trailing slash should be trimmed: got %s", c2.baseURL)
	}
}

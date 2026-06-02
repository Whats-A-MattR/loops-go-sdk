// Tests that the SDK's endpoints match the Loops OpenAPI spec (openapi.json).
package loops

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type specPaths map[string]map[string]interface{}

func loadSpecPaths(t *testing.T) specPaths {
	t.Helper()
	dir := "."
	b, err := os.ReadFile(filepath.Join(dir, "openapi.json"))
	if err != nil {
		t.Skipf("openapi.json not found (run from repo root): %v", err)
		return nil
	}
	var spec struct {
		Paths specPaths `json:"paths"`
	}
	if err := json.Unmarshal(b, &spec); err != nil {
		t.Fatalf("parse openapi.json: %v", err)
	}
	return spec.Paths
}

var expectedEndpoints = []struct {
	Method string
	Path   string
}{
	{"GET", "/v1/api-key"},
	{"GET", "/v1/campaigns"},
	{"POST", "/v1/campaigns"},
	{"GET", "/v1/campaigns/{campaignId}"},
	{"POST", "/v1/campaigns/{campaignId}"},
	{"GET", "/v1/components"},
	{"GET", "/v1/components/{componentId}"},
	{"POST", "/v1/contacts/create"},
	{"PUT", "/v1/contacts/update"},
	{"GET", "/v1/contacts/find"},
	{"POST", "/v1/contacts/delete"},
	{"POST", "/v1/contacts/properties"},
	{"GET", "/v1/contacts/properties"},
	{"GET", "/v1/contacts/suppression"},
	{"DELETE", "/v1/contacts/suppression"},
	{"GET", "/v1/dedicated-sending-ips"},
	{"GET", "/v1/email-messages/{emailMessageId}"},
	{"POST", "/v1/email-messages/{emailMessageId}"},
	{"GET", "/v1/lists"},
	{"GET", "/v1/themes"},
	{"GET", "/v1/themes/{themeId}"},
	{"POST", "/v1/events/send"},
	{"POST", "/v1/transactional"},
	{"GET", "/v1/transactional"},
	{"GET", "/v2/transactional"},
	{"POST", "/v2/transactional"},
	{"GET", "/v2/transactional/{transactionalId}"},
	{"POST", "/v2/transactional/{transactionalId}"},
	{"POST", "/v2/transactional/{transactionalId}/draft"},
	{"POST", "/v2/transactional/{transactionalId}/publish"},
	{"POST", "/v1/uploads"},
	{"POST", "/v1/uploads/{id}/complete"},
}

func TestOpenAPI_SDKEndpointsExistInSpec(t *testing.T) {
	paths := loadSpecPaths(t)
	if paths == nil {
		return
	}
	for _, ep := range expectedEndpoints {
		pathOps, ok := paths[ep.Path]
		if !ok {
			t.Errorf("path %q not in OpenAPI spec", ep.Path)
			continue
		}
		methodLower := strings.ToLower(ep.Method)
		if _, ok := pathOps[methodLower]; !ok {
			t.Errorf("method %s for path %q not in OpenAPI spec (allowed: %v)", ep.Method, ep.Path, keys(pathOps))
		}
	}
}

func TestOpenAPI_UploadMaxContentLength(t *testing.T) {
	dir := "."
	b, err := os.ReadFile(filepath.Join(dir, "openapi.json"))
	if err != nil {
		t.Skipf("openapi.json not found: %v", err)
	}
	var spec struct {
		Components struct {
			Schemas map[string]struct {
				Properties map[string]struct {
					Description string `json:"description"`
				} `json:"properties"`
			} `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(b, &spec); err != nil {
		t.Fatalf("parse openapi.json: %v", err)
	}
	schema, ok := spec.Components.Schemas["CreateUploadRequest"]
	if !ok {
		t.Fatal("CreateUploadRequest schema not found in spec")
	}
	prop, ok := schema.Properties["contentLength"]
	if !ok {
		t.Fatal("contentLength property not found in CreateUploadRequest")
	}
	if !strings.Contains(prop.Description, "4,000,000") {
		t.Errorf("spec contentLength description changed — SDK max-bytes validation may need updating; description: %q", prop.Description)
	}
}

func keys(m map[string]interface{}) []string {
	var k []string
	for s := range m {
		k = append(k, s)
	}
	return k
}

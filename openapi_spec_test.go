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
	{"GET", "/api-key"},
	{"POST", "/contacts/create"},
	{"PUT", "/contacts/update"},
	{"GET", "/contacts/find"},
	{"POST", "/contacts/delete"},
	{"POST", "/contacts/properties"},
	{"GET", "/contacts/properties"},
	{"GET", "/dedicated-sending-ips"},
	{"GET", "/lists"},
	{"POST", "/events/send"},
	{"POST", "/transactional"},
	{"GET", "/transactional"},
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

func keys(m map[string]interface{}) []string {
	var k []string
	for s := range m {
		k = append(k, s)
	}
	return k
}

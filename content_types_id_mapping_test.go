package loops

import (
	"encoding/json"
	"testing"
)

func TestContentTypes_IDMappings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		check func(t *testing.T, b []byte)
	}{
		{
			name: "Theme",
			body: `{"id":"theme_123"}`,
			check: func(t *testing.T, b []byte) {
				var v Theme
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.ThemeID != "theme_123" {
					t.Fatalf("ThemeID = %q, want %q", v.ThemeID, "theme_123")
				}
			},
		},
		{
			name: "ThemeResponse",
			body: `{"id":"theme_123"}`,
			check: func(t *testing.T, b []byte) {
				var v ThemeResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.ThemeID != "theme_123" {
					t.Fatalf("ThemeID = %q, want %q", v.ThemeID, "theme_123")
				}
			},
		},
		{
			name: "Component",
			body: `{"id":"component_123"}`,
			check: func(t *testing.T, b []byte) {
				var v Component
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.ComponentID != "component_123" {
					t.Fatalf("ComponentID = %q, want %q", v.ComponentID, "component_123")
				}
			},
		},
		{
			name: "ComponentResponse",
			body: `{"id":"component_123"}`,
			check: func(t *testing.T, b []byte) {
				var v ComponentResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.ComponentID != "component_123" {
					t.Fatalf("ComponentID = %q, want %q", v.ComponentID, "component_123")
				}
			},
		},
		{
			name: "CampaignListItem",
			body: `{"id":"campaign_123","emailMessageId":"email_message_123"}`,
			check: func(t *testing.T, b []byte) {
				var v CampaignListItem
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.CampaignID != "campaign_123" {
					t.Fatalf("CampaignID = %q, want %q", v.CampaignID, "campaign_123")
				}
				if v.EmailMessageID == nil || *v.EmailMessageID != "email_message_123" {
					t.Fatalf("EmailMessageID = %v, want %q", v.EmailMessageID, "email_message_123")
				}
			},
		},
		{
			name: "CreateCampaignResponse",
			body: `{"id":"campaign_123"}`,
			check: func(t *testing.T, b []byte) {
				var v CreateCampaignResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.CampaignID != "campaign_123" {
					t.Fatalf("CampaignID = %q, want %q", v.CampaignID, "campaign_123")
				}
			},
		},
		{
			name: "CampaignResponse",
			body: `{"id":"campaign_123","emailMessageId":"email_message_123"}`,
			check: func(t *testing.T, b []byte) {
				var v CampaignResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.CampaignID != "campaign_123" {
					t.Fatalf("CampaignID = %q, want %q", v.CampaignID, "campaign_123")
				}
				if v.EmailMessageID == nil || *v.EmailMessageID != "email_message_123" {
					t.Fatalf("EmailMessageID = %v, want %q", v.EmailMessageID, "email_message_123")
				}
			},
		},
		{
			name: "EmailMessageResponse",
			body: `{"id":"email_message_123","campaignId":"campaign_123"}`,
			check: func(t *testing.T, b []byte) {
				var v EmailMessageResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.EmailMessageID != "email_message_123" {
					t.Fatalf("EmailMessageID = %q, want %q", v.EmailMessageID, "email_message_123")
				}
				if v.CampaignID == nil || *v.CampaignID != "campaign_123" {
					t.Fatalf("CampaignID = %v, want %q", v.CampaignID, "campaign_123")
				}
			},
		},
		{
			name: "TransactionalEmailResource",
			body: `{"id":"transactional_123","name":"example","dataVariables":[]}`,
			check: func(t *testing.T, b []byte) {
				var v TransactionalEmailResource
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.TransactionalID != "transactional_123" {
					t.Fatalf("TransactionalID = %q, want %q", v.TransactionalID, "transactional_123")
				}
			},
		},
		{
			name: "TransactionalResourceResponse",
			body: `{"id":"transactional_123","name":"example","dataVariables":[]}`,
			check: func(t *testing.T, b []byte) {
				var v TransactionalResourceResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.TransactionalID != "transactional_123" {
					t.Fatalf("TransactionalID = %q, want %q", v.TransactionalID, "transactional_123")
				}
			},
		},
		{
			name: "TransactionalDraftResponse",
			body: `{"id":"transactional_123","name":"example","dataVariables":[]}`,
			check: func(t *testing.T, b []byte) {
				var v TransactionalDraftResponse
				if err := json.Unmarshal(b, &v); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if v.TransactionalID != "transactional_123" {
					t.Fatalf("TransactionalID = %q, want %q", v.TransactionalID, "transactional_123")
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.check(t, []byte(tt.body))
		})
	}
}

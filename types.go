// Package loops provides a server-side Go SDK for the Loops API.
// Types and request/response structs match the Loops OpenAPI spec (https://app.loops.so/openapi.json).
package loops

// DefaultBaseURL is the default Loops API base URL from the OpenAPI spec.
const DefaultBaseURL = "https://app.loops.so/api/v1"

// --- API key (GET /api-key) ---

// APIKeyResponse is the 200 response for GET /api-key (OpenAPI: success + teamName required).
type APIKeyResponse struct {
	Success  bool   `json:"success"`
	TeamName string `json:"teamName"`
}

// APIKeyErrorResponse is the 401 response for GET /api-key.
type APIKeyErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// --- Contact (schema) ---

// Contact represents a contact (OpenAPI schema Contact).
type Contact struct {
	ID           string          `json:"id,omitempty"`
	Email        string          `json:"email,omitempty"`
	FirstName    *string         `json:"firstName,omitempty"`
	LastName     *string         `json:"lastName,omitempty"`
	Source       string          `json:"source,omitempty"`
	Subscribed   bool            `json:"subscribed"`
	UserGroup    string          `json:"userGroup,omitempty"`
	UserID       *string         `json:"userId,omitempty"`
	MailingLists map[string]bool `json:"mailingLists,omitempty"`
	OptInStatus  *string         `json:"optInStatus,omitempty"` // "accepted" | "pending" | "rejected"
}

// --- Contact create/update (ContactRequest, ContactUpdateRequest) ---

// ContactRequest is the body for POST /contacts/create (email required; OpenAPI ContactRequest).
type ContactRequest struct {
	Email        string                 `json:"email"`
	FirstName    string                 `json:"firstName,omitempty"`
	LastName     string                 `json:"lastName,omitempty"`
	Subscribed   *bool                  `json:"subscribed,omitempty"`
	UserGroup    string                 `json:"userGroup,omitempty"`
	UserID       string                 `json:"userId,omitempty"`
	MailingLists map[string]bool        `json:"mailingLists,omitempty"`
	Extra        map[string]interface{} `json:"-"`
}

// ContactUpdateRequest is the body for PUT /contacts/update (email or userId; OpenAPI ContactUpdateRequest).
type ContactUpdateRequest struct {
	Email        string                 `json:"email,omitempty"`
	FirstName    string                 `json:"firstName,omitempty"`
	LastName     string                 `json:"lastName,omitempty"`
	Subscribed   *bool                  `json:"subscribed,omitempty"`
	UserGroup    string                 `json:"userGroup,omitempty"`
	UserID       string                 `json:"userId,omitempty"`
	MailingLists map[string]bool        `json:"mailingLists,omitempty"`
	Extra        map[string]interface{} `json:"-"`
}

// ContactSuccessResponse is the 200 response for create/update (OpenAPI: success, id required).
type ContactSuccessResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// ContactFailureResponse is used for 400/404/409 (OpenAPI: success, message required).
type ContactFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ContactDeleteRequest is the body for POST /contacts/delete. Spec says include only one of email or userId.
type ContactDeleteRequest struct {
	Email  string `json:"email,omitempty"`
	UserID string `json:"userId,omitempty"`
}

// ContactDeleteResponse is the 200 response for delete (OpenAPI: success, message required).
type ContactDeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Contact properties ---

// ContactPropertyCreateRequest is the body for POST /contacts/properties (OpenAPI: name, type required).
type ContactPropertyCreateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ContactProperty is a single property (OpenAPI: key, label, type required).
type ContactProperty struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// ContactPropertySuccessResponse is the 200 response for create property.
type ContactPropertySuccessResponse struct {
	Success bool `json:"success"`
}

// ContactPropertyFailureResponse is the 400 response (OpenAPI: success, message required).
type ContactPropertyFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Mailing lists ---

// MailingList (OpenAPI schema: id, name, description, isPublic required).
type MailingList struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

// --- Events ---

// EventRequest is the body for POST /events/send (eventName required; OpenAPI EventRequest).
type EventRequest struct {
	Email           string                 `json:"email,omitempty"`
	UserID          string                 `json:"userId,omitempty"`
	EventName       string                 `json:"eventName"`
	EventProperties map[string]interface{} `json:"eventProperties,omitempty"`
	MailingLists    map[string]bool        `json:"mailingLists,omitempty"`
	Extra           map[string]interface{} `json:"-"`
}

// EventSuccessResponse is the 200 response (OpenAPI: success required).
type EventSuccessResponse struct {
	Success bool `json:"success"`
}

// EventFailureResponse is the 400 response (OpenAPI: success, message required).
type EventFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// IdempotencyKeyFailureResponse is the 409 response for duplicate Idempotency-Key (OpenAPI: success, message required).
type IdempotencyKeyFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Transactional ---

// TransactionalRequest is the body for POST /transactional (OpenAPI: email, transactionalId required).
type TransactionalRequest struct {
	Email           string                    `json:"email"`
	TransactionalID string                    `json:"transactionalId"`
	AddToAudience   *bool                     `json:"addToAudience,omitempty"`
	DataVariables   map[string]interface{}    `json:"dataVariables,omitempty"`
	Attachments     []TransactionalAttachment `json:"attachments,omitempty"`
}

// TransactionalAttachment (OpenAPI: filename, contentType, data required).
type TransactionalAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Data        string `json:"data"` // base64-encoded
}

// TransactionalSuccessResponse is the 200 response (OpenAPI: success required).
type TransactionalSuccessResponse struct {
	Success bool `json:"success"`
}

// TransactionalFailureResponse is used for 400/404 (OpenAPI variants have success, message; some add path, error, transactionalId).
type TransactionalFailureResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	Path            string `json:"path,omitempty"`
	Error           *struct {
		Path    string `json:"path,omitempty"`
		Message string `json:"message,omitempty"`
		Reason  string `json:"reason,omitempty"`
	} `json:"error,omitempty"`
	TransactionalID string `json:"transactionalId,omitempty"`
}

// TransactionalEmail (OpenAPI: id, name, lastUpdated, dataVariables required).
type TransactionalEmail struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	LastUpdated    string   `json:"lastUpdated"`
	DataVariables  []string `json:"dataVariables"`
}

// ListTransactionalsPagination (OpenAPI ListTransactionalsResponse.pagination).
type ListTransactionalsPagination struct {
	TotalResults    int     `json:"totalResults"`
	ReturnedResults int     `json:"returnedResults"`
	PerPage         int     `json:"perPage"`
	TotalPages      int     `json:"totalPages"`
	NextCursor      *string `json:"nextCursor"`
	NextPage        *string `json:"nextPage"`
}

// ListTransactionalsResponse is the 200 response for GET /transactional.
type ListTransactionalsResponse struct {
	Pagination ListTransactionalsPagination `json:"pagination"`
	Data       []TransactionalEmail         `json:"data"`
}

package loops

// ListPagination is the pagination shape used by paginated list endpoints.
type ListPagination struct {
	TotalResults    int     `json:"totalResults"`
	ReturnedResults int     `json:"returnedResults"`
	PerPage         int     `json:"perPage"`
	TotalPages      int     `json:"totalPages"`
	NextCursor      *string `json:"nextCursor"`
	NextPage        *string `json:"nextPage"`
}

// ThemeStyles represents style attributes returned for a theme.
type ThemeStyles struct {
	BackgroundColor       string  `json:"backgroundColor,omitempty"`
	BackgroundXPadding    float64 `json:"backgroundXPadding,omitempty"`
	BackgroundYPadding    float64 `json:"backgroundYPadding,omitempty"`
	BodyColor             string  `json:"bodyColor,omitempty"`
	BodyXPadding          float64 `json:"bodyXPadding,omitempty"`
	BodyYPadding          float64 `json:"bodyYPadding,omitempty"`
	BodyFontFamily        string  `json:"bodyFontFamily,omitempty"`
	BodyFontCategory      string  `json:"bodyFontCategory,omitempty"`
	BorderColor           string  `json:"borderColor,omitempty"`
	BorderWidth           float64 `json:"borderWidth,omitempty"`
	BorderRadius          float64 `json:"borderRadius,omitempty"`
	ButtonBodyColor       string  `json:"buttonBodyColor,omitempty"`
	ButtonBodyXPadding    float64 `json:"buttonBodyXPadding,omitempty"`
	ButtonBodyYPadding    float64 `json:"buttonBodyYPadding,omitempty"`
	ButtonBorderColor     string  `json:"buttonBorderColor,omitempty"`
	ButtonBorderWidth     float64 `json:"buttonBorderWidth,omitempty"`
	ButtonBorderRadius    float64 `json:"buttonBorderRadius,omitempty"`
	ButtonTextColor       string  `json:"buttonTextColor,omitempty"`
	ButtonTextFormat      float64 `json:"buttonTextFormat,omitempty"`
	ButtonTextFontSize    float64 `json:"buttonTextFontSize,omitempty"`
	DividerColor          string  `json:"dividerColor,omitempty"`
	DividerBorderWidth    float64 `json:"dividerBorderWidth,omitempty"`
	TextBaseColor         string  `json:"textBaseColor,omitempty"`
	TextBaseFontSize      float64 `json:"textBaseFontSize,omitempty"`
	TextBaseLineHeight    float64 `json:"textBaseLineHeight,omitempty"`
	TextBaseLetterSpacing float64 `json:"textBaseLetterSpacing,omitempty"`
	TextLinkColor         string  `json:"textLinkColor,omitempty"`
	Heading1Color         string  `json:"heading1Color,omitempty"`
	Heading1FontSize      float64 `json:"heading1FontSize,omitempty"`
	Heading1LineHeight    float64 `json:"heading1LineHeight,omitempty"`
	Heading1LetterSpacing float64 `json:"heading1LetterSpacing,omitempty"`
	Heading2Color         string  `json:"heading2Color,omitempty"`
	Heading2FontSize      float64 `json:"heading2FontSize,omitempty"`
	Heading2LineHeight    float64 `json:"heading2LineHeight,omitempty"`
	Heading2LetterSpacing float64 `json:"heading2LetterSpacing,omitempty"`
	Heading3Color         string  `json:"heading3Color,omitempty"`
	Heading3FontSize      float64 `json:"heading3FontSize,omitempty"`
	Heading3LineHeight    float64 `json:"heading3LineHeight,omitempty"`
	Heading3LetterSpacing float64 `json:"heading3LetterSpacing,omitempty"`
}

// Theme represents an email theme.
type Theme struct {
	ThemeID   string      `json:"themeId"`
	Name      string      `json:"name"`
	Styles    ThemeStyles `json:"styles"`
	IsDefault bool        `json:"isDefault"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

// ListThemesResponse is the 200 response for GET /themes.
type ListThemesResponse struct {
	Success    bool           `json:"success"`
	Pagination ListPagination `json:"pagination"`
	Data       []Theme        `json:"data"`
}

// ThemeResponse is the 200 response for GET /themes/{themeId}.
type ThemeResponse struct {
	Success   bool        `json:"success"`
	ThemeID   string      `json:"themeId"`
	Name      string      `json:"name"`
	Styles    ThemeStyles `json:"styles"`
	IsDefault bool        `json:"isDefault"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

// ThemeFailureResponse is used for theme request failures.
type ThemeFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Component represents an email component.
type Component struct {
	ComponentID string `json:"componentId"`
	Name        string `json:"name"`
	LMX         string `json:"lmx"`
}

// ListComponentsResponse is the 200 response for GET /components.
type ListComponentsResponse struct {
	Success    bool           `json:"success"`
	Pagination ListPagination `json:"pagination"`
	Data       []Component    `json:"data"`
}

// ComponentResponse is the 200 response for GET /components/{componentId}.
type ComponentResponse struct {
	Success     bool   `json:"success"`
	ComponentID string `json:"componentId"`
	Name        string `json:"name"`
	LMX         string `json:"lmx"`
}

// ComponentFailureResponse is used for component request failures.
type ComponentFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CampaignListItem is a single campaign returned from GET /campaigns.
type CampaignListItem struct {
	CampaignID     string  `json:"campaignId"`
	EmailMessageID *string `json:"emailMessageId"`
	Name           string  `json:"name"`
	Subject        string  `json:"subject"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

// ListCampaignsResponse is the 200 response for GET /campaigns.
type ListCampaignsResponse struct {
	Success    bool               `json:"success"`
	Pagination ListPagination     `json:"pagination"`
	Data       []CampaignListItem `json:"data"`
}

// CreateCampaignRequest is the body for POST /campaigns.
type CreateCampaignRequest struct {
	Name string `json:"name"`
}

// CreateCampaignResponse is the 201 response for POST /campaigns.
type CreateCampaignResponse struct {
	Success                       bool    `json:"success"`
	CampaignID                    string  `json:"campaignId"`
	Name                          string  `json:"name"`
	Status                        string  `json:"status"`
	CreatedAt                     string  `json:"createdAt"`
	UpdatedAt                     string  `json:"updatedAt"`
	EmailMessageID                string  `json:"emailMessageId"`
	EmailMessageContentRevisionID *string `json:"emailMessageContentRevisionId"`
}

// UpdateCampaignRequest is the body for POST /campaigns/{campaignId}.
type UpdateCampaignRequest struct {
	Name string `json:"name"`
}

// CampaignResponse is the 200 response for campaign reads and updates.
type CampaignResponse struct {
	Success        bool    `json:"success"`
	CampaignID     string  `json:"campaignId"`
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
	EmailMessageID *string `json:"emailMessageId"`
}

// CampaignFailureResponse is used for campaign request failures.
type CampaignFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateEmailMessageRequest is the body for POST /email-messages/{emailMessageId}.
type UpdateEmailMessageRequest struct {
	ExpectedRevisionID string `json:"expectedRevisionId,omitempty"`
	Subject            string `json:"subject,omitempty"`
	PreviewText        string `json:"previewText,omitempty"`
	FromName           string `json:"fromName,omitempty"`
	FromEmail          string `json:"fromEmail,omitempty"`
	ReplyToEmail       string `json:"replyToEmail,omitempty"`
	LMX                string `json:"lmx,omitempty"`
}

// EmailMessageResponse is the 200 response for email message reads and updates.
type EmailMessageResponse struct {
	Success           bool    `json:"success"`
	EmailMessageID    string  `json:"emailMessageId"`
	CampaignID        *string `json:"campaignId"`
	Subject           string  `json:"subject"`
	PreviewText       string  `json:"previewText"`
	FromName          string  `json:"fromName"`
	FromEmail         string  `json:"fromEmail"`
	ReplyToEmail      string  `json:"replyToEmail"`
	LMX               string  `json:"lmx"`
	ContentRevisionID *string `json:"contentRevisionId"`
	UpdatedAt         string  `json:"updatedAt"`
}

// EmailMessageFailureResponse is used for email message request failures.
type EmailMessageFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

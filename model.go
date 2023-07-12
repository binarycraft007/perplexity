package perplexity

type GetSidResponse struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int      `json:"pingTimeout"`
	PingInterval int      `json:"pingInterval"`
}

type AskRequest struct {
	Source                SearchSource `json:"source"`
	Version               string       `json:"version,omitempty"`
	Token                 string       `json:"token"`
	FrontendUUID          string       `json:"frontend_uuid"`
	FrontendSessionID     string       `json:"frontend_session_id,omitempty"`
	LastBackendUUID       string       `json:"last_backend_uuid,omitempty"`
	UseInhouseModel       bool         `json:"use_inhouse_model,omitempty"`
	ReadWriteToken        string       `json:"read_write_token,omitempty"`
	ConversationalEnabled bool         `json:"conversational_enabled"`
	AndroidDeviceID       string       `json:"android_device_id,omitempty"`
	Language              string       `json:"language,omitempty"`
	Timezone              string       `json:"timezone,omitempty"`
	SearchFocus           SearchFocus  `json:"search_focus,omitempty"`
	Gpt4                  bool         `json:"gpt4,omitempty"`
	Mode                  SearchMode   `json:"mode,omitempty"`
}

type AskResponse struct {
	Status             string `json:"status"`
	UUID               string `json:"uuid"`
	ReadWriteToken     string `json:"read_write_token"`
	QueryStr           string `json:"query_str"`
	StepType           string `json:"step_type"`
	RelatedQueries     []any  `json:"related_queries"`
	Text               string `json:"text"`
	Personalized       bool   `json:"personalized"`
	Mode               string `json:"mode"`
	Gpt4               bool   `json:"gpt4"`
	BackendUUID        string `json:"backend_uuid"`
	SearchFocus        string `json:"search_focus"`
	Label              string `json:"label"`
	ContextUUID        string `json:"context_uuid"`
	ThreadTitle        string `json:"thread_title"`
	AuthorUsername     any    `json:"author_username"`
	AuthorImage        any    `json:"author_image"`
	S3SocialPreviewURL string `json:"s3_social_preview_url"`
}

type AnswerDetails struct {
	Answer     string `json:"answer"`
	WebResults []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Snippet  string `json:"snippet"`
		Client   string `json:"client"`
		MetaData any    `json:"meta_data"`
	} `json:"web_results"`
	Chunks          []string `json:"chunks"`
	EntityLinks     any      `json:"entity_links"`
	ExtraWebResults []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Snippet  string `json:"snippet"`
		Client   string `json:"client"`
		MetaData any    `json:"meta_data"`
	} `json:"extra_web_results"`
	DeletedUrls   []any  `json:"deleted_urls"`
	SearchFocus   string `json:"search_focus"`
	ImageMetadata []any  `json:"image_metadata"`
}

type QueryProgress struct {
	Status         string `json:"status"`
	UUID           string `json:"uuid"`
	ReadWriteToken any    `json:"read_write_token"`
	Text           string `json:"text"`
	Final          bool   `json:"final"`
	BackendUUID    string `json:"backend_uuid"`
}

type CompleteResponse struct {
	Status         string `json:"status,omitempty"`
	UUID           string `json:"uuid,omitempty"`
	ReadWriteToken string `json:"read_write_token,omitempty"`
}

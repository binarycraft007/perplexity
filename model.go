package perplexity

type GetSidResponse struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int      `json:"pingTimeout"`
	PingInterval int      `json:"pingInterval"`
}

type AskRequest struct {
	Source                string `json:"source"`
	Version               string `json:"version"`
	Token                 string `json:"token"`
	FrontendUUID          string `json:"frontend_uuid"`
	LastBackendUUID       string `json:"last_backend_uuid,omitempty"`
	UseInhouseModel       bool   `json:"use_inhouse_model"`
	ReadWriteToken        string `json:"read_write_token,omitempty"`
	ConversationalEnabled bool   `json:"conversational_enabled"`
	AndroidDeviceID       string `json:"android_device_id"`
}

type AskResponse struct {
	Status             string   `json:"status"`
	UUID               string   `json:"uuid"`
	ReadWriteToken     string   `json:"read_write_token"`
	Mode               string   `json:"mode"`
	Label              string   `json:"label"`
	SearchFocus        string   `json:"search_focus"`
	StepType           string   `json:"step_type"`
	RelatedQueries     []string `json:"related_queries"`
	Gpt4               bool     `json:"gpt4"`
	BackendUUID        string   `json:"backend_uuid"`
	QueryStr           string   `json:"query_str"`
	Text               string   `json:"text"`
	ContextUUID        string   `json:"context_uuid"`
	ThreadTitle        string   `json:"thread_title"`
	AuthorUsername     any      `json:"author_username"`
	AuthorImage        any      `json:"author_image"`
	S3SocialPreviewURL string   `json:"s3_social_preview_url"`
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

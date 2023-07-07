package provider

type Webhook struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	OUCode     string   `json:"ouCode"`
	Enabled    bool     `json:"enabled"`
	EventTypes []string `json:"eventTypes"`
	URL        string   `json:"url"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
}

type WebhookResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	OUCode     string   `json:"ouCode"`
	Enabled    bool     `json:"enabled"`
	EventTypes []string `json:"eventTypes"`
	URL        string   `json:"url"`
	Headers    []Header `json:"headers"`
	Created    string   `json:"created"`
	Updated    string   `json:"updated"`
}

type WebhookConfig struct {
	Name       string   `json:"name"`
	OUCode     string   `json:"ouCode"`
	Enabled    bool     `json:"enabled"`
	EventTypes []string `json:"eventTypes"`
	Headers    []Header `json:"headers"`
	URL        string   `json:"url"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

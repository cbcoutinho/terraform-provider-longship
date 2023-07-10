package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (c *Client) GetWebhooks() ([]Webhook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/webhooks", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	webhooks := []Webhook{}
	err = json.Unmarshal(body, &webhooks)
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}

func (c *Client) GetWebhook(id string) (*WebhookResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/webhooks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	webhook := WebhookResponse{}
	err = json.Unmarshal(body, &webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

func (c *Client) CreateWebhook(webhook WebhookConfig) (*WebhookResponse, error) {
	rb, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/webhooks", c.HostURL), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newWebhook := WebhookResponse{}
	err = json.Unmarshal(body, &newWebhook)
	if err != nil {
		return nil, err
	}

	return &newWebhook, nil
}

func (c *Client) UpdateWebhook(id string, webhook WebhookConfig) (*WebhookResponse, error) {
	rb, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/webhooks/%s", c.HostURL, id), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newWebhook := WebhookResponse{}
	err = json.Unmarshal(body, &newWebhook)
	if err != nil {
		return nil, err
	}

	return &newWebhook, nil

}

func (c *Client) DeleteWebhook(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/webhooks/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

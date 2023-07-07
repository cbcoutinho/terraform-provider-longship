package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type webhooksDataSourceModel struct {
	Webhooks []WebhookDataSourceModel `tfsdk:"webhooks"`
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

func (c *Client) GetWebhook(id string) (*WebhookConfig, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/webhooks/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	webhook := WebhookConfig{}
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

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/webhooks", c.HostURL), strings.NewReader(string(rb)))

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

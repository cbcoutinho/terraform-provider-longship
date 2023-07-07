package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
}

type AuthStruct struct {
	TenantKey      string
	ApplicationKey string
}

func NewClient(host, tenantKey, applicationKey *string) (*Client, error) {

	// If username or password not provided, return empty client
	if host == nil || tenantKey == nil || applicationKey == nil {
		return nil, fmt.Errorf("Misconfigured client, missing host, tenant key, or Application key")
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    *host,
	}

	c.Auth = AuthStruct{
		TenantKey:      *tenantKey,
		ApplicationKey: *applicationKey,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

	req.Header.Set("Ocp-Apim-Subscription-Key", c.Auth.TenantKey)
	req.Header.Set("x-api-key", c.Auth.ApplicationKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

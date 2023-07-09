package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Chargepoint struct {
	ID            string `json:"id"`
	ChargepointID string `json:"chargePointId"`
	DateDeleted   string `json:"dateDeleted"`
	DisplayName   string `json:"displayName"`
	RoamingName   string `json:"roamingName"`
}

func (c *Client) GetChargepoints() ([]Chargepoint, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/chargepoints", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	chargepoints := []Chargepoint{}
	err = json.Unmarshal(body, &chargepoints)
	if err != nil {
		return nil, err
	}

	return chargepoints, nil
}

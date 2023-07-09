package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Chargepoint struct {
	ID                    string `json:"id"`
	ChargepointID         string `json:"chargePointId"`
	DateDeleted           string `json:"dateDeleted"`
	DisplayName           string `json:"displayName"`
	RoamingName           string `json:"roamingName"`
	ChargeBoxSerialNumber string `json:"chargeBoxSerialNumber"`
	ChargepointVendor     string `json:"chargePointVendor"`
	Evses                 []Evse `json:"evses"`
}

type Evse struct {
	EvseID     string      `json:"evse_id"`
	Connectors []Connector `json:"connectors"`
}

type Connector struct {
	ID                 string `json:"id"`
	OperationalStatus  string `json:"operationalStatus"`
	Standard           string `json:"standard"`
	Format             string `json:"format"`
	PowerType          string `json:"powerType"`
	MaxVoltage         int64  `json:"maxVoltage"`
	MaxAmperage        int64  `json:"maxAmperage"`
	MaxElectricalPower int64  `json:"maxElectricalPower"`
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

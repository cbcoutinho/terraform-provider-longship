package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OrganizationalUnit struct {
	ID               string           `json:"id"`
	ParentID         string           `json:"parentId"`
	Name             string           `json:"name"`
	FinancialDetails FinancialDetails `json:"financialDetails"`
}

type FinancialDetails struct {
	BeneficiaryName string `json:"beneficiaryName"`
	IBAN            string `json:"iban"`
	BIC             string `json:"bic"`
}

func (c *Client) GetOrganizationalUnits() ([]OrganizationalUnit, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/organizationalunits", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	organizationalUnits := []OrganizationalUnit{}
	err = json.Unmarshal(body, &organizationalUnits)
	if err != nil {
		return nil, err
	}

	return organizationalUnits, nil
}

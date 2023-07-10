package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OrganizationalUnit struct {
	ID                        string           `json:"id"`
	ParentID                  string           `json:"parentId"`
	Name                      string           `json:"name"`
	Code                      string           `json:"code"`
	ExternalReference         string           `json:"external_reference"`
	GridOwnerReference        string           `json:"grid_owner_reference"`
	TenantReference           string           `json:"tenant_reference"`
	CustomerReference         string           `json:"customer_reference"`
	Address                   string           `json:"address"`
	State                     string           `json:"state"`
	Country                   string           `json:"country"`
	City                      string           `json:"city"`
	HouseNumber               string           `json:"house_number"`
	PostalCode                string           `json:"postal_code"`
	HotlinePhoneNumber        string           `json:"hotline_phone_number"`
	CompanyEmail              string           `json:"company_email"`
	PrimaryContactPerson      string           `json:"primary_contact_person"`
	PrimaryContactPersonEmail string           `json:"primary_contact_person_email"`
	DirectPaymentProfileId    string           `json:"direct_payment_profile_id"`
	MspOuID                   string           `json:"msp_ou_id"`
	MspOuName                 string           `json:"msp_ou_name"`
	MspOuCode                 string           `json:"msp_ou_code"`
	MspExternalID             string           `json:"msp_external_id"`
	FinancialDetails          FinancialDetails `json:"financialDetails"`
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

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &OrganizationalUnitsDataSource{}
	_ datasource.DataSourceWithConfigure = &OrganizationalUnitsDataSource{}
)

// OrganizationalUnitsDataSource is the data source implementation.
type OrganizationalUnitsDataSource struct {
	client *Client
}

type OrganizationalUnitsDataSourceModel struct {
	OrganizationalUnits []OrganizationalUnitDataSourceModel `tfsdk:"organizational_units"`
}

type OrganizationalUnitDataSourceModel struct {
	ID                        types.String                    `tfsdk:"id"`
	ParentID                  types.String                    `tfsdk:"parent_id"`
	Name                      types.String                    `tfsdk:"name"`
	Code                      types.String                    `tfsdk:"code"`
	ExternalReference         types.String                    `tfsdk:"external_reference"`
	GridOwnerReference        types.String                    `tfsdk:"grid_owner_reference"`
	TenantReference           types.String                    `tfsdk:"tenant_reference"`
	CustomerReference         types.String                    `tfsdk:"customer_reference"`
	Address                   types.String                    `tfsdk:"address"`
	State                     types.String                    `tfsdk:"state"`
	Country                   types.String                    `tfsdk:"country"`
	City                      types.String                    `tfsdk:"city"`
	HouseNumber               types.String                    `tfsdk:"house_number"`
	PostalCode                types.String                    `tfsdk:"postal_code"`
	HotlinePhoneNumber        types.String                    `tfsdk:"hotline_phone_number"`
	CompanyEmail              types.String                    `tfsdk:"company_email"`
	PrimaryContactPerson      types.String                    `tfsdk:"primary_contact_person"`
	PrimaryContactPersonEmail types.String                    `tfsdk:"primary_contact_person_email"`
	DirectPaymentProfileId    types.String                    `tfsdk:"direct_payment_profile_id"`
	MspOuID                   types.String                    `tfsdk:"msp_ou_id"`
	MspOuName                 types.String                    `tfsdk:"msp_ou_name"`
	MspOuCode                 types.String                    `tfsdk:"msp_ou_code"`
	MspExternalID             types.String                    `tfsdk:"msp_external_id"`
	FinancialDetails          FinancialDetailsDataSourceModel `tfsdk:"financial_details"`
}

type FinancialDetailsDataSourceModel struct {
	BeneficiaryName types.String `tfsdk:"beneficiary_name"`
	IBAN            types.String `tfsdk:"iban"`
	BIC             types.String `tfsdk:"bic"`
}

func NewOrganizationalUnitsDataSource() datasource.DataSource {
	return &OrganizationalUnitsDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *OrganizationalUnitsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *OrganizationalUnitsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_units"
}

func (d *OrganizationalUnitsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organizational_units": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"parent_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"code": schema.StringAttribute{
							Computed: true,
						},
						"external_reference": schema.StringAttribute{
							Computed: true,
						},
						"grid_owner_reference": schema.StringAttribute{
							Computed: true,
						},
						"tenant_reference": schema.StringAttribute{
							Computed: true,
						},
						"customer_reference": schema.StringAttribute{
							Computed: true,
						},
						"address": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"country": schema.StringAttribute{
							Computed: true,
						},
						"city": schema.StringAttribute{
							Computed: true,
						},
						"house_number": schema.StringAttribute{
							Computed: true,
						},
						"postal_code": schema.StringAttribute{
							Computed: true,
						},
						"hotline_phone_number": schema.StringAttribute{
							Computed: true,
						},
						"company_email": schema.StringAttribute{
							Computed: true,
						},
						"primary_contact_person": schema.StringAttribute{
							Computed: true,
						},
						"primary_contact_person_email": schema.StringAttribute{
							Computed: true,
						},
						"direct_payment_profile_id": schema.StringAttribute{
							Computed: true,
						},
						"msp_ou_id": schema.StringAttribute{
							Computed: true,
						},
						"msp_ou_name": schema.StringAttribute{
							Computed: true,
						},
						"msp_ou_code": schema.StringAttribute{
							Computed: true,
						},
						"msp_external_id": schema.StringAttribute{
							Computed: true,
						},
						"financial_details": schema.ObjectAttribute{
							Computed: true,
							AttributeTypes: map[string]attr.Type{
								"beneficiary_name": types.StringType,
								"iban":             types.StringType,
								"bic":              types.StringType,
							},
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationalUnitsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state OrganizationalUnitsDataSourceModel

	organizationalUnits, err := d.client.GetOrganizationalUnits()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Longship Organizational Units",
			err.Error(),
		)
		return
	}

	tflog.Warn(ctx, fmt.Sprintf("Reading all organizational units: %+v", organizationalUnits))

	// Map response body to model
	for _, ou := range organizationalUnits {

		tflog.Warn(ctx, fmt.Sprintf("Organizational Unit: %+v", ou))

		ouState := OrganizationalUnitDataSourceModel{
			ID:                        types.StringValue(ou.ID),
			ParentID:                  types.StringValue(ou.ParentID),
			Name:                      types.StringValue(ou.Name),
			Code:                      types.StringValue(ou.Code),
			ExternalReference:         types.StringValue(ou.ExternalReference),
			GridOwnerReference:        types.StringValue(ou.GridOwnerReference),
			TenantReference:           types.StringValue(ou.TenantReference),
			CustomerReference:         types.StringValue(ou.CustomerReference),
			Address:                   types.StringValue(ou.Address),
			State:                     types.StringValue(ou.State),
			Country:                   types.StringValue(ou.Country),
			City:                      types.StringValue(ou.City),
			HouseNumber:               types.StringValue(ou.HouseNumber),
			PostalCode:                types.StringValue(ou.PostalCode),
			HotlinePhoneNumber:        types.StringValue(ou.HotlinePhoneNumber),
			CompanyEmail:              types.StringValue(ou.CompanyEmail),
			PrimaryContactPerson:      types.StringValue(ou.PrimaryContactPerson),
			PrimaryContactPersonEmail: types.StringValue(ou.PrimaryContactPersonEmail),
			DirectPaymentProfileId:    types.StringValue(ou.DirectPaymentProfileId),
			MspOuID:                   types.StringValue(ou.MspOuID),
			MspOuName:                 types.StringValue(ou.MspOuName),
			MspOuCode:                 types.StringValue(ou.MspOuCode),
			MspExternalID:             types.StringValue(ou.MspExternalID),

			FinancialDetails: FinancialDetailsDataSourceModel{
				BeneficiaryName: types.StringValue(ou.FinancialDetails.BeneficiaryName),
				IBAN:            types.StringValue(ou.FinancialDetails.IBAN),
				BIC:             types.StringValue(ou.FinancialDetails.BIC),
			},
		}

		state.OrganizationalUnits = append(state.OrganizationalUnits, ouState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

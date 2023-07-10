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
	ID               types.String                    `tfsdk:"id"`
	ParentID         types.String                    `tfsdk:"parent_id"`
	Name             types.String                    `tfsdk:"name"`
	FinancialDetails FinancialDetailsDataSourceModel `tfsdk:"financial_details"`
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
			ID:       types.StringValue(ou.ID),
			ParentID: types.StringValue(ou.ParentID),
			Name:     types.StringValue(ou.Name),
		}

		ouState.FinancialDetails = FinancialDetailsDataSourceModel{
			BeneficiaryName: types.StringValue(ou.FinancialDetails.BeneficiaryName),
			IBAN:            types.StringValue(ou.FinancialDetails.IBAN),
			BIC:             types.StringValue(ou.FinancialDetails.BIC),
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

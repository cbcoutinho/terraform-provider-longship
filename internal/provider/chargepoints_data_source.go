package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ChargepointsDataSource{}
	_ datasource.DataSourceWithConfigure = &ChargepointsDataSource{}
)

// ChargepointsDataSource is the data source implementation.
type ChargepointsDataSource struct {
	client *Client
}

type ChargepointsDataSourceModel struct {
	Chargepoints []ChargepointDataSourceModel `tfsdk:"chargepoints"`
}

type ChargepointDataSourceModel struct {
	ID                    types.String          `tfsdk:"id"`
	ChargepointID         types.String          `tfsdk:"chargepoint_id"`
	DateDeleted           types.String          `tfsdk:"date_deleted"`
	DisplayName           types.String          `tfsdk:"display_name"`
	RoamingName           types.String          `tfsdk:"roaming_name"`
	ChargeBoxSerialNumber types.String          `tfsdk:"charge_box_serial_number"`
	ChargepointVendor     types.String          `tfsdk:"chargepoint_vendor"`
	Evses                 []EvseDataSourceModel `tfsdk:"evses"`
}

type EvseDataSourceModel struct {
	EvseID     types.String               `tfsdk:"evse_id"`
	Connectors []ConnectorDataSourceModel `tfsdk:"connectors"`
}

type ConnectorDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	OperationalStatus  types.String `tfsdk:"operational_status"`
	Standard           types.String `tfsdk:"standard"`
	Format             types.String `tfsdk:"format"`
	PowerType          types.String `tfsdk:"power_type"`
	MaxVoltage         types.Int64  `tfsdk:"max_voltage"`
	MaxAmperage        types.Int64  `tfsdk:"max_amperage"`
	MaxElectricalPower types.Int64  `tfsdk:"max_electrical_power"`
}

// NewChargepointsDataSource is a helper function to simplify the provider implementation.
func NewChargepointsDataSource() datasource.DataSource {
	return &ChargepointsDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *ChargepointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ChargepointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chargepoints"
}

func (d *ChargepointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"chargepoints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"chargepoint_id": schema.StringAttribute{
							Computed: true,
						},
						"date_deleted": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"roaming_name": schema.StringAttribute{
							Computed: true,
						},
						"charge_box_serial_number": schema.StringAttribute{
							Computed: true,
						},
						"chargepoint_vendor": schema.StringAttribute{
							Computed: true,
						},
						"evses": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"evse_id": schema.StringAttribute{
										Computed: true,
									},
									"connectors": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Computed: true,
												},
												"operational_status": schema.StringAttribute{
													Computed: true,
												},
												"standard": schema.StringAttribute{
													Computed: true,
												},
												"format": schema.StringAttribute{
													Computed: true,
												},
												"power_type": schema.StringAttribute{
													Computed: true,
												},
												"max_voltage": schema.Int64Attribute{
													Computed: true,
												},
												"max_amperage": schema.Int64Attribute{
													Computed: true,
												},
												"max_electrical_power": schema.Int64Attribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ChargepointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ChargepointsDataSourceModel

	chargepoints, err := d.client.GetChargepoints()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Longship Webhooks",
			err.Error(),
		)
		return
	}

	for _, chargepoint := range chargepoints {
		chargepointState := ChargepointDataSourceModel{
			ID:                    types.StringValue(chargepoint.ID),
			ChargepointID:         types.StringValue(chargepoint.ChargepointID),
			DateDeleted:           types.StringValue(chargepoint.DateDeleted),
			DisplayName:           types.StringValue(chargepoint.DisplayName),
			RoamingName:           types.StringValue(chargepoint.RoamingName),
			ChargeBoxSerialNumber: types.StringValue(chargepoint.ChargeBoxSerialNumber),
			ChargepointVendor:     types.StringValue(chargepoint.ChargepointVendor),
		}

		evsesState := []EvseDataSourceModel{}
		for _, evse := range chargepoint.Evses {
			connectorsState := []ConnectorDataSourceModel{}
			for _, connector := range evse.Connectors {
				connectorState := ConnectorDataSourceModel{
					ID:                 types.StringValue(connector.ID),
					OperationalStatus:  types.StringValue(connector.OperationalStatus),
					Standard:           types.StringValue(connector.Standard),
					Format:             types.StringValue(connector.Format),
					PowerType:          types.StringValue(connector.PowerType),
					MaxVoltage:         types.Int64Value(connector.MaxVoltage),
					MaxAmperage:        types.Int64Value(connector.MaxAmperage),
					MaxElectricalPower: types.Int64Value(connector.MaxElectricalPower),
				}
				connectorsState = append(connectorsState, connectorState)
			}
			evsesState = append(evsesState, EvseDataSourceModel{
				EvseID:     types.StringValue(evse.EvseID),
				Connectors: connectorsState,
			})
		}
		chargepointState.Evses = evsesState

		state.Chargepoints = append(state.Chargepoints, chargepointState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

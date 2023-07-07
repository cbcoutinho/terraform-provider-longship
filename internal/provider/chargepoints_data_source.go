package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &chargepointsDataSource{}
)

// NewChargepointsDataSource is a helper function to simplify the provider implementation.
func NewChargepointsDataSource() datasource.DataSource {
	return &chargepointsDataSource{}
}

// chargepointsDataSource is the data source implementation.
type chargepointsDataSource struct {
	client *http.Client
}

// Metadata returns the data source type name.
func (d *chargepointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chargepoints"
}

// Read refreshes the Terraform state with the latest data.
func (d *chargepointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

func (d *chargepointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"chargepoints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"id": schema.StringAttribute{
							Computed: true,
						},
						"charge_point_id": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"reimburse_tariff_price": schema.StringAttribute{
							Computed: true,
						},
						"reimburse_tariff_id": schema.StringAttribute{
							Computed: true,
						},
						"has_gues_usage": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &WebhooksDataSource{}
	_ datasource.DataSourceWithConfigure = &WebhooksDataSource{}
)

type WebhooksDataSource struct {
	client *Client
}

type webhooksDataSourceModel struct {
	Webhooks []WebhookDataSourceModel `tfsdk:"webhooks"`
}

type WebhookDataSourceModel struct {
	ID         types.String   `tfsdk:"id"`
	Name       types.String   `tfsdk:"name"`
	Enabled    types.Bool     `tfsdk:"enabled"`
	EventTypes []types.String `tfsdk:"event_types"`
	Created    types.String   `tfsdk:"created"`
	Updated    types.String   `tfsdk:"updated"`
}

func NewWebhooksDataSource() datasource.DataSource {
	return &WebhooksDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *WebhooksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *WebhooksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks"
}

// Schema defines the schema for the data source.
func (d *WebhooksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of webhooks",
		Attributes: map[string]schema.Attribute{
			"webhooks": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the webhook.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the webhook.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Webhook enabled or not.",
							Computed:    true,
						},
						"event_types": schema.ListAttribute{
							Description: "Notifications triggered with this webhook.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"created": schema.StringAttribute{
							Description: "Timestamp of when webhook was created.",
							Computed:    true,
						},
						"updated": schema.StringAttribute{
							Description: "Timestamp of when webhook was last updated.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *WebhooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state webhooksDataSourceModel

	webhooks, err := d.client.GetWebhooks()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Longship Webhooks",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, webhook := range webhooks {
		webhookState := WebhookDataSourceModel{
			ID:      types.StringValue(webhook.ID),
			Name:    types.StringValue(webhook.Name),
			Enabled: types.BoolValue(webhook.Enabled),
			Created: types.StringValue(webhook.Created),
			Updated: types.StringValue(webhook.Updated),
		}

		for _, eventType := range webhook.EventTypes {
			webhookState.EventTypes = append(webhookState.EventTypes, types.StringValue(eventType))
		}

		state.Webhooks = append(state.Webhooks, webhookState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

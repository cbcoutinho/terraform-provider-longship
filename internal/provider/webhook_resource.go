package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &webhookResource{}
	_ resource.ResourceWithConfigure   = &webhookResource{}
	_ resource.ResourceWithImportState = &webhookResource{}
)

type WebhookResourceModel struct {
	ID         types.String   `tfsdk:"id"`
	Name       types.String   `tfsdk:"name"`
	OUCode     types.String   `tfsdk:"ou_code"`
	Enabled    types.Bool     `tfsdk:"enabled"`
	EventTypes []types.String `tfsdk:"event_types"`
	URL        types.String   `tfsdk:"url"`
	Headers    []HeaderModel  `tfsdk:"headers"`
	Created    types.String   `tfsdk:"created"`
	Updated    types.String   `tfsdk:"updated"`
}

type HeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// Configure adds the provider configured client to the data source.
func (d *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	tflog.Debug(ctx, "Retrieving Longship API client")

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

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Schema defines the schema for the resource.
func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"ou_code": schema.StringAttribute{
				Required: true,
			},
			"enabled": schema.BoolAttribute{
				Required: true,
			},
			"event_types": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"headers": schema.ListNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan WebhookResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Planning webhookResource: %s", plan))

	var eventTypes []string
	for _, eventType := range plan.EventTypes {
		eventTypes = append(eventTypes, eventType.ValueString())
	}

	headers := []Header{}
	for _, header := range plan.Headers {
		headers = append(headers, Header{
			Name:  header.Name.ValueString(),
			Value: header.Value.ValueString(),
		})
	}

	config := WebhookConfig{
		Name:       plan.Name.ValueString(),
		OUCode:     plan.OUCode.ValueString(),
		Enabled:    plan.Enabled.ValueBool(),
		EventTypes: eventTypes,
		Headers:    headers,
		URL:        plan.URL.ValueString(),
	}

	tflog.Info(ctx, fmt.Sprintf("Creating webhook: %+v", config))

	webhook, err := r.client.CreateWebhook(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Create webhook response: %+v", webhook))

	plan.ID = types.StringValue(webhook.ID)
	plan.Name = types.StringValue(webhook.Name)
	plan.OUCode = types.StringValue(webhook.OUCode)
	plan.Enabled = types.BoolValue(webhook.Enabled)
	plan.URL = types.StringValue(webhook.URL)
	for idx, eventType := range webhook.EventTypes {
		plan.EventTypes[idx] = types.StringValue(eventType)
	}
	for idx, header := range webhook.Headers {
		plan.Headers[idx] = HeaderModel{
			Name:  types.StringValue(header.Name),
			Value: types.StringValue(header.Value),
		}
	}

	plan.Updated = types.StringValue(webhook.Updated)
	plan.Created = types.StringValue(webhook.Created)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state WebhookResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading all webhooks")

	exists := false
	webhooks, err := r.client.GetWebhooks()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading webhooks",
			"Could not read webhook, unexpected error: "+err.Error(),
		)
		return
	}

	for _, w := range webhooks {
		tflog.Info(ctx, fmt.Sprintf("Found webhook[%s]: %+v", w.ID, w))
		if state.ID.ValueString() == w.ID {
			exists = true
		}
	}

	if !exists {
		tflog.Info(ctx, "Webhook does not exist!")
		state = WebhookResourceModel{}

		// Set refreshed state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Reading webhook id: %s", state.ID.ValueString()))

	// Get refreshed webhook value from Longship
	webhook, err := r.client.GetWebhook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Longship Webhook",
			"Could not read Longship webhook ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("GET /v1/webhooks/%s response: %+v", state.ID.ValueString(), webhook))

	// Overwrite attributes with refreshed state
	state.Name = types.StringValue(webhook.Name)
	state.OUCode = types.StringValue(webhook.OUCode)
	state.Enabled = types.BoolValue(webhook.Enabled)
	state.URL = types.StringValue(webhook.URL)
	state.Updated = types.StringValue(webhook.Updated)
	state.Created = types.StringValue(webhook.Created)

	state.EventTypes = []types.String{}
	for _, eventType := range webhook.EventTypes {
		state.EventTypes = append(state.EventTypes, types.StringValue(eventType))
	}
	state.Headers = []HeaderModel{}
	for _, header := range webhook.Headers {
		state.Headers = append(state.Headers, HeaderModel{
			Name:  types.StringValue(header.Name),
			Value: types.StringValue(header.Value),
		})
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan WebhookResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Updating webhook id: %s", plan.ID.ValueString()))

	var eventTypes []string
	for _, eventType := range plan.EventTypes {
		eventTypes = append(eventTypes, eventType.ValueString())
	}

	var headers []Header
	for _, header := range plan.Headers {
		headers = append(headers, Header{
			Name:  header.Name.ValueString(),
			Value: header.Value.ValueString(),
		})
	}

	config := WebhookConfig{
		Name:       plan.Name.ValueString(),
		OUCode:     plan.OUCode.ValueString(),
		Enabled:    plan.Enabled.ValueBool(),
		EventTypes: eventTypes,
		Headers:    headers,
		URL:        plan.URL.ValueString(),
	}

	webhook, err := r.client.UpdateWebhook(plan.ID.ValueString(), config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(webhook.ID)
	plan.Name = types.StringValue(webhook.Name)
	plan.OUCode = types.StringValue(webhook.OUCode)
	plan.Enabled = types.BoolValue(webhook.Enabled)
	plan.URL = types.StringValue(webhook.URL)
	for idx, eventType := range webhook.EventTypes {
		plan.EventTypes[idx] = types.StringValue(eventType)
	}
	for idx, header := range webhook.Headers {
		plan.Headers[idx] = HeaderModel{
			Name:  types.StringValue(header.Name),
			Value: types.StringValue(header.Value),
		}
	}

	plan.Updated = types.StringValue(webhook.Updated)
	plan.Created = types.StringValue(webhook.Created)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state WebhookResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting webhook id: %s", state.ID.ValueString()))

	// Get refreshed webhook value from Longship
	err := r.client.DeleteWebhook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Longship Webhook",
			"Could not read Longship webhook ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	tflog.Info(ctx, fmt.Sprintf("Importing webhook id: %s", path.Root("id")))

	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

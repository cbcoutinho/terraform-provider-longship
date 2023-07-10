package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &longshipProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &longshipProvider{
			version: version,
		}
	}
}

// longshipProvider is the provider implementation.
type longshipProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *longshipProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "longship"
	resp.Version = p.version
}

type longshipProviderModel struct {
	Host           types.String `tfsdk:"host"`
	TenantKey      types.String `tfsdk:"tenant_key"`
	ApplicationKey types.String `tfsdk:"application_key"`
}

// Schema defines the provider-level schema for configuration data.
func (p *longshipProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Longship API. May also be provided via LONGSHIP_HOST environment variable.",
				Optional:    true,
			},
			"tenant_key": schema.StringAttribute{
				Description: "Tenant key for Longship API. May also be provided via LONGSHIP_TENANT_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"application_key": schema.StringAttribute{
				Description: "Application key for Longship API. May also be provided via LONGSHIP_APPLICATION_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Longship API client for data sources and resources.
func (p *longshipProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Debug(ctx, "Configuring Longship client")

	// Retrieve provider data from configuration
	var config longshipProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Longship API Host",
			"The provider cannot create the Longship API client as there is an unknown configuration value for the Longship API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the `LONGSHIP_HOST` environment variable.",
		)
	}

	if config.TenantKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant_key"),
			"Unknown Longship API Tenant Key",
			"The provider cannot create the Longship API client as there is an unknown configuration value for the Longship API tenant key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the `LONGSHIP_TENANT_KEY` environment variable.",
		)
	}

	if config.ApplicationKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_key"),
			"Unknown Longship API Application Key",
			"The provider cannot create the Longship API client as there is an unknown configuration value for the Longship API application key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the `LONGSHIP_APPLICATION_KEY` environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("LONGSHIP_HOST")
	tenantKey := os.Getenv("LONGSHIP_TENANT_KEY")
	applicationKey := os.Getenv("LONGSHIP_APPLICATION_KEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.TenantKey.IsNull() {
		tenantKey = config.TenantKey.ValueString()
	}

	if !config.ApplicationKey.IsNull() {
		applicationKey = config.ApplicationKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Longship API Host",
			"The provider cannot create the Longship API client as there is a missing or empty value for the Longship API host. "+
				"Set the host value in the configuration or use the LONGSHIP_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if tenantKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenantKey"),
			"Missing Longship API Tenant Key",
			"The provider cannot create the Longship API client as there is a missing or empty value for the Longship API Tenant Key. "+
				"Set the tenant_key value in the configuration or use the LONGSHIP_TENANT_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if applicationKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("applicationKey"),
			"Missing Longship API Application Key",
			"The provider cannot create the Longship API client as there is a missing or empty value for the Longship API Application Key. "+
				"Set the application_key value in the configuration or use the LONGSHIP_APPLICATION_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "longship_host", host)
	ctx = tflog.SetField(ctx, "longship_tenant_key", tenantKey)
	ctx = tflog.SetField(ctx, "longship_application_key", applicationKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "longship_tenant_key")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "longship_application_key")

	tflog.Debug(ctx, "Creating Longship client")

	// Create a new Longship client using the configuration values
	client, err := NewClient(&host, &tenantKey, &applicationKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Longship API Client",
			"An unexpected error occurred when creating the Longship API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Longship Client Error: "+err.Error(),
		)
		return
	}

	// Make the Longship client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Longship client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *longshipProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewChargepointsDataSource,
		NewWebhooksDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *longshipProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebhookResource,
	}
}

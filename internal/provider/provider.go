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
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

// Ensure KeaCBProvider satisfies various provider interfaces.
var _ provider.Provider = &KeaCBProvider{}

// KeaCBProvider defines the provider implementation for Kea configuration-backend.
type KeaCBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KeaProviderModel describes the provider data model.
type KeaProviderModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata : Defines the provider metadata.
func (p *KeaCBProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kea"
	resp.Version = p.version
}

// Schema : Defines the provider schema.
func (p *KeaCBProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "Kea ctrl-agent username. Defaults to env var `KEA_USERNAME` if not specified.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Kea ctrl-agent password. Defaults to env var `KEA_PASSWORD` if not specified.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

// Configure : Defines the provider configuration.
func (p *KeaCBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Define an empty configuration.
	var config KeaProviderModel

	// Read/populate the provider data from the configuration.
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// If there are any diagnostics, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Catch any unknown attribute errors here and stop.
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the username and password from the environment, but override if
	// the practitioner provided a value in the configuration.
	username := os.Getenv("KEA_USERNAME")
	password := os.Getenv("KEA_PASSWORD")

	if username == "" {
		username = config.Username.ValueString()
	}
	if password == "" {
		password = config.Password.ValueString()
	}

	// After all is set, check if the username and password are empty.
	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Kea DHCP API Username",
			"The provider cannot create the Kea DHCP API client as there is a missing or empty value for "+
				"the Kea DHCP API username. Set the username value in the configuration or use the KEA_USERNAME "+
				"environment variable. If either is already set, ensure the value is not empty.",
		)
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Kea DHCP API Password",
			"The provider cannot create the Kea DHCP API client as there is a missing or empty value for "+
				"the Kea DHCP API password. Set the password value in the configuration or use the KEA_PASSWORD "+
				"environment variable. If either is already set, ensure the value is not empty.",
		)
	}
	// Stop here if there are any errors.
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Kea DHCP API client.
	client := kea.New(kea.WithAuth(username, password))

	// Make the Kea DHCP client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources : Defines the provider resources.
func (p *KeaCBProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRemoteSubnet4Resource,
		NewRemoteOptionDef4Resource,
		NewReservationResource,
	}
}

// DataSources : Defines the provider data sources.
func (p *KeaCBProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRemoteSubnet4DataSource,
		NewRemoteOptionDef4DataSource,
		NewReservationDataSource,
	}
}

// New : Creates a new provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KeaCBProvider{
			version: version,
		}
	}
}

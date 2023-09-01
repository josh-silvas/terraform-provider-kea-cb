// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

var (
	// Ensure provider defined types fully satisfy framework interfaces.
	_ datasource.DataSource              = &remoteOptionDef4DataSource{}
	_ datasource.DataSourceWithConfigure = &remoteOptionDef4DataSource{}
)

// NewRemoteOptionDef4DataSource : Creates a new empty data source client.
func NewRemoteOptionDef4DataSource() datasource.DataSource {
	return &remoteOptionDef4DataSource{}
}

type (
	// remoteOptionDef4DataSource defines the data source client.
	remoteOptionDef4DataSource struct {
		client *kea.Client
	}

	// remoteOptionDef4DataSourceSchema describes the data source data model.
	// Maps to the source schema data.
	remoteOptionDef4DataSourceSchema struct {
		Hostname    types.String `tfsdk:"hostname"`
		Name        types.String `tfsdk:"name"`
		Code        types.Int64  `tfsdk:"code"`
		Type        types.String `tfsdk:"type"`
		Array       types.Bool   `tfsdk:"array"`
		RecordTypes types.String `tfsdk:"record_types"`
		Space       types.String `tfsdk:"space"`
		Encapsulate types.String `tfsdk:"encapsulate"`
	}
)

// Metadata : Defines the data source metadata.
func (d *remoteOptionDef4DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_option_def4_data_source"
}

// Schema : Defines the data source schema.
func (d *remoteOptionDef4DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reservation data source",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the kea server to connect to. e.g. `kea.example.com`",
				Required:            true,
			},
			"code": schema.Int64Attribute{
				MarkdownDescription: "DHCP option code. e.g. `222`",
				Required:            true,
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The DHCP space for the option-def. e.g. `dhcp4`.",
				Required:            true,
			},
			"name":         schema.StringAttribute{Computed: true},
			"type":         schema.StringAttribute{Computed: true},
			"array":        schema.BoolAttribute{Computed: true},
			"record_types": schema.StringAttribute{Computed: true},
			"encapsulate":  schema.StringAttribute{Computed: true},
		},
	}
}

// Configure : Configures the data source client.
func (d *remoteOptionDef4DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Fetch the Kea DHCP client from the provider.
	client, ok := req.ProviderData.(*kea.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kea.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

// Read : Reads the data source data into the Terraform state.
func (d *remoteOptionDef4DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Define an empty configuration.
	var config remoteOptionDef4DataSourceSchema

	// Read Terraform configuration data into the model
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Get", "`hostname` field is required")
	}

	//  If the Code value is empty, add an error to the diagnostics.
	if config.Code.IsNull() || config.Code.IsUnknown() {
		resp.Diagnostics.AddError("RemoteOptionDef4Get", "`code` field is required")
	}

	//  If the Space value is empty, add an error to the diagnostics.
	if config.Space.IsNull() || config.Space.IsUnknown() || config.Space.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Get", "`space` field is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	respData, err := d.client.RemoteOptionDef4Get(
		config.Hostname.ValueString(),
		config.Space.ValueString(),
		int(config.Code.ValueInt64()),
	)
	if err != nil {
		// Only return an error if the error is NOT subnet not found.
		if !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"RemoteOptionDef4Get",
				fmt.Sprintf("Unable to read remote-option-def4, got error: %s", err),
			)
			return
		}
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	if respData == nil {
		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}

	// Marshalling the response data taken from Kea, and write
	// it into the TF  model.
	config.Type = types.StringValue(respData.Type)
	config.Array = types.BoolValue(respData.Array)
	config.RecordTypes = types.StringValue(respData.RecordTypes)
	config.Encapsulate = types.StringValue(respData.Encapsulate)

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

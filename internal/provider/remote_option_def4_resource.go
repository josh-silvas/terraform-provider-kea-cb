package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

var (
	// Ensure provider defined types fully satisfy framework interfaces.
	_ resource.Resource                = &remoteOptionDef4Resource{}
	_ resource.ResourceWithImportState = &remoteOptionDef4Resource{}
)

// NewRemoteOptionDef4Resource : Creates a new empty resource client.
func NewRemoteOptionDef4Resource() resource.Resource {
	return &remoteOptionDef4Resource{}
}

type (
	// remoteOptionDef4Resource defines the resource implementation.
	remoteOptionDef4Resource struct {
		client *kea.Client
	}

	// remoteOptionDef4ResourceSchema describes the resource data model.
	remoteOptionDef4ResourceSchema struct {
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

// Metadata : Returns the resource type name and supported features.
func (r *remoteOptionDef4Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_option_def4_resource"
}

// Schema : Returns the resource schema.
func (r *remoteOptionDef4Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Remote OptionDef4 resource",

		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the kea server to connect to. e.g. `kea.example.com`",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "DHCP option name. e.g. `location-identifier`",
				Required:            true,
			},
			"code": schema.Int64Attribute{
				MarkdownDescription: "DHCP option code. e.g. `222`",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "DHCP option type. e.g. `string`, `uint32`",
				Required:            true,
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The DHCP space for the option-def. e.g. `dhcp4`.",
				Required:            true,
			},
			"array": schema.BoolAttribute{
				MarkdownDescription: "The false value of the array parameter determines that the option does NOT comprise an array of uint32 values but is, instead, a single value..",
				Optional:            true,
			},
			"record_types": schema.StringAttribute{
				MarkdownDescription: "The record_types value should be non-empty if type is set to \"record\"; otherwise it must be left blank. ",
				Optional:            true,
			},
			"encapsulate": schema.StringAttribute{
				MarkdownDescription: "The name of the option space in which the sub-options are defined.",
				Optional:            true,
			},
		},
	}
}

// Configure : Configures the resource client data and populates the client interface.
func (r *remoteOptionDef4Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = client
}

// Create : Creates a new resource.
func (r *remoteOptionDef4Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config remoteOptionDef4ResourceSchema

	// Read Terraform configuration data into the model
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`hostname` field is required")
	}

	//  If the Code value is empty, add an error to the diagnostics.
	if config.Code.IsNull() || config.Code.IsUnknown() {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`code` field is required")
	}

	//  If the Name value is empty, add an error to the diagnostics.
	if config.Name.IsNull() || config.Name.IsUnknown() || config.Name.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`name` field is required")
	}

	//  If the Type value is empty, add an error to the diagnostics.
	if config.Type.IsNull() || config.Type.IsUnknown() || config.Type.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`type` field is required")
	}

	//  If the Space value is empty, add an error to the diagnostics.
	if config.Space.IsNull() || config.Space.IsUnknown() || config.Space.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`space` field is required")
	}

	// If there are any diagnostics, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	def := kea.RemoteOptionDef4{
		Name:  config.Name.ValueString(),
		Code:  int(config.Code.ValueInt64()),
		Type:  config.Type.ValueString(),
		Space: config.Space.ValueString(),
	}

	if !config.Array.IsNull() && !config.Array.IsUnknown() {
		def.Array = config.Array.ValueBool()
	}
	if !config.RecordTypes.IsNull() && !config.RecordTypes.IsUnknown() {
		def.RecordTypes = config.RecordTypes.ValueString()
	}
	if !config.Encapsulate.IsNull() && !config.Encapsulate.IsUnknown() {
		def.Encapsulate = config.Encapsulate.ValueString()
	}

	// nolint: contextcheck
	if err := r.client.RemoteOptionDef4Set(config.Hostname.ValueString(), def); err != nil {
		resp.Diagnostics.AddError(
			"RemoteOptionDef4Set",
			fmt.Sprintf("Unable to create option-def4 in Kea, got error: %s | %v", err, def),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Read : Reads the resource data into the Terraform state.
func (r *remoteOptionDef4Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config remoteOptionDef4ResourceSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`hostname` field is required")
	}

	//  If the Code value is empty, add an error to the diagnostics.
	if config.Code.IsNull() || config.Code.IsUnknown() {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`code` field is required")
	}

	//  If the Space value is empty, add an error to the diagnostics.
	if config.Space.IsNull() || config.Space.IsUnknown() || config.Space.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Set", "`space` field is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	respData, err := r.client.RemoteOptionDef4Get(
		config.Hostname.ValueString(),
		config.Space.ValueString(),
		int(config.Code.ValueInt64()),
	)
	if err != nil {
		// Only return an error if the error is NOT subnet not found.
		if !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"RemoteOptionDef4Set",
				fmt.Sprintf("Unable to read remote-option-def4, got error: %s", err),
			)
			return
		}
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Marshalling the response data taken from Kea, and write
	// it into the TF  model.
	if respData.Type != "" {
		config.Type = types.StringValue(respData.Type)
	}
	if respData.Array {
		config.Array = types.BoolValue(respData.Array)
	}
	if respData.RecordTypes != "" {
		config.RecordTypes = types.StringValue(respData.RecordTypes)
	}
	if respData.Encapsulate != "" {
		config.Encapsulate = types.StringValue(respData.Encapsulate)
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Update : Updates an existing resource.
func (r *remoteOptionDef4Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config remoteOptionDef4ResourceSchema

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Update", "`hostname` field is required")
	}

	//  If the Code value is empty, add an error to the diagnostics.
	if config.Code.IsNull() || config.Code.IsUnknown() {
		resp.Diagnostics.AddError("RemoteOptionDef4Update", "`code` field is required")
	}

	//  If the Name value is empty, add an error to the diagnostics.
	if config.Name.IsNull() || config.Name.IsUnknown() || config.Name.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Update", "`name` field is required")
	}

	//  If the Type value is empty, add an error to the diagnostics.
	if config.Type.IsNull() || config.Type.IsUnknown() || config.Type.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Update", "`type` field is required")
	}

	//  If the Space value is empty, add an error to the diagnostics.
	if config.Space.IsNull() || config.Space.IsUnknown() || config.Space.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Update", "`space` field is required")
	}

	// If there are any diagnostics, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	def := kea.RemoteOptionDef4{
		Name:  config.Name.ValueString(),
		Code:  int(config.Code.ValueInt64()),
		Type:  config.Type.ValueString(),
		Space: config.Space.ValueString(),
	}

	if !config.Array.IsNull() && !config.Array.IsUnknown() {
		def.Array = config.Array.ValueBool()
	}
	if !config.RecordTypes.IsNull() && !config.RecordTypes.IsUnknown() {
		def.RecordTypes = config.RecordTypes.ValueString()
	}
	if !config.Encapsulate.IsNull() && !config.Encapsulate.IsUnknown() {
		def.Encapsulate = config.Encapsulate.ValueString()
	}

	// nolint: contextcheck
	if err := r.client.RemoteOptionDef4Set(config.Hostname.ValueString(), def); err != nil {
		resp.Diagnostics.AddError(
			"RemoteOptionDef4Update",
			fmt.Sprintf("Unable to update remote-option-def4 in Kea, got error: %s | %v", err, def),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Delete : Deletes an existing resource.
func (r *remoteOptionDef4Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config remoteOptionDef4ResourceSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Del", "`hostname` field is required")
	}

	//  If the Code value is empty, add an error to the diagnostics.
	if config.Code.IsNull() || config.Code.IsUnknown() {
		resp.Diagnostics.AddError("RemoteOptionDef4Del", "`code` field is required")
	}

	//  If the Space value is empty, add an error to the diagnostics.
	if config.Space.IsNull() || config.Space.IsUnknown() || config.Space.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteOptionDef4Del", "`space` field is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	if err := r.client.RemoteOptionDef4Del(config.Hostname.ValueString(), config.Space.ValueString(), int(config.Code.ValueInt64())); err != nil {
		resp.Diagnostics.AddError(
			"RemoteOptionDef4Del",
			fmt.Sprintf("Unable to delete remote-option-def4, got error: %s", err),
		)
		return
	}
}

// ImportState : Imports an existing resource by a unique identifier.
func (r *remoteOptionDef4Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("code"), req, resp)
}

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

var (
	// Ensure provider defined types fully satisfy framework interfaces.
	_ resource.Resource                = &reservationResource{}
	_ resource.ResourceWithImportState = &reservationResource{}
)

// NewReservationResource : Creates a new empty resource client.
func NewReservationResource() resource.Resource {
	return &reservationResource{}
}

type (
	// remoteSubnet4Resource defines the resource implementation.
	reservationResource struct {
		client *kea.Client
	}

	// reservationResourceSchema describes the resource data model.
	reservationResourceSchema struct {
		SubnetID            types.Int64                      `tfsdk:"subnet_id"`
		Hostname            types.String                     `tfsdk:"hostname"`
		ReservationHostname types.String                     `tfsdk:"reservation_hostname"`
		BootFileName        types.String                     `tfsdk:"boot_file_name"`
		ClientID            types.String                     `tfsdk:"client_id"`
		CircuitID           types.String                     `tfsdk:"circuit_id"`
		DuID                types.String                     `tfsdk:"duid"`
		FlexID              types.String                     `tfsdk:"flex_id"`
		IPAddress           types.String                     `tfsdk:"ip_address"`
		HwAddress           types.String                     `tfsdk:"hw_address"`
		NextServer          types.String                     `tfsdk:"next_server"`
		OptionData          []reservationOptionResourceModel `tfsdk:"option_data"`
		UserContext         types.Map                        `tfsdk:"user_context"`
	}

	// reservationOptionResourceModel : Represents a single option-data entry in Kea.
	reservationOptionResourceModel struct {
		Code       types.Int64  `tfsdk:"code"`
		Data       types.String `tfsdk:"data"`
		Name       types.String `tfsdk:"name"`
		AlwaysSend types.Bool   `tfsdk:"always_send"`
	}
)

// Metadata : Returns the resource type name and supported features.
func (r *reservationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation_resource"
}

// Schema : Returns the resource schema.
func (r *reservationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reservation resource",

		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the kea server to connect to. e.g. `kea.example.com`",
				Required:            true,
			},
			"subnet_id": schema.Int64Attribute{
				MarkdownDescription: "Subnet ID of the subnet to reserve in Kea. e.g. `1921682300`",
				Required:            true,
			},
			"reservation_hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname to define this reservation. e.g. `switch.example.com`",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "IP address for this reservation.",
				Required:            true,
			},
			"hw_address": schema.StringAttribute{
				MarkdownDescription: "Hw-address/MAC address for this reservation.",
				Required:            true,
			},
			"boot_file_name": schema.StringAttribute{
				MarkdownDescription: "Boot-file-name for this reservation.",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client-Id for this reservation.",
				Optional:            true,
			},
			"circuit_id": schema.StringAttribute{
				MarkdownDescription: "Circuit-Id for this reservation.",
				Optional:            true,
			},
			"duid": schema.StringAttribute{
				MarkdownDescription: "Du-Id for this reservation.",
				Optional:            true,
			},
			"flex_id": schema.StringAttribute{
				MarkdownDescription: "Flex-Id for this reservation.",
				Optional:            true,
			},
			"next_server": schema.StringAttribute{
				MarkdownDescription: "Next-Server for this reservation.",
				Optional:            true,
			},
			"option_data": schema.ListNestedAttribute{
				MarkdownDescription: "List of option-data to configure on the pool. e.g. `[{code = 6, name = \"domain-name-servers\", data = \"8.8.8.8, 4.2.2.2\"}]`",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code":        schema.Int64Attribute{Required: true},
						"name":        schema.StringAttribute{Required: true},
						"data":        schema.StringAttribute{Required: true},
						"always_send": schema.BoolAttribute{Required: true},
					},
				},
			},
			"user_context": schema.MapAttribute{
				MarkdownDescription: "Arbitrary string data to tie to the subnet. e.g. `{site = \"AUS\", name = \"Austin, Tx\"}`",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

// Configure : Configures the resource client data and populates the client interface.
func (r *reservationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *reservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config reservationResourceSchema

	// Read Terraform configuration data into the model
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.SubnetID.IsNull() || config.SubnetID.IsUnknown() {
		resp.Diagnostics.AddError("ReservationAdd", "`subnet_id` is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`hostname` field is required")
	}

	//  If the ReservationHostname value is empty, add an error to the diagnostics.
	if config.ReservationHostname.IsNull() || config.ReservationHostname.IsUnknown() || config.ReservationHostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`reservation_hostname` field is required")
	}

	//  If the IPAddress value is empty, add an error to the diagnostics.
	if config.IPAddress.IsNull() || config.IPAddress.IsUnknown() || config.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`ip_address` field is required")
	}

	//  If the HwAddress value is empty, add an error to the diagnostics.
	if config.HwAddress.IsNull() || config.HwAddress.IsUnknown() || config.HwAddress.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`hw_address` field is required")
	}

	// If there are any diagnostics, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	resv := kea.Reservation{
		Hostname:  config.ReservationHostname.ValueString(),
		IPAddress: config.IPAddress.ValueString(),
		HwAddress: config.HwAddress.ValueString(),
		SubnetID:  int(config.SubnetID.ValueInt64()),
		OptionData: func() []kea.OptionData {
			fr := make([]kea.OptionData, 0)
			for _, o := range config.OptionData {
				code := int(o.Code.ValueInt64())
				fr = append(fr, kea.OptionData{
					Code:       &code,
					Name:       o.Name.ValueString(),
					Data:       o.Data.ValueString(),
					AlwaysSend: o.AlwaysSend.ValueBool(),
				})
			}
			return fr
		}(),
		UserContext: func() map[string]any {
			fr := make(map[string]any)
			for k, v := range config.UserContext.Elements() {
				fr[k] = v.String()
			}
			return fr
		}(),
	}
	if !config.BootFileName.IsNull() && !config.BootFileName.IsUnknown() {
		resv.BootFileName = config.BootFileName.ValueString()
	}
	if !config.ClientID.IsNull() && !config.ClientID.IsUnknown() {
		resv.ClientID = config.ClientID.ValueString()
	}
	if !config.CircuitID.IsNull() && !config.CircuitID.IsUnknown() {
		resv.CircuitID = config.CircuitID.ValueString()
	}
	if !config.DuID.IsNull() && !config.DuID.IsUnknown() {
		resv.DuID = config.DuID.ValueString()
	}
	if !config.FlexID.IsNull() && !config.FlexID.IsUnknown() {
		resv.FlexID = config.FlexID.ValueString()
	}
	if !config.NextServer.IsNull() && !config.NextServer.IsUnknown() {
		resv.NextServer = config.NextServer.ValueString()
	}

	// nolint: contextcheck
	if err := r.client.ReservationAdd(config.Hostname.ValueString(), resv); err != nil {
		resp.Diagnostics.AddError(
			"ReservationAdd",
			fmt.Sprintf("Unable to create reservation in Kea, got error: %s | %v", err, resv),
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
func (r *reservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config reservationResourceSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.SubnetID.IsNull() || config.SubnetID.IsUnknown() {
		resp.Diagnostics.AddError("ReservationAdd", "`subnet_id` is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`hostname` field is required")
	}

	//  If the ReservationHostname value is empty, add an error to the diagnostics.
	if config.ReservationHostname.IsNull() || config.ReservationHostname.IsUnknown() || config.ReservationHostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationAdd", "`reservation_hostname` field is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	respData, err := r.client.ReservationGet(
		config.Hostname.ValueString(),
		func() string {
			if !config.IPAddress.IsNull() && !config.IPAddress.IsUnknown() && config.IPAddress.ValueString() != "" {
				return config.IPAddress.ValueString()
			}
			return config.HwAddress.ValueString()
		}(),
		int(config.SubnetID.ValueInt64()),
	)
	if err != nil {
		// Only return an error if the error is NOT subnet not found.
		if !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"ReservationGet",
				fmt.Sprintf("Unable to read example, got error: %s", err),
			)
			return
		}
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Marshalling the response data taken from Kea, and write
	// it into the TF Subnets model.
	nextServer := ""
	if respData.NextServer != "0.0.0.0" {
		nextServer = respData.NextServer
	}
	config.ReservationHostname = types.StringValue(respData.Hostname)
	config.BootFileName = types.StringValue(respData.BootFileName)
	config.ClientID = types.StringValue(respData.ClientID)
	config.CircuitID = types.StringValue(respData.CircuitID)
	config.DuID = types.StringValue(respData.DuID)
	config.FlexID = types.StringValue(respData.FlexID)
	config.IPAddress = types.StringValue(respData.IPAddress)
	config.HwAddress = types.StringValue(respData.HwAddress)
	config.NextServer = types.StringValue(nextServer)

	config.OptionData = func() []reservationOptionResourceModel {
		r := make([]reservationOptionResourceModel, 0)
		for _, v := range respData.OptionData {
			code := 0
			if v.Code != nil {
				code = *v.Code
			}
			r = append(r, reservationOptionResourceModel{
				Code:       types.Int64Value(int64(code)),
				Data:       types.StringValue(v.Data),
				Name:       types.StringValue(v.Name),
				AlwaysSend: types.BoolValue(v.AlwaysSend),
			})
		}
		return r
	}()
	if respData.UserContext != nil {
		config.UserContext = func() types.Map {
			fr := make(map[string]attr.Value)
			for k, v := range respData.UserContext {
				fr[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
			mv, diags := types.MapValue(types.StringType, fr)
			resp.Diagnostics.Append(diags...)
			return mv
		}()
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Update : Updates an existing resource.
func (r *reservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config reservationResourceSchema

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.SubnetID.IsNull() || config.SubnetID.IsUnknown() {
		resp.Diagnostics.AddError("ReservationUpdate", "`subnet_id` is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationUpdate", "`hostname` field is required")
	}

	//  If the ReservationHostname value is empty, add an error to the diagnostics.
	if config.ReservationHostname.IsNull() || config.ReservationHostname.IsUnknown() || config.ReservationHostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationUpdate", "`reservation_hostname` field is required")
	}

	//  If the IPAddress value is empty, add an error to the diagnostics.
	if config.IPAddress.IsNull() || config.IPAddress.IsUnknown() || config.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationUpdate", "`ip_address` field is required")
	}

	//  If the HwAddress value is empty, add an error to the diagnostics.
	if config.HwAddress.IsNull() || config.HwAddress.IsUnknown() || config.HwAddress.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationUpdate", "`hw_address` field is required")
	}

	// If there are any diagnostics, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	resv := kea.Reservation{
		Hostname:  config.ReservationHostname.ValueString(),
		IPAddress: config.IPAddress.ValueString(),
		HwAddress: config.HwAddress.ValueString(),
		SubnetID:  int(config.SubnetID.ValueInt64()),
		OptionData: func() []kea.OptionData {
			fr := make([]kea.OptionData, 0)
			for _, o := range config.OptionData {
				code := int(o.Code.ValueInt64())
				fr = append(fr, kea.OptionData{
					Code:       &code,
					Name:       o.Name.ValueString(),
					Data:       o.Data.ValueString(),
					AlwaysSend: o.AlwaysSend.ValueBool(),
				})
			}
			return fr
		}(),
		UserContext: func() map[string]any {
			fr := make(map[string]any)
			for k, v := range config.UserContext.Elements() {
				fr[k] = v.String()
			}
			return fr
		}(),
	}
	if !config.BootFileName.IsNull() && !config.BootFileName.IsUnknown() {
		resv.BootFileName = config.BootFileName.ValueString()
	}
	if !config.ClientID.IsNull() && !config.ClientID.IsUnknown() {
		resv.ClientID = config.ClientID.ValueString()
	}
	if !config.CircuitID.IsNull() && !config.CircuitID.IsUnknown() {
		resv.CircuitID = config.CircuitID.ValueString()
	}
	if !config.DuID.IsNull() && !config.DuID.IsUnknown() {
		resv.DuID = config.DuID.ValueString()
	}
	if !config.FlexID.IsNull() && !config.FlexID.IsUnknown() {
		resv.FlexID = config.FlexID.ValueString()
	}
	if !config.NextServer.IsNull() && !config.NextServer.IsUnknown() {
		resv.NextServer = config.NextServer.ValueString()
	}

	// nolint: contextcheck
	if err := r.client.ReservationUpdate(config.Hostname.ValueString(), resv); err != nil {
		resp.Diagnostics.AddError(
			"ReservationUpdate",
			fmt.Sprintf("Unable to update reservation in Kea, got error: %s | %v", err, resv),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Delete : Deletes an existing resource.
func (r *reservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config reservationResourceSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.SubnetID.IsNull() || config.SubnetID.IsUnknown() {
		resp.Diagnostics.AddError("ReservationDel", "`subnet_id` is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationDel", "`hostname` field is required")
	}

	//  If the ReservationHostname value is empty, add an error to the diagnostics.
	if config.ReservationHostname.IsNull() || config.ReservationHostname.IsUnknown() || config.ReservationHostname.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationDel", "`reservation_hostname` field is required")
	}

	//  If the IPAddress value is empty, add an error to the diagnostics.
	if config.IPAddress.IsNull() || config.IPAddress.IsUnknown() || config.IPAddress.ValueString() == "" {
		resp.Diagnostics.AddError("ReservationDel", "`ip_address` field is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	if err := r.client.ReservationDel(config.Hostname.ValueString(), config.IPAddress.ValueString(), int(config.SubnetID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError(
			"ReservationDel",
			fmt.Sprintf("Unable to delete reservation, got error: %s", err),
		)
		return
	}
}

// ImportState : Imports an existing resource by a unique identifier.
func (r *reservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("subnet"), req, resp)
}

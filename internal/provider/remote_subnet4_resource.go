package provider

import (
	"context"
	"fmt"
	"strconv"
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
	_            resource.Resource                = &remoteSubnet4Resource{}
	_            resource.ResourceWithImportState = &remoteSubnet4Resource{}
	cidrToIDRepl                                  = strings.NewReplacer(".", "", "/", "", " ", "")
)

// NewRemoteSubnet4Resource : Creates a new empty resource client.
func NewRemoteSubnet4Resource() resource.Resource {
	return &remoteSubnet4Resource{}
}

type (
	// remoteSubnet4Resource defines the resource implementation.
	remoteSubnet4Resource struct {
		client *kea.Client
	}

	// remoteSubnet4ResourceSchema describes the resource data model.
	remoteSubnet4ResourceSchema struct {
		Hostname   types.String                       `tfsdk:"hostname"`
		ID         types.Int64                        `tfsdk:"id"`
		OptionData []remoteSubnet4OptionResourceModel `tfsdk:"option_data"`
		Pools      []remoteSubnet4PoolResourceModel   `tfsdk:"pools"`
		Relay      []remoteSubnet4RelayResourceModel  `tfsdk:"relay"`
		Subnet     types.String                       `tfsdk:"subnet"`
	}

	// remoteSubnet4OptionResourceModel : Represents a single option-data entry in Kea.
	remoteSubnet4OptionResourceModel struct {
		Code       types.Int64  `tfsdk:"code"`
		Data       types.String `tfsdk:"data"`
		Name       types.String `tfsdk:"name"`
		AlwaysSend types.Bool   `tfsdk:"always_send"`
	}

	// remoteSubnet4PoolResourceModel : Represents a single pool entry in Kea.
	remoteSubnet4PoolResourceModel struct {
		Pool types.String `tfsdk:"pool"`
	}

	// remoteSubnet4RelayResourceModel : Represents a single ip-address relay entry in Kea.
	remoteSubnet4RelayResourceModel struct {
		IPAddress types.String `tfsdk:"ip_address"`
	}
)

// Metadata : Returns the resource type name and supported features.
func (r *remoteSubnet4Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_subnet4_resource"
}

// Schema : Returns the resource schema.
func (r *remoteSubnet4Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Remote Subnet4 resource",

		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the kea server to connect to. e.g. `kea.example.com`",
				Required:            true,
			},
			"id": schema.Int64Attribute{Computed: true},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "Subnet4 prefix to configure in Kea. e.g. `192.168.230.0/24`",
				Required:            true,
			},
			"pools": schema.ListNestedAttribute{
				MarkdownDescription: "List of pools to configure in the subnet. e.g. `['192.168.230.10-192.168.230.200']",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"pool": schema.StringAttribute{Required: true},
					},
				},
			},
			"relay": schema.ListNestedAttribute{
				MarkdownDescription: "List of relay IPs to configure in Kea. e.g. `['192.168.230.1']`",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.StringAttribute{Required: true},
					},
				},
			},
			"option_data": schema.ListNestedAttribute{
				MarkdownDescription: "List of option-data to configure on the pool. e.g. `[{code = 6, name = \"domain-name-servers\", data = \"8.8.8.8, 4.2.2.2\"}]`",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code":        schema.Int64Attribute{Required: true},
						"name":        schema.StringAttribute{Required: true},
						"data":        schema.StringAttribute{Required: true},
						"always_send": schema.BoolAttribute{Optional: true},
					},
				},
			},
		},
	}
}

// Configure : Configures the resource client data and populates the client interface.
func (r *remoteSubnet4Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *remoteSubnet4Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config remoteSubnet4ResourceSchema

	// Read Terraform configuration data into the model
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.Subnet.IsNull() || config.Subnet.IsUnknown() || config.Subnet.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Subnet4 prefix is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Hostname is required")
	}

	// If there are any diagnostics, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	subnet := config.Subnet.ValueString()
	id, err := strconv.Atoi(cidrToIDRepl.Replace(strings.Split(subnet, "/")[0]))
	if err != nil {
		resp.Diagnostics.AddError(
			"RemoteSubnet4Create",
			fmt.Sprintf("Unable to parse subnet4 prefix `%s`, got error: %s", subnet, err),
		)
		return
	}

	newSubnet := kea.NewRemoteSubnet4{
		ID:     id,
		Subnet: config.Subnet.ValueString(),
		Pools: func() []kea.Pool {
			fr := make([]kea.Pool, 0)
			for _, p := range config.Pools {
				fr = append(fr, kea.Pool{Pool: p.Pool.ValueString()})
			}
			return fr
		}(),
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
		Relay: func() kea.Relay {
			fr := kea.Relay{}
			if len(config.Relay) == 0 {
				return fr
			}
			for _, ip := range config.Relay {
				fr.IPAddresses = append(fr.IPAddresses, ip.IPAddress.ValueString())
			}
			return fr
		}(),
	}

	// nolint: contextcheck
	respData, err := r.client.RemoteSubnet4Set(config.Hostname.ValueString(), []kea.NewRemoteSubnet4{newSubnet})
	if err != nil {
		resp.Diagnostics.AddError(
			"RemoteSubnet4Create",
			fmt.Sprintf("Unable to create subnet with new id=%d in Kea, got error: %s | %v", id, err, newSubnet),
		)
		return
	}
	if len(respData) != 1 {
		resp.Diagnostics.AddError(
			"RemoteSubnet4Update",
			fmt.Sprintf("Unable to create subnet4 with new id=%d in Kea,, got error: %s", id, err),
		)
		return
	}

	// Marshalling the response data taken from Kea, and write
	// it into the TF Subnets model.
	res := respData[0]
	config.ID = types.Int64Value(int64(res.ID))
	config.Subnet = types.StringValue(res.Subnet)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Read : Reads the resource data into the Terraform state.
func (r *remoteSubnet4Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var config remoteSubnet4ResourceSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.Subnet.IsNull() || config.Subnet.IsUnknown() || config.Subnet.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Subnet4 prefix is required")
	}
	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Hostname is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	respData, err := r.client.RemoteSubnet4GetByPrefix(config.Hostname.ValueString(), config.Subnet.ValueString())
	if err != nil {
		// Only return an error if the error is NOT subnet not found.
		if !strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"RemoteSubnet4GetByPrefix",
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
	config.ID = types.Int64Value(int64(respData.ID))
	config.OptionData = func() []remoteSubnet4OptionResourceModel {
		r := make([]remoteSubnet4OptionResourceModel, 0)
		for _, v := range respData.OptionData {
			code := 0
			if v.Code != nil {
				code = *v.Code
			}
			r = append(r, remoteSubnet4OptionResourceModel{
				Code:       types.Int64Value(int64(code)),
				Data:       types.StringValue(v.Data),
				Name:       types.StringValue(v.Name),
				AlwaysSend: types.BoolValue(v.AlwaysSend),
			})
		}
		return r
	}()
	config.Pools = func() []remoteSubnet4PoolResourceModel {
		fr := make([]remoteSubnet4PoolResourceModel, 0)
		for _, v := range respData.Pools {
			fr = append(fr, remoteSubnet4PoolResourceModel{Pool: types.StringValue(v.Pool)})
		}
		return fr
	}()
	config.Relay = func() []remoteSubnet4RelayResourceModel {
		fr := make([]remoteSubnet4RelayResourceModel, 0)
		for _, v := range respData.Relay.IPAddresses {
			fr = append(fr, remoteSubnet4RelayResourceModel{IPAddress: types.StringValue(v)})
		}
		return fr
	}()
	config.Subnet = types.StringValue(respData.Subnet)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Update : Updates an existing resource.
func (r *remoteSubnet4Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config remoteSubnet4ResourceSchema

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.Subnet.IsNull() || config.Subnet.IsUnknown() || config.Subnet.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Subnet4 prefix is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Hostname is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	subnet := config.Subnet.ValueString()
	id, err := strconv.Atoi(cidrToIDRepl.Replace(strings.Split(subnet, "/")[0]))
	if err != nil {
		resp.Diagnostics.AddError(
			"RemoteSubnet4Create",
			fmt.Sprintf("Unable to parse subnet4 prefix `%s`, got error: %s", subnet, err),
		)
		return
	}

	// nolint: contextcheck
	respData, err := r.client.RemoteSubnet4Set(
		config.Hostname.ValueString(),
		[]kea.NewRemoteSubnet4{
			{
				ID:     id,
				Subnet: config.Subnet.ValueString(),
				Pools: func() []kea.Pool {
					fr := make([]kea.Pool, 0)
					for _, p := range config.Pools {
						fr = append(fr, kea.Pool{Pool: p.Pool.ValueString()})
					}
					return fr
				}(),
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
				Relay: func() kea.Relay {
					fr := kea.Relay{}
					if len(config.Relay) == 0 {
						return fr
					}
					for _, ip := range config.Relay {
						fr.IPAddresses = append(fr.IPAddresses, ip.IPAddress.ValueString())
					}
					return fr
				}(),
			},
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"RemoteSubnet4Update",
			fmt.Sprintf("Unable to update subnet4 in Kea, got error: %s", err),
		)
		return
	}
	if len(respData) != 1 {
		resp.Diagnostics.AddError(
			"RemoteSubnet4Update",
			fmt.Sprintf("Unable to update subnet4 in Kea,, got error: %s", err),
		)
		return
	}

	// Marshalling the response data taken from Kea, and write
	// it into the TF Subnets model.
	res := respData[0]
	config.ID = types.Int64Value(int64(res.ID))
	config.Subnet = types.StringValue(res.Subnet)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Delete : Deletes an existing resource.
func (r *remoteSubnet4Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var config remoteSubnet4ResourceSchema

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &config)...)

	// If the subnet value is empty, add an error to the diagnostics.
	if config.Subnet.IsNull() || config.Subnet.IsUnknown() || config.Subnet.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Subnet4 prefix is required")
	}

	//  If the hostname value is empty, add an error to the diagnostics.
	if config.Hostname.IsNull() || config.Hostname.IsUnknown() || config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError("RemoteSubnet4Read", "Hostname is required")
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	if _, err := r.client.RemoteSubnet4DelByPrefix(config.Hostname.ValueString(), config.Subnet.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"RemoteSubnet4DelByPrefix",
			fmt.Sprintf("Unable to delete prefix, got error: %s", err),
		)
		return
	}
}

// ImportState : Imports an existing resource by a unique identifier.
func (r *remoteSubnet4Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("subnet"), req, resp)
}

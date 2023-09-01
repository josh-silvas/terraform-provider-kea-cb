// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

var (
	// Ensure provider defined types fully satisfy framework interfaces.
	_ datasource.DataSource              = &reservationDataSource{}
	_ datasource.DataSourceWithConfigure = &reservationDataSource{}
)

// NewReservationDataSource : Creates a new empty data source client.
func NewReservationDataSource() datasource.DataSource {
	return &reservationDataSource{}
}

type (
	// reservationDataSource defines the data source client.
	reservationDataSource struct {
		client *kea.Client
	}

	// reservationDataSourceSchema describes the data source data model.
	// Maps to the source schema data.
	reservationDataSourceSchema struct {
		IPOrMac             types.String                       `tfsdk:"ip_or_mac_address"`
		SubnetID            types.Int64                        `tfsdk:"subnet_id"`
		Hostname            types.String                       `tfsdk:"hostname"`
		ReservationHostname types.String                       `tfsdk:"reservation_hostname"`
		BootFileName        types.String                       `tfsdk:"boot_file_name"`
		ClientID            types.String                       `tfsdk:"client_id"`
		CircuitID           types.String                       `tfsdk:"circuit_id"`
		DuID                types.String                       `tfsdk:"duid"`
		FlexID              types.String                       `tfsdk:"flex_id"`
		IPAddress           types.String                       `tfsdk:"ip_address"`
		HwAddress           types.String                       `tfsdk:"hw_address"`
		NextServer          types.String                       `tfsdk:"next_server"`
		OptionData          []reservationDataSourceOptionModel `tfsdk:"option_data"`
		UserContext         types.Map                          `tfsdk:"user_context"`
	}

	// reservationDataSourceOptionModel : Represents a single option-data entry in Kea.
	reservationDataSourceOptionModel struct {
		Code       types.Int64  `tfsdk:"code"`
		Data       types.String `tfsdk:"data"`
		Name       types.String `tfsdk:"name"`
		AlwaysSend types.Bool   `tfsdk:"always_send"`
	}
)

// Metadata : Defines the data source metadata.
func (d *reservationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation_data_source"
}

// Schema : Defines the data source schema.
func (d *reservationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reservation data source",
		Attributes: map[string]schema.Attribute{
			"subnet_id": schema.Int64Attribute{
				MarkdownDescription: "Subnet4 ID to fetch the reservations from. e.g. 1921682300`",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the kea server to connect to. e.g. `kea.example.com`",
				Required:            true,
			},
			"ip_or_mac_address": schema.StringAttribute{
				MarkdownDescription: "IP address or mac-address to fetch for this reservation. e.g. 192.168.230.50`",
				Required:            true,
			},
			"reservation_hostname": schema.StringAttribute{Computed: true},
			"boot_file_name":       schema.StringAttribute{Computed: true},
			"client_id":            schema.StringAttribute{Computed: true},
			"circuit_id":           schema.StringAttribute{Computed: true},
			"duid":                 schema.StringAttribute{Computed: true},
			"flex_id":              schema.StringAttribute{Computed: true},
			"ip_address":           schema.StringAttribute{Computed: true},
			"hw_address":           schema.StringAttribute{Computed: true},
			"next_server":          schema.StringAttribute{Computed: true},
			"option_data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code":        schema.Int64Attribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"data":        schema.StringAttribute{Computed: true},
						"always_send": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"user_context": schema.MapAttribute{ElementType: types.StringType, Computed: true},
		},
	}
}

// Configure : Configures the data source client.
func (d *reservationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *reservationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Define an empty configuration.
	var config reservationDataSourceSchema

	// Read Terraform configuration data into the model
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// Validate that only one of `prefix` or `subnet_id` is specified.
	if config.IPOrMac.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"`ip_or_mac_address` must be specified.",
		)
	}

	// Validate that a `hostname` is specified.
	if config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"A `hostname` must be specified. DNS name or IP address of the Kea DHCP server.",
		)
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// nolint: contextcheck
	respData, err := d.client.ReservationGet(
		config.Hostname.ValueString(),
		config.IPOrMac.ValueString(),
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

	config.ReservationHostname = types.StringValue(respData.Hostname)
	config.BootFileName = types.StringValue(respData.BootFileName)
	config.ClientID = types.StringValue(respData.ClientID)
	config.CircuitID = types.StringValue(respData.CircuitID)
	config.DuID = types.StringValue(respData.DuID)
	config.FlexID = types.StringValue(respData.FlexID)
	config.IPAddress = types.StringValue(respData.IPAddress)
	config.HwAddress = types.StringValue(respData.HwAddress)
	config.NextServer = types.StringValue(respData.NextServer)

	config.OptionData = func() []reservationDataSourceOptionModel {
		r := make([]reservationDataSourceOptionModel, 0)
		for _, v := range respData.OptionData {
			code := 0
			if v.Code != nil {
				code = *v.Code
			}
			r = append(r, reservationDataSourceOptionModel{
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

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientNzbvortexResourceName   = "download_client_nzbvortex"
	downloadClientNzbvortexImplementation = "Nzbvortex"
	downloadClientNzbvortexConfigContract = "NzbvortexSettings"
	downloadClientNzbvortexProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientNzbvortexResource{}
	_ resource.ResourceWithImportState = &DownloadClientNzbvortexResource{}
)

func NewDownloadClientNzbvortexResource() resource.Resource {
	return &DownloadClientNzbvortexResource{}
}

// DownloadClientNzbvortexResource defines the download client implementation.
type DownloadClientNzbvortexResource struct {
	client *whisparr.APIClient
}

// DownloadClientNzbvortex describes the download client data model.
type DownloadClientNzbvortex struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	APIKey                   types.String `tfsdk:"api_key"`
	MovieCategory            types.String `tfsdk:"movie_category"`
	RecentMoviePriority      types.Int64  `tfsdk:"recent_movie_priority"`
	OlderMoviePriority       types.Int64  `tfsdk:"older_movie_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientNzbvortex) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		APIKey:                   d.APIKey,
		MovieCategory:            d.MovieCategory,
		RecentMoviePriority:      d.RecentMoviePriority,
		OlderMoviePriority:       d.OlderMoviePriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
	}
}

func (d *DownloadClientNzbvortex) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.APIKey = client.APIKey
	d.MovieCategory = client.MovieCategory
	d.RecentMoviePriority = client.RecentMoviePriority
	d.OlderMoviePriority = client.OlderMoviePriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientNzbvortexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientNzbvortexResourceName
}

func (r *DownloadClientNzbvortexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Nzbvortex resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/whisparr/settings#download-clients) and [Nzbvortex](https://wiki.servarr.com/whisparr/supported#nzbvortex).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_completed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_failed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Movie priority. `-1` Low, `0` Normal, `1` High.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-1, 0, 1),
				},
			},
			"older_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Movie priority. `-1` Low, `0` Normal, `1` High.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-1, 0, 1),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
			},
			"movie_category": schema.StringAttribute{
				MarkdownDescription: "Movie category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientNzbvortexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientNzbvortexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientNzbvortex
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientNzbvortex

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientNzbvortex current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientNzbvortex
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbvortexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientNzbvortex

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientNzbvortex current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientNzbvortexResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientNzbvortexResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientNzbvortexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientNzbvortexResourceName+": "+req.ID)
}

func (d *DownloadClientNzbvortex) write(ctx context.Context, downloadClient *whisparr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.GetEnable()),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.GetRemoveCompletedDownloads()),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.GetRemoveFailedDownloads()),
		Priority:                 types.Int64Value(int64(downloadClient.GetPriority())),
		ID:                       types.Int64Value(int64(downloadClient.GetId())),
		Name:                     types.StringValue(downloadClient.GetName()),
	}
	genericDownloadClient.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, downloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.GetFields())
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientNzbvortex) read(ctx context.Context) *whisparr.DownloadClientResource {
	tags := make([]*int32, len(d.Tags.Elements()))
	tfsdk.ValueAs(ctx, d.Tags, &tags)

	client := whisparr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetRemoveCompletedDownloads(d.RemoveCompletedDownloads.ValueBool())
	client.SetRemoveFailedDownloads(d.RemoveFailedDownloads.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(downloadClientNzbvortexConfigContract)
	client.SetImplementation(downloadClientNzbvortexImplementation)
	client.SetName(d.Name.ValueString())
	client.SetProtocol(downloadClientNzbvortexProtocol)
	client.SetTags(tags)
	client.SetFields(d.toDownloadClient().readFields(ctx))

	return client
}

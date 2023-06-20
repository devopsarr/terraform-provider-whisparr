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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientNzbgetResourceName   = "download_client_nzbget"
	downloadClientNzbgetImplementation = "Nzbget"
	downloadClientNzbgetConfigContract = "NzbgetSettings"
	downloadClientNzbgetProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientNzbgetResource{}
	_ resource.ResourceWithImportState = &DownloadClientNzbgetResource{}
)

func NewDownloadClientNzbgetResource() resource.Resource {
	return &DownloadClientNzbgetResource{}
}

// DownloadClientNzbgetResource defines the download client implementation.
type DownloadClientNzbgetResource struct {
	client *whisparr.APIClient
}

// DownloadClientNzbget describes the download client data model.
type DownloadClientNzbget struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	MovieCategory            types.String `tfsdk:"movie_category"`
	RecentMoviePriority      types.Int64  `tfsdk:"recent_movie_priority"`
	OlderMoviePriority       types.Int64  `tfsdk:"older_movie_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	AddPaused                types.Bool   `tfsdk:"add_paused"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientNzbget) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Username:                 d.Username,
		Password:                 d.Password,
		MovieCategory:            d.MovieCategory,
		RecentMoviePriority:      d.RecentMoviePriority,
		OlderMoviePriority:       d.OlderMoviePriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		AddPaused:                d.AddPaused,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		Implementation:           types.StringValue(downloadClientNzbgetImplementation),
		ConfigContract:           types.StringValue(downloadClientNzbgetConfigContract),
		Protocol:                 types.StringValue(downloadClientNzbgetProtocol),
	}
}

func (d *DownloadClientNzbget) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.MovieCategory = client.MovieCategory
	d.RecentMoviePriority = client.RecentMoviePriority
	d.OlderMoviePriority = client.OlderMoviePriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.AddPaused = client.AddPaused
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientNzbgetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientNzbgetResourceName
}

func (r *DownloadClientNzbgetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client NZBGet resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/whisparr/settings#download-clients) and [NZBGet](https://wiki.servarr.com/whisparr/supported#nzbget).",
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
			"add_paused": schema.BoolAttribute{
				MarkdownDescription: "Add paused flag.",
				Optional:            true,
				Computed:            true,
			},
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
				MarkdownDescription: "Recent Movie priority. `-100` VeryLow, `-50` Low, `0` Normal, `50` High, `100` VeryHigh, `900` Force.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-100, -50, 0, 50, 100, 900),
				},
			},
			"older_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Movie priority. `-100` VeryLow, `-50` Low, `0` Normal, `50` High, `100` VeryHigh, `900` Force.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-100, -50, 0, 50, 100, 900),
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
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"movie_category": schema.StringAttribute{
				MarkdownDescription: "Movie category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientNzbgetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientNzbgetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientNzbget

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientNzbget
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientNzbgetResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientNzbgetResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbgetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientNzbget

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientNzbget current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientNzbgetResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientNzbgetResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbgetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientNzbget

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientNzbget
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientNzbgetResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientNzbgetResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientNzbgetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientNzbget

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientNzbget current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientNzbgetResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientNzbgetResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientNzbgetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientNzbgetResourceName+": "+req.ID)
}

func (d *DownloadClientNzbget) write(ctx context.Context, downloadClient *whisparr.DownloadClientResource) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientNzbget) read(ctx context.Context) *whisparr.DownloadClientResource {
	return d.toDownloadClient().read(ctx)
}

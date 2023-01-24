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
	downloadClientQbittorrentResourceName   = "download_client_qbittorrent"
	downloadClientQbittorrentImplementation = "QBittorrent"
	downloadClientQbittorrentConfigContract = "QBittorrentSettings"
	downloadClientQbittorrentProtocol       = "torrent"
)

var downloadClientQbittorrentInitialStates = []int64{0, 1, 2}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientQbittorrentResource{}
	_ resource.ResourceWithImportState = &DownloadClientQbittorrentResource{}
)

func NewDownloadClientQbittorrentResource() resource.Resource {
	return &DownloadClientQbittorrentResource{}
}

// DownloadClientQbittorrentResource defines the download client implementation.
type DownloadClientQbittorrentResource struct {
	client *whisparr.APIClient
}

// DownloadClientQbittorrent describes the download client data model.
type DownloadClientQbittorrent struct {
	Tags                     types.Set    `tfsdk:"tags"`
	MovieImportedCategory    types.String `tfsdk:"movie_imported_category"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	MovieCategory            types.String `tfsdk:"movie_category"`
	MovieDirectory           types.String `tfsdk:"movie_directory"`
	RecentMoviePriority      types.Int64  `tfsdk:"recent_movie_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	OlderMoviePriority       types.Int64  `tfsdk:"older_movie_priority"`
	InitialState             types.Int64  `tfsdk:"initial_state"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
	FirstAndLast             types.Bool   `tfsdk:"first_and_last"`
	SequentialOrder          types.Bool   `tfsdk:"sequential_order"`
}

func (d DownloadClientQbittorrent) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Username:                 d.Username,
		Password:                 d.Password,
		MovieCategory:            d.MovieCategory,
		MovieDirectory:           d.MovieDirectory,
		RecentMoviePriority:      d.RecentMoviePriority,
		OlderMoviePriority:       d.OlderMoviePriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		MovieImportedCategory:    d.MovieImportedCategory,
		InitialState:             d.InitialState,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		FirstAndLast:             d.FirstAndLast,
		SequentialOrder:          d.SequentialOrder,
	}
}

func (d *DownloadClientQbittorrent) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.MovieCategory = client.MovieCategory
	d.MovieDirectory = client.MovieDirectory
	d.RecentMoviePriority = client.RecentMoviePriority
	d.OlderMoviePriority = client.OlderMoviePriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.MovieImportedCategory = client.MovieImportedCategory
	d.InitialState = client.InitialState
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
	d.FirstAndLast = client.FirstAndLast
	d.SequentialOrder = client.SequentialOrder
}

func (r *DownloadClientQbittorrentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientQbittorrentResourceName
}

func (r *DownloadClientQbittorrentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client qBittorrent resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/whisparr/settings#download-clients) and [qBittorrent](https://wiki.servarr.com/whisparr/supported#qbittorrent).",
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
			"sequential_order": schema.BoolAttribute{
				MarkdownDescription: "Sequential order flag.",
				Optional:            true,
				Computed:            true,
			},
			"first_and_last": schema.BoolAttribute{
				MarkdownDescription: "First and last flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"older_movie_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Movie priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"initial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(downloadClientQbittorrentInitialStates...),
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
			"movie_imported_category": schema.StringAttribute{
				MarkdownDescription: "Movie imported category.",
				Optional:            true,
				Computed:            true,
			},
			"movie_directory": schema.StringAttribute{
				MarkdownDescription: "Movie directory.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientQbittorrentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientQbittorrentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientQbittorrent

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientQbittorrent
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientQbittorrentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientQbittorrent

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientQbittorrent current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientQbittorrentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientQbittorrent

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientQbittorrent
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientQbittorrentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientQbittorrent

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientQbittorrent current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientQbittorrentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientQbittorrentResourceName+": "+req.ID)
}

func (d *DownloadClientQbittorrent) write(ctx context.Context, downloadClient *whisparr.DownloadClientResource) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.GetEnable()),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.GetRemoveCompletedDownloads()),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.GetRemoveFailedDownloads()),
		Priority:                 types.Int64Value(int64(downloadClient.GetPriority())),
		ID:                       types.Int64Value(int64(downloadClient.GetId())),
		Name:                     types.StringValue(downloadClient.GetName()),
	}
	genericDownloadClient.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, downloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientQbittorrent) read(ctx context.Context) *whisparr.DownloadClientResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	client := whisparr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetRemoveCompletedDownloads(d.RemoveCompletedDownloads.ValueBool())
	client.SetRemoveFailedDownloads(d.RemoveFailedDownloads.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(downloadClientQbittorrentConfigContract)
	client.SetImplementation(downloadClientQbittorrentImplementation)
	client.SetName(d.Name.ValueString())
	client.SetProtocol(downloadClientQbittorrentProtocol)
	client.SetTags(tags)
	client.SetFields(d.toDownloadClient().readFields(ctx))

	return client
}

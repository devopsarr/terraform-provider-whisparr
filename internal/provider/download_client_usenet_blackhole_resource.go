package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientUsenetBlackholeResourceName   = "download_client_usenet_blackhole"
	downloadClientUsenetBlackholeImplementation = "UsenetBlackhole"
	downloadClientUsenetBlackholeConfigContract = "UsenetBlackholeSettings"
	downloadClientUsenetBlackholeProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientUsenetBlackholeResource{}
	_ resource.ResourceWithImportState = &DownloadClientUsenetBlackholeResource{}
)

func NewDownloadClientUsenetBlackholeResource() resource.Resource {
	return &DownloadClientUsenetBlackholeResource{}
}

// DownloadClientUsenetBlackholeResource defines the download client implementation.
type DownloadClientUsenetBlackholeResource struct {
	client *whisparr.APIClient
}

// DownloadClientUsenetBlackhole describes the download client data model.
type DownloadClientUsenetBlackhole struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	NzbFolder                types.String `tfsdk:"nzb_folder"`
	WatchFolder              types.String `tfsdk:"watch_folder"`
	Priority                 types.Int64  `tfsdk:"priority"`
	ID                       types.Int64  `tfsdk:"id"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientUsenetBlackhole) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		NzbFolder:                d.NzbFolder,
		WatchFolder:              d.WatchFolder,
		Priority:                 d.Priority,
		ID:                       d.ID,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
	}
}

func (d *DownloadClientUsenetBlackhole) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.NzbFolder = client.NzbFolder
	d.WatchFolder = client.WatchFolder
	d.Priority = client.Priority
	d.ID = client.ID
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientUsenetBlackholeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientUsenetBlackholeResourceName
}

func (r *DownloadClientUsenetBlackholeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Usenet Blackhole resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/whisparr/settings#download-clients) and [UsenetBlackhole](https://wiki.servarr.com/whisparr/supported#usenetblackhole).",
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
			"nzb_folder": schema.StringAttribute{
				MarkdownDescription: "Usenet folder.",
				Required:            true,
			},
			"watch_folder": schema.StringAttribute{
				MarkdownDescription: "Watch folder flag.",
				Required:            true,
			},
		},
	}
}

func (r *DownloadClientUsenetBlackholeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientUsenetBlackholeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientUsenetBlackhole

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientUsenetBlackhole
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientUsenetBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientUsenetBlackholeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUsenetBlackholeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientUsenetBlackhole

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientUsenetBlackhole current value
	response, _, err := r.client.DownloadClientApi.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientUsenetBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientUsenetBlackholeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUsenetBlackholeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientUsenetBlackhole

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientUsenetBlackhole
	request := client.read(ctx)

	response, _, err := r.client.DownloadClientApi.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientUsenetBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientUsenetBlackholeResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientUsenetBlackholeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientUsenetBlackhole

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientUsenetBlackhole current value
	_, err := r.client.DownloadClientApi.DeleteDownloadClient(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientUsenetBlackholeResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientUsenetBlackholeResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientUsenetBlackholeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientUsenetBlackholeResourceName+": "+req.ID)
}

func (d *DownloadClientUsenetBlackhole) write(ctx context.Context, downloadClient *whisparr.DownloadClientResource) {
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

func (d *DownloadClientUsenetBlackhole) read(ctx context.Context) *whisparr.DownloadClientResource {
	tags := make([]*int32, len(d.Tags.Elements()))
	tfsdk.ValueAs(ctx, d.Tags, &tags)

	client := whisparr.NewDownloadClientResource()
	client.SetEnable(d.Enable.ValueBool())
	client.SetRemoveCompletedDownloads(d.RemoveCompletedDownloads.ValueBool())
	client.SetRemoveFailedDownloads(d.RemoveFailedDownloads.ValueBool())
	client.SetPriority(int32(d.Priority.ValueInt64()))
	client.SetId(int32(d.ID.ValueInt64()))
	client.SetConfigContract(downloadClientUsenetBlackholeConfigContract)
	client.SetImplementation(downloadClientUsenetBlackholeImplementation)
	client.SetName(d.Name.ValueString())
	client.SetProtocol(downloadClientUsenetBlackholeProtocol)
	client.SetTags(tags)
	client.SetFields(d.toDownloadClient().readFields(ctx))

	return client
}

package provider

import (
	"context"
	"fmt"
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
	indexerTorrentRssResourceName   = "indexer_torrent_rss"
	indexerTorrentRssImplementation = "TorrentRssIndexer"
	indexerTorrentRssConfigContract = "TorrentRssIndexerSettings"
	indexerTorrentRssProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerTorrentRssResource{}
	_ resource.ResourceWithImportState = &IndexerTorrentRssResource{}
)

func NewIndexerTorrentRssResource() resource.Resource {
	return &IndexerTorrentRssResource{}
}

// IndexerTorrentRssResource defines the TorrentRss indexer implementation.
type IndexerTorrentRssResource struct {
	client *whisparr.APIClient
}

// IndexerTorrentRss describes the TorrentRss indexer data model.
type IndexerTorrentRss struct {
	Tags             types.Set     `tfsdk:"tags"`
	MultiLanguages   types.Set     `tfsdk:"multi_languages"`
	RequiredFlags    types.Set     `tfsdk:"required_flags"`
	Name             types.String  `tfsdk:"name"`
	BaseURL          types.String  `tfsdk:"base_url"`
	Cookie           types.String  `tfsdk:"cookie"`
	Priority         types.Int64   `tfsdk:"priority"`
	ID               types.Int64   `tfsdk:"id"`
	DownloadClientID types.Int64   `tfsdk:"download_client_id"`
	MinimumSeeders   types.Int64   `tfsdk:"minimum_seeders"`
	SeedTime         types.Int64   `tfsdk:"seed_time"`
	SeedRatio        types.Float64 `tfsdk:"seed_ratio"`
	AllowZeroSize    types.Bool    `tfsdk:"allow_zero_size"`
	EnableRss        types.Bool    `tfsdk:"enable_rss"`
}

func (i IndexerTorrentRss) toIndexer() *Indexer {
	return &Indexer{
		EnableRss:        i.EnableRss,
		AllowZeroSize:    i.AllowZeroSize,
		Priority:         i.Priority,
		DownloadClientID: i.DownloadClientID,
		ID:               i.ID,
		Name:             i.Name,
		Cookie:           i.Cookie,
		MinimumSeeders:   i.MinimumSeeders,
		SeedTime:         i.SeedTime,
		SeedRatio:        i.SeedRatio,
		BaseURL:          i.BaseURL,
		Tags:             i.Tags,
		MultiLanguages:   i.MultiLanguages,
		RequiredFlags:    i.RequiredFlags,
	}
}

func (i *IndexerTorrentRss) fromIndexer(indexer *Indexer) {
	i.EnableRss = indexer.EnableRss
	i.AllowZeroSize = indexer.AllowZeroSize
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.Cookie = indexer.Cookie
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
	i.MultiLanguages = indexer.MultiLanguages
	i.RequiredFlags = indexer.RequiredFlags
}

func (r *IndexerTorrentRssResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerTorrentRssResourceName
}

func (r *IndexerTorrentRssResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Torrent RSS resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/whisparr/settings#indexers) and [Torrent RSS](https://wiki.servarr.com/whisparr/supported#torrentrssindexer).",
		Attributes: map[string]schema.Attribute{
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"download_client_id": schema.Int64Attribute{
				MarkdownDescription: "Download client ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "IndexerTorrentRss name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerTorrentRss ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"allow_zero_size": schema.BoolAttribute{
				MarkdownDescription: "Allow zero size files.",
				Optional:            true,
				Computed:            true,
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Optional:            true,
				Computed:            true,
			},
			"seed_time": schema.Int64Attribute{
				MarkdownDescription: "Seed time.",
				Optional:            true,
				Computed:            true,
			},
			"seed_ratio": schema.Float64Attribute{
				MarkdownDescription: "Seed ratio.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
			"cookie": schema.StringAttribute{
				MarkdownDescription: "Cookie.",
				Optional:            true,
				Computed:            true,
			},
			"multi_languages": schema.SetAttribute{
				MarkdownDescription: "Languages list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"required_flags": schema.SetAttribute{
				MarkdownDescription: "Flag list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerTorrentRssResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerTorrentRssResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerTorrentRss
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorrentRssResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerTorrentRss current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorrentRssResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerTorrentRss
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+indexerTorrentRssResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerTorrentRssResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerTorrentRss

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerTorrentRss current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerTorrentRssResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerTorrentRssResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerTorrentRssResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerTorrentRssResourceName+": "+req.ID)
}

func (i *IndexerTorrentRss) write(ctx context.Context, indexer *whisparr.IndexerResource) {
	genericIndexer := Indexer{
		EnableRss:        types.BoolValue(indexer.GetEnableRss()),
		Priority:         types.Int64Value(int64(indexer.GetPriority())),
		DownloadClientID: types.Int64Value(int64(indexer.GetDownloadClientId())),
		ID:               types.Int64Value(int64(indexer.GetId())),
		Name:             types.StringValue(indexer.GetName()),
	}
	genericIndexer.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, indexer.Tags)
	genericIndexer.writeFields(ctx, indexer.GetFields())
	i.fromIndexer(&genericIndexer)
}

func (i *IndexerTorrentRss) read(ctx context.Context) *whisparr.IndexerResource {
	tags := make([]*int32, len(i.Tags.Elements()))
	tfsdk.ValueAs(ctx, i.Tags, &tags)

	indexer := whisparr.NewIndexerResource()
	indexer.SetEnableRss(i.EnableRss.ValueBool())
	indexer.SetPriority(int32(i.Priority.ValueInt64()))
	indexer.SetDownloadClientId(int32(i.DownloadClientID.ValueInt64()))
	indexer.SetId(int32(i.ID.ValueInt64()))
	indexer.SetConfigContract(indexerTorrentRssConfigContract)
	indexer.SetImplementation(indexerTorrentRssImplementation)
	indexer.SetName(i.Name.ValueString())
	indexer.SetProtocol(indexerTorrentRssProtocol)
	indexer.SetTags(tags)
	indexer.SetFields(i.toIndexer().readFields(ctx))

	return indexer
}

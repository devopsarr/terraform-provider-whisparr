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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerIptorrentsResourceName   = "indexer_iptorrents"
	indexerIptorrentsImplementation = "IPTorrents"
	indexerIptorrentsConfigContract = "IPTorrentsSettings"
	indexerIptorrentsProtocol       = "torrent"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerIptorrentsResource{}
	_ resource.ResourceWithImportState = &IndexerIptorrentsResource{}
)

func NewIndexerIptorrentsResource() resource.Resource {
	return &IndexerIptorrentsResource{}
}

// IndexerIptorrentsResource defines the Iptorrents indexer implementation.
type IndexerIptorrentsResource struct {
	client *whisparr.APIClient
}

// IndexerIptorrents describes the Iptorrents indexer data model.
type IndexerIptorrents struct {
	Tags             types.Set     `tfsdk:"tags"`
	MultiLanguages   types.Set     `tfsdk:"multi_languages"`
	RequiredFlags    types.Set     `tfsdk:"required_flags"`
	Name             types.String  `tfsdk:"name"`
	BaseURL          types.String  `tfsdk:"base_url"`
	Priority         types.Int64   `tfsdk:"priority"`
	ID               types.Int64   `tfsdk:"id"`
	DownloadClientID types.Int64   `tfsdk:"download_client_id"`
	MinimumSeeders   types.Int64   `tfsdk:"minimum_seeders"`
	SeedTime         types.Int64   `tfsdk:"seed_time"`
	SeedRatio        types.Float64 `tfsdk:"seed_ratio"`
	EnableRss        types.Bool    `tfsdk:"enable_rss"`
}

func (i IndexerIptorrents) toIndexer() *Indexer {
	return &Indexer{
		EnableRss:        i.EnableRss,
		Priority:         i.Priority,
		DownloadClientID: i.DownloadClientID,
		ID:               i.ID,
		Name:             i.Name,
		MinimumSeeders:   i.MinimumSeeders,
		SeedTime:         i.SeedTime,
		SeedRatio:        i.SeedRatio,
		BaseURL:          i.BaseURL,
		Tags:             i.Tags,
		MultiLanguages:   i.MultiLanguages,
		RequiredFlags:    i.RequiredFlags,
		Implementation:   types.StringValue(indexerIptorrentsImplementation),
		ConfigContract:   types.StringValue(indexerIptorrentsConfigContract),
		Protocol:         types.StringValue(indexerIptorrentsProtocol),
	}
}

func (i *IndexerIptorrents) fromIndexer(indexer *Indexer) {
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.MinimumSeeders = indexer.MinimumSeeders
	i.SeedTime = indexer.SeedTime
	i.SeedRatio = indexer.SeedRatio
	i.BaseURL = indexer.BaseURL
	i.Tags = indexer.Tags
	i.MultiLanguages = indexer.MultiLanguages
	i.RequiredFlags = indexer.RequiredFlags
}

func (r *IndexerIptorrentsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerIptorrentsResourceName
}

func (r *IndexerIptorrentsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer IP Torrents resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/whisparr/settings#indexers) and [IP Torrents](https://wiki.servarr.com/whisparr/supported#iptorrents).",
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
				MarkdownDescription: "IndexerIptorrents name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerIptorrents ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
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

func (r *IndexerIptorrentsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerIptorrentsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerIptorrents

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerIptorrents
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerIptorrentsResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerIptorrentsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerIptorrentsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerIptorrents

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerIptorrents current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerIptorrentsResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerIptorrentsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerIptorrentsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerIptorrents

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerIptorrents
	request := indexer.read(ctx)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+indexerIptorrentsResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerIptorrentsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerIptorrentsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var indexer *IndexerIptorrents

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerIptorrents current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerIptorrentsResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerIptorrentsResourceName+": "+strconv.Itoa(int(indexer.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerIptorrentsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerIptorrentsResourceName+": "+req.ID)
}

func (i *IndexerIptorrents) write(ctx context.Context, indexer *whisparr.IndexerResource) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerIptorrents) read(ctx context.Context) *whisparr.IndexerResource {
	return i.toIndexer().read(ctx)
}

package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerOmgwtfnzbsResourceName   = "indexer_omgwtfnzbs"
	indexerOmgwtfnzbsImplementation = "Omgwtfnzbs"
	indexerOmgwtfnzbsConfigContract = "OmgwtfnzbsSettings"
	indexerOmgwtfnzbsProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerOmgwtfnzbsResource{}
	_ resource.ResourceWithImportState = &IndexerOmgwtfnzbsResource{}
)

func NewIndexerOmgwtfnzbsResource() resource.Resource {
	return &IndexerOmgwtfnzbsResource{}
}

// IndexerOmgwtfnzbsResource defines the Omgwtfnzbs indexer implementation.
type IndexerOmgwtfnzbsResource struct {
	client *whisparr.APIClient
}

// IndexerOmgwtfnzbs describes the Omgwtfnzbs indexer data model.
type IndexerOmgwtfnzbs struct {
	Tags                    types.Set    `tfsdk:"tags"`
	MultiLanguages          types.Set    `tfsdk:"multi_languages"`
	APIKey                  types.String `tfsdk:"api_key"`
	Name                    types.String `tfsdk:"name"`
	Username                types.String `tfsdk:"username"`
	ID                      types.Int64  `tfsdk:"id"`
	DownloadClientID        types.Int64  `tfsdk:"download_client_id"`
	Priority                types.Int64  `tfsdk:"priority"`
	Delay                   types.Int64  `tfsdk:"delay"`
	EnableRss               types.Bool   `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool   `tfsdk:"enable_interactive_search"`
	EnableAutomaticSearch   types.Bool   `tfsdk:"enable_automatic_search"`
}

func (i IndexerOmgwtfnzbs) toIndexer() *Indexer {
	return &Indexer{
		Delay:                   i.Delay,
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		Priority:                i.Priority,
		DownloadClientID:        i.DownloadClientID,
		ID:                      i.ID,
		Name:                    i.Name,
		Username:                i.Username,
		APIKey:                  i.APIKey,
		Tags:                    i.Tags,
		MultiLanguages:          i.MultiLanguages,
		Implementation:          types.StringValue(indexerOmgwtfnzbsImplementation),
		ConfigContract:          types.StringValue(indexerOmgwtfnzbsConfigContract),
		Protocol:                types.StringValue(indexerOmgwtfnzbsProtocol),
	}
}

func (i *IndexerOmgwtfnzbs) fromIndexer(indexer *Indexer) {
	i.Delay = indexer.Delay
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.Priority = indexer.Priority
	i.DownloadClientID = indexer.DownloadClientID
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.Username = indexer.Username
	i.APIKey = indexer.APIKey
	i.Tags = indexer.Tags
	i.MultiLanguages = indexer.MultiLanguages
}

func (r *IndexerOmgwtfnzbsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerOmgwtfnzbsResourceName
}

func (r *IndexerOmgwtfnzbsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Omgwtfnzbs resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/whisparr/settings#indexers) and [Omgwtfnzbs](https://wiki.servarr.com/whisparr/supported#omgwtfnzbs).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
			},
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
				MarkdownDescription: "IndexerOmgwtfnzbs name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerOmgwtfnzbs ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"delay": schema.Int64Attribute{
				MarkdownDescription: "Delay.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"multi_languages": schema.SetAttribute{
				MarkdownDescription: "Languages list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerOmgwtfnzbsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerOmgwtfnzbsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerOmgwtfnzbs

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerOmgwtfnzbs
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerOmgwtfnzbsResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerOmgwtfnzbsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerOmgwtfnzbsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerOmgwtfnzbs

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerOmgwtfnzbs current value
	response, _, err := r.client.IndexerApi.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerOmgwtfnzbsResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerOmgwtfnzbsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerOmgwtfnzbsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerOmgwtfnzbs

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerOmgwtfnzbs
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerApi.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerOmgwtfnzbsResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerOmgwtfnzbsResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerOmgwtfnzbsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerOmgwtfnzbs current value
	_, err := r.client.IndexerApi.DeleteIndexer(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerOmgwtfnzbsResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerOmgwtfnzbsResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerOmgwtfnzbsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerOmgwtfnzbsResourceName+": "+req.ID)
}

func (i *IndexerOmgwtfnzbs) write(ctx context.Context, indexer *whisparr.IndexerResource, diags *diag.Diagnostics) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer, diags)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerOmgwtfnzbs) read(ctx context.Context, diags *diag.Diagnostics) *whisparr.IndexerResource {
	return i.toIndexer().read(ctx, diags)
}

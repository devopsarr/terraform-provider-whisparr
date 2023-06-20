package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	importListTMDBCollectionResourceName   = "import_list_tmdb_collection"
	importListTMDBCollectionImplementation = "TMDbCollectionImport"
	importListTMDBCollectionConfigContract = "TMDbCollectionSettings"
	importListTMDBCollectionType           = "tmdb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListTMDBCollectionResource{}
	_ resource.ResourceWithImportState = &ImportListTMDBCollectionResource{}
)

func NewImportListTMDBCollectionResource() resource.Resource {
	return &ImportListTMDBCollectionResource{}
}

// ImportListTMDBCollectionResource defines the import list implementation.
type ImportListTMDBCollectionResource struct {
	client *whisparr.APIClient
}

// ImportListTMDBCollection describes the import list data model.
type ImportListTMDBCollection struct {
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	RootFolderPath      types.String `tfsdk:"root_folder_path"`
	CollectionID        types.String `tfsdk:"collection_id"`
	ListOrder           types.Int64  `tfsdk:"list_order"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	ShouldMonitor       types.Bool   `tfsdk:"should_monitor"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	EnableAuto          types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd         types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListTMDBCollection) toImportList() *ImportList {
	return &ImportList{
		Tags:                i.Tags,
		Name:                i.Name,
		ShouldMonitor:       i.ShouldMonitor,
		MinimumAvailability: i.MinimumAvailability,
		RootFolderPath:      i.RootFolderPath,
		CollectionID:        i.CollectionID,
		ListOrder:           i.ListOrder,
		ID:                  i.ID,
		QualityProfileID:    i.QualityProfileID,
		Enabled:             i.Enabled,
		EnableAuto:          i.EnableAuto,
		SearchOnAdd:         i.SearchOnAdd,
		Implementation:      types.StringValue(importListTMDBCollectionImplementation),
		ConfigContract:      types.StringValue(importListTMDBCollectionConfigContract),
		ListType:            types.StringValue(importListTMDBCollectionType),
	}
}

func (i *ImportListTMDBCollection) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.CollectionID = importList.CollectionID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListTMDBCollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListTMDBCollectionResourceName
}

func (r *ImportListTMDBCollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List TMDB Collection resource.\nFor more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [TMDB Collection](https://wiki.servarr.com/whisparr/supported#tmdbcollectionimport).",
		Attributes: map[string]schema.Attribute{
			"enable_auto": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Optional:            true,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled flag.",
				Optional:            true,
				Computed:            true,
			},
			"search_on_add": schema.BoolAttribute{
				MarkdownDescription: "Search on add flag.",
				Optional:            true,
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"list_order": schema.Int64Attribute{
				MarkdownDescription: "List order.",
				Optional:            true,
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Required:            true,
			},
			"should_monitor": schema.BoolAttribute{
				MarkdownDescription: "Should monitor.",
				Required:            true,
			},
			"minimum_availability": schema.StringAttribute{
				MarkdownDescription: "Minimum availability.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("tba", "announced", "inCinemas", "released", "deleted"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"collection_id": schema.StringAttribute{
				MarkdownDescription: "Collection ID.",
				Required:            true,
			},
		},
	}
}

func (r *ImportListTMDBCollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListTMDBCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListTMDBCollection

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListTMDBCollection
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListTMDBCollectionResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListTMDBCollectionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListTMDBCollection

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListTMDBCollection current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTMDBCollectionResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListTMDBCollectionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListTMDBCollection

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListTMDBCollection
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListTMDBCollectionResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListTMDBCollectionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListTMDBCollection current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, importListTMDBCollectionResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListTMDBCollectionResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListTMDBCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListTMDBCollectionResourceName+": "+req.ID)
}

func (i *ImportListTMDBCollection) write(ctx context.Context, importList *whisparr.ImportListResource) {
	genericImportList := i.toImportList()
	genericImportList.write(ctx, importList)
	i.fromImportList(genericImportList)
}

func (i *ImportListTMDBCollection) read(ctx context.Context) *whisparr.ImportListResource {
	return i.toImportList().read(ctx)
}

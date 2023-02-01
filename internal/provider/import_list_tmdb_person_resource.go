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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	importListTMDBPersonResourceName   = "import_list_tmdb_person"
	importListTMDBPersonImplementation = "TMDbPersonImport"
	importListTMDBPersonConfigContract = "TMDbPersonSettings"
	importListTMDBPersonType           = "tmdb"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListTMDBPersonResource{}
	_ resource.ResourceWithImportState = &ImportListTMDBPersonResource{}
)

func NewImportListTMDBPersonResource() resource.Resource {
	return &ImportListTMDBPersonResource{}
}

// ImportListTMDBPersonResource defines the import list implementation.
type ImportListTMDBPersonResource struct {
	client *whisparr.APIClient
}

// ImportListTMDBPerson describes the import list data model.
type ImportListTMDBPerson struct {
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	RootFolderPath      types.String `tfsdk:"root_folder_path"`
	PersonID            types.String `tfsdk:"person_id"`
	ListOrder           types.Int64  `tfsdk:"list_order"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	ShouldMonitor       types.Bool   `tfsdk:"should_monitor"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	EnableAuto          types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd         types.Bool   `tfsdk:"search_on_add"`
	PersonCast          types.Bool   `tfsdk:"cast"`
	PersonCastDirector  types.Bool   `tfsdk:"cast_director"`
	PersonCastProducer  types.Bool   `tfsdk:"cast_producer"`
	PersonCastSound     types.Bool   `tfsdk:"cast_sound"`
	PersonCastWriting   types.Bool   `tfsdk:"cast_writing"`
}

func (i ImportListTMDBPerson) toImportList() *ImportList {
	return &ImportList{
		Tags:                i.Tags,
		Name:                i.Name,
		ShouldMonitor:       i.ShouldMonitor,
		MinimumAvailability: i.MinimumAvailability,
		RootFolderPath:      i.RootFolderPath,
		PersonID:            i.PersonID,
		ListOrder:           i.ListOrder,
		ID:                  i.ID,
		QualityProfileID:    i.QualityProfileID,
		Enabled:             i.Enabled,
		EnableAuto:          i.EnableAuto,
		SearchOnAdd:         i.SearchOnAdd,
		PersonCast:          i.PersonCast,
		PersonCastDirector:  i.PersonCastDirector,
		PersonCastProducer:  i.PersonCastProducer,
		PersonCastSound:     i.PersonCastSound,
		PersonCastWriting:   i.PersonCastWriting,
	}
}

func (i *ImportListTMDBPerson) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.PersonID = importList.PersonID
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
	i.PersonCast = importList.PersonCast
	i.PersonCastDirector = importList.PersonCastDirector
	i.PersonCastProducer = importList.PersonCastProducer
	i.PersonCastSound = importList.PersonCastSound
	i.PersonCastWriting = importList.PersonCastWriting
}

func (r *ImportListTMDBPersonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListTMDBPersonResourceName
}

func (r *ImportListTMDBPersonResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List TMDB Person resource.\nFor more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [TMDB Person](https://wiki.servarr.com/whisparr/supported#tmdbpersonimport).",
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
			"cast": schema.BoolAttribute{
				MarkdownDescription: "Include cast.",
				Optional:            true,
				Computed:            true,
			},
			"cast_director": schema.BoolAttribute{
				MarkdownDescription: "Include cast director.",
				Optional:            true,
				Computed:            true,
			},
			"cast_producer": schema.BoolAttribute{
				MarkdownDescription: "Include cast producer.",
				Optional:            true,
				Computed:            true,
			},
			"cast_sound": schema.BoolAttribute{
				MarkdownDescription: "Include cast sound.",
				Optional:            true,
				Computed:            true,
			},
			"cast_writing": schema.BoolAttribute{
				MarkdownDescription: "Include cast writing.",
				Optional:            true,
				Computed:            true,
			},
			"person_id": schema.StringAttribute{
				MarkdownDescription: "Person ID.",
				Required:            true,
			},
		},
	}
}

func (r *ImportListTMDBPersonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListTMDBPersonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListTMDBPerson

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListTMDBPerson
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListTMDBPersonResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListTMDBPersonResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBPersonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListTMDBPerson

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListTMDBPerson current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTMDBPersonResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListTMDBPersonResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBPersonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListTMDBPerson

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListTMDBPerson
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListTMDBPersonResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListTMDBPersonResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTMDBPersonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListTMDBPerson

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListTMDBPerson current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTMDBPersonResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListTMDBPersonResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListTMDBPersonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListTMDBPersonResourceName+": "+req.ID)
}

func (i *ImportListTMDBPerson) write(ctx context.Context, importList *whisparr.ImportListResource) {
	genericImportList := ImportList{
		Name:                types.StringValue(importList.GetName()),
		ShouldMonitor:       types.BoolValue(importList.GetShouldMonitor()),
		MinimumAvailability: types.StringValue(string(importList.GetMinimumAvailability())),
		RootFolderPath:      types.StringValue(importList.GetRootFolderPath()),
		ListOrder:           types.Int64Value(int64(importList.GetListOrder())),
		ID:                  types.Int64Value(int64(importList.GetId())),
		QualityProfileID:    types.Int64Value(int64(importList.GetQualityProfileId())),
		Enabled:             types.BoolValue(importList.GetEnabled()),
		EnableAuto:          types.BoolValue(importList.GetEnableAuto()),
		SearchOnAdd:         types.BoolValue(importList.GetSearchOnAdd()),
	}
	genericImportList.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, importList.Tags)
	genericImportList.writeFields(ctx, importList.GetFields())
	i.fromImportList(&genericImportList)
}

func (i *ImportListTMDBPerson) read(ctx context.Context) *whisparr.ImportListResource {
	tags := make([]*int32, len(i.Tags.Elements()))
	tfsdk.ValueAs(ctx, i.Tags, &tags)

	list := whisparr.NewImportListResource()
	list.SetShouldMonitor(i.ShouldMonitor.ValueBool())
	list.SetMinimumAvailability(whisparr.MovieStatusType(i.MinimumAvailability.ValueString()))
	list.SetRootFolderPath(i.RootFolderPath.ValueString())
	list.SetQualityProfileId(int32(i.QualityProfileID.ValueInt64()))
	list.SetListOrder(int32(i.ListOrder.ValueInt64()))
	list.SetEnableAuto(i.EnableAuto.ValueBool())
	list.SetEnabled(i.Enabled.ValueBool())
	list.SetSearchOnAdd(i.SearchOnAdd.ValueBool())
	list.SetListType(importListTMDBPersonType)
	list.SetConfigContract(importListTMDBPersonConfigContract)
	list.SetImplementation(importListTMDBPersonImplementation)
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetName(i.Name.ValueString())
	list.SetTags(tags)
	list.SetFields(i.toImportList().readFields(ctx))

	return list
}

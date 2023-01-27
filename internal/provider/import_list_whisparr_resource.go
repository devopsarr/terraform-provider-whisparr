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
	importListWhisparrResourceName   = "import_list_whisparr"
	importListWhisparrImplementation = "WhisparrImport"
	importListWhisparrConfigContract = "WhisparrSettings"
	importListWhisparrType           = "advanced"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListWhisparrResource{}
	_ resource.ResourceWithImportState = &ImportListWhisparrResource{}
)

func NewImportListWhisparrResource() resource.Resource {
	return &ImportListWhisparrResource{}
}

// ImportListWhisparrResource defines the import list implementation.
type ImportListWhisparrResource struct {
	client *whisparr.APIClient
}

// ImportListWhisparr describes the import list data model.
type ImportListWhisparr struct {
	ProfileIds          types.Set    `tfsdk:"profile_ids"`
	TagIds              types.Set    `tfsdk:"tag_ids"`
	Tags                types.Set    `tfsdk:"tags"`
	Name                types.String `tfsdk:"name"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	RootFolderPath      types.String `tfsdk:"root_folder_path"`
	BaseURL             types.String `tfsdk:"base_url"`
	APIKey              types.String `tfsdk:"api_key"`
	ListOrder           types.Int64  `tfsdk:"list_order"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	ShouldMonitor       types.Bool   `tfsdk:"should_monitor"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	EnableAuto          types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd         types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListWhisparr) toImportList() *ImportList {
	return &ImportList{
		ProfileIds:          i.ProfileIds,
		TagIds:              i.TagIds,
		Tags:                i.Tags,
		Name:                i.Name,
		ShouldMonitor:       i.ShouldMonitor,
		MinimumAvailability: i.MinimumAvailability,
		RootFolderPath:      i.RootFolderPath,
		BaseURL:             i.BaseURL,
		APIKey:              i.APIKey,
		ListOrder:           i.ListOrder,
		ID:                  i.ID,
		QualityProfileID:    i.QualityProfileID,
		Enabled:             i.Enabled,
		EnableAuto:          i.EnableAuto,
		SearchOnAdd:         i.SearchOnAdd,
	}
}

func (i *ImportListWhisparr) fromImportList(importList *ImportList) {
	i.ProfileIds = importList.ProfileIds
	i.TagIds = importList.TagIds
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.BaseURL = importList.BaseURL
	i.APIKey = importList.APIKey
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListWhisparrResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListWhisparrResourceName
}

func (r *ImportListWhisparrResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Whisparr resource.\nFor more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [Whisparr](https://wiki.servarr.com/whisparr/supported#whisparrimport).",
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
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
			"profile_ids": schema.SetAttribute{
				MarkdownDescription: "Profile IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *ImportListWhisparrResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListWhisparrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListWhisparr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListWhisparr
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListWhisparrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListWhisparrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListWhisparr

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListWhisparr current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListWhisparrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListWhisparrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListWhisparr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListWhisparr
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListWhisparrResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListWhisparrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListWhisparr

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListWhisparr current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListWhisparrResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListWhisparrResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListWhisparrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListWhisparrResourceName+": "+req.ID)
}

func (i *ImportListWhisparr) write(ctx context.Context, importList *whisparr.ImportListResource) {
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
	genericImportList.writeFields(ctx, importList.Fields)
	i.fromImportList(&genericImportList)
}

func (i *ImportListWhisparr) read(ctx context.Context) *whisparr.ImportListResource {
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
	list.SetListType(importListWhisparrType)
	list.SetConfigContract(importListWhisparrConfigContract)
	list.SetImplementation(importListWhisparrImplementation)
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetName(i.Name.ValueString())
	list.SetTags(tags)
	list.SetFields(i.toImportList().readFields(ctx))

	return list
}

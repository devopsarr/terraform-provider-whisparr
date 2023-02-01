package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	importListTraktUserResourceName   = "import_list_trakt_user"
	importListTraktUserImplementation = "TraktUserImport"
	importListTraktUserConfigContract = "TraktUserSettings"
	importListTraktUserType           = "trakt"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ImportListTraktUserResource{}
	_ resource.ResourceWithImportState = &ImportListTraktUserResource{}
)

func NewImportListTraktUserResource() resource.Resource {
	return &ImportListTraktUserResource{}
}

// ImportListTraktUserResource defines the import list implementation.
type ImportListTraktUserResource struct {
	client *whisparr.APIClient
}

// ImportListTraktUser describes the import list data model.
type ImportListTraktUser struct {
	Tags                      types.Set    `tfsdk:"tags"`
	Name                      types.String `tfsdk:"name"`
	MinimumAvailability       types.String `tfsdk:"minimum_availability"`
	RootFolderPath            types.String `tfsdk:"root_folder_path"`
	AuthUser                  types.String `tfsdk:"auth_user"`
	TraktAdditionalParameters types.String `tfsdk:"trakt_additional_parameters"`
	AccessToken               types.String `tfsdk:"access_token"`
	RefreshToken              types.String `tfsdk:"refresh_token"`
	Expires                   types.String `tfsdk:"expires"`
	Username                  types.String `tfsdk:"username"`
	TraktListType             types.Int64  `tfsdk:"trakt_list_type"`
	Limit                     types.Int64  `tfsdk:"limit"`
	ListOrder                 types.Int64  `tfsdk:"list_order"`
	ID                        types.Int64  `tfsdk:"id"`
	QualityProfileID          types.Int64  `tfsdk:"quality_profile_id"`
	ShouldMonitor             types.Bool   `tfsdk:"should_monitor"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	EnableAuto                types.Bool   `tfsdk:"enable_auto"`
	SearchOnAdd               types.Bool   `tfsdk:"search_on_add"`
}

func (i ImportListTraktUser) toImportList() *ImportList {
	return &ImportList{
		Tags:                      i.Tags,
		Name:                      i.Name,
		ShouldMonitor:             i.ShouldMonitor,
		MinimumAvailability:       i.MinimumAvailability,
		RootFolderPath:            i.RootFolderPath,
		ListOrder:                 i.ListOrder,
		RefreshToken:              i.RefreshToken,
		AccessToken:               i.AccessToken,
		Expires:                   i.Expires,
		AuthUser:                  i.AuthUser,
		TraktAdditionalParameters: i.TraktAdditionalParameters,
		Username:                  i.Username,
		TraktListType:             i.TraktListType,
		Limit:                     i.Limit,
		ID:                        i.ID,
		QualityProfileID:          i.QualityProfileID,
		Enabled:                   i.Enabled,
		EnableAuto:                i.EnableAuto,
		SearchOnAdd:               i.SearchOnAdd,
	}
}

func (i *ImportListTraktUser) fromImportList(importList *ImportList) {
	i.Tags = importList.Tags
	i.Name = importList.Name
	i.ShouldMonitor = importList.ShouldMonitor
	i.MinimumAvailability = importList.MinimumAvailability
	i.RootFolderPath = importList.RootFolderPath
	i.RefreshToken = importList.RefreshToken
	i.AccessToken = importList.AccessToken
	i.Expires = importList.Expires
	i.AuthUser = importList.AuthUser
	i.TraktAdditionalParameters = importList.TraktAdditionalParameters
	i.Username = importList.Username
	i.TraktListType = importList.TraktListType
	i.Limit = importList.Limit
	i.ListOrder = importList.ListOrder
	i.ID = importList.ID
	i.QualityProfileID = importList.QualityProfileID
	i.Enabled = importList.Enabled
	i.EnableAuto = importList.EnableAuto
	i.SearchOnAdd = importList.SearchOnAdd
}

func (r *ImportListTraktUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListTraktUserResourceName
}

func (r *ImportListTraktUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Import List Trakt User resource.\nFor more information refer to [Import List](https://wiki.servarr.com/whisparr/settings#import-lists) and [Trakt User](https://wiki.servarr.com/whisparr/supported#traktuserimport).",
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
			"limit": schema.Int64Attribute{
				MarkdownDescription: "limit.",
				Required:            true,
			},
			"trakt_list_type": schema.Int64Attribute{
				MarkdownDescription: "Trakt list type.`0` UserWatchList, `1` UserWatchedList, `2` UserCollectionList.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Required:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Optional:            true,
				Computed:            true,
			},
			"trakt_additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Trakt additional parameters.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *ImportListTraktUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *ImportListTraktUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ImportListTraktUser
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.CreateImportList(ctx).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+importListTraktUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get ImportListTraktUser current value
	response, _, err := r.client.ImportListApi.GetImportListById(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListTraktUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.Plan.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update ImportListTraktUser
	request := importList.read(ctx)

	response, _, err := r.client.ImportListApi.UpdateImportList(ctx, strconv.Itoa(int(request.GetId()))).ImportListResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+importListTraktUserResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	importList.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &importList)...)
}

func (r *ImportListTraktUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var importList *ImportListTraktUser

	resp.Diagnostics.Append(req.State.Get(ctx, &importList)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete ImportListTraktUser current value
	_, err := r.client.ImportListApi.DeleteImportList(ctx, int32(importList.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListTraktUserResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+importListTraktUserResourceName+": "+strconv.Itoa(int(importList.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *ImportListTraktUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+importListTraktUserResourceName+": "+req.ID)
}

func (i *ImportListTraktUser) write(ctx context.Context, importList *whisparr.ImportListResource) {
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

func (i *ImportListTraktUser) read(ctx context.Context) *whisparr.ImportListResource {
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
	list.SetListType(importListTraktUserType)
	list.SetConfigContract(importListTraktUserConfigContract)
	list.SetImplementation(importListTraktUserImplementation)
	list.SetId(int32(i.ID.ValueInt64()))
	list.SetName(i.Name.ValueString())
	list.SetTags(tags)
	list.SetFields(i.toImportList().readFields(ctx))

	return list
}

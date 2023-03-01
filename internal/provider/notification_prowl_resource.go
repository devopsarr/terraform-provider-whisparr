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
	notificationProwlResourceName   = "notification_prowl"
	notificationProwlImplementation = "Prowl"
	notificationProwlConfigContract = "ProwlSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationProwlResource{}
	_ resource.ResourceWithImportState = &NotificationProwlResource{}
)

func NewNotificationProwlResource() resource.Resource {
	return &NotificationProwlResource{}
}

// NotificationProwlResource defines the notification implementation.
type NotificationProwlResource struct {
	client *whisparr.APIClient
}

// NotificationProwl describes the notification data model.
type NotificationProwl struct {
	Tags                        types.Set    `tfsdk:"tags"`
	Name                        types.String `tfsdk:"name"`
	APIKey                      types.String `tfsdk:"api_key"`
	Priority                    types.Int64  `tfsdk:"priority"`
	ID                          types.Int64  `tfsdk:"id"`
	OnGrab                      types.Bool   `tfsdk:"on_grab"`
	OnMovieFileDeleteForUpgrade types.Bool   `tfsdk:"on_movie_file_delete_for_upgrade"`
	OnMovieFileDelete           types.Bool   `tfsdk:"on_movie_file_delete"`
	IncludeHealthWarnings       types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate         types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue               types.Bool   `tfsdk:"on_health_issue"`
	OnMovieDelete               types.Bool   `tfsdk:"on_movie_delete"`
	OnUpgrade                   types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                  types.Bool   `tfsdk:"on_download"`
}

func (n NotificationProwl) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		APIKey:                      n.APIKey,
		Priority:                    n.Priority,
		Name:                        n.Name,
		ID:                          n.ID,
		OnGrab:                      n.OnGrab,
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade,
		OnMovieFileDelete:           n.OnMovieFileDelete,
		IncludeHealthWarnings:       n.IncludeHealthWarnings,
		OnApplicationUpdate:         n.OnApplicationUpdate,
		OnHealthIssue:               n.OnHealthIssue,
		OnMovieDelete:               n.OnMovieDelete,
		OnUpgrade:                   n.OnUpgrade,
		OnDownload:                  n.OnDownload,
		ConfigContract:              types.StringValue(notificationProwlConfigContract),
		Implementation:              types.StringValue(notificationProwlImplementation),
	}
}

func (n *NotificationProwl) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.APIKey = notification.APIKey
	n.Priority = notification.Priority
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnMovieFileDeleteForUpgrade = notification.OnMovieFileDeleteForUpgrade
	n.OnMovieFileDelete = notification.OnMovieFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnMovieDelete = notification.OnMovieDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationProwlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationProwlResourceName
}

func (r *NotificationProwlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Prowl resource.\nFor more information refer to [Notification](https://wiki.servarr.com/whisparr/settings#connect) and [Prowl](https://wiki.servarr.com/whisparr/supported#prowl).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_movie_delete": schema.BoolAttribute{
				MarkdownDescription: "On movie delete flag.",
				Required:            true,
			},
			"on_movie_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_movie_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On movie file delete for upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationProwl name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.`-2` Very Low, `-1` Low, `0` Normal, `1` High, `2` Emergency.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(-2, -1, 0, 1, 2),
				},
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationProwlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationProwlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationProwl
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationProwlResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationProwlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationProwl current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationProwlResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationProwlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationProwl
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationProwlResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationProwlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationProwl

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationProwl current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationProwlResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationProwlResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationProwlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationProwlResourceName+": "+req.ID)
}

func (n *NotificationProwl) write(ctx context.Context, notification *whisparr.NotificationResource) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification)
	n.fromNotification(genericNotification)
}

func (n *NotificationProwl) read(ctx context.Context) *whisparr.NotificationResource {
	return n.toNotification().read(ctx)
}

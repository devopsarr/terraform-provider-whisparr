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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationTraktResourceName   = "notification_trakt"
	notificationTraktImplementation = "Trakt"
	notificationTraktConfigContract = "TraktSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationTraktResource{}
	_ resource.ResourceWithImportState = &NotificationTraktResource{}
)

func NewNotificationTraktResource() resource.Resource {
	return &NotificationTraktResource{}
}

// NotificationTraktResource defines the notification implementation.
type NotificationTraktResource struct {
	client *whisparr.APIClient
}

// NotificationTrakt describes the notification data model.
type NotificationTrakt struct {
	Tags                        types.Set    `tfsdk:"tags"`
	AuthUser                    types.String `tfsdk:"auth_user"`
	AccessToken                 types.String `tfsdk:"access_token"`
	RefreshToken                types.String `tfsdk:"refresh_token"`
	Expires                     types.String `tfsdk:"expires"`
	Name                        types.String `tfsdk:"name"`
	ID                          types.Int64  `tfsdk:"id"`
	OnMovieFileDeleteForUpgrade types.Bool   `tfsdk:"on_movie_file_delete_for_upgrade"`
	OnMovieFileDelete           types.Bool   `tfsdk:"on_movie_file_delete"`
	IncludeHealthWarnings       types.Bool   `tfsdk:"include_health_warnings"`
	OnMovieDelete               types.Bool   `tfsdk:"on_movie_delete"`
	OnUpgrade                   types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                  types.Bool   `tfsdk:"on_download"`
}

func (n NotificationTrakt) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		AuthUser:                    n.AuthUser,
		Name:                        n.Name,
		AccessToken:                 n.AccessToken,
		RefreshToken:                n.RefreshToken,
		ID:                          n.ID,
		Expires:                     n.Expires,
		OnMovieFileDeleteForUpgrade: n.OnMovieFileDeleteForUpgrade,
		OnMovieFileDelete:           n.OnMovieFileDelete,
		IncludeHealthWarnings:       n.IncludeHealthWarnings,
		OnMovieDelete:               n.OnMovieDelete,
		OnUpgrade:                   n.OnUpgrade,
		OnDownload:                  n.OnDownload,
		ConfigContract:              types.StringValue(notificationTraktConfigContract),
		Implementation:              types.StringValue(notificationTraktImplementation),
	}
}

func (n *NotificationTrakt) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.AuthUser = notification.AuthUser
	n.Name = notification.Name
	n.AccessToken = notification.AccessToken
	n.RefreshToken = notification.RefreshToken
	n.Expires = notification.Expires
	n.ID = notification.ID
	n.OnMovieFileDeleteForUpgrade = notification.OnMovieFileDeleteForUpgrade
	n.OnMovieFileDelete = notification.OnMovieFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnMovieDelete = notification.OnMovieDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationTraktResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationTraktResourceName
}

func (r *NotificationTraktResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Trakt resource.\nFor more information refer to [Notification](https://wiki.servarr.com/whisparr/settings#connect) and [Trakt](https://wiki.servarr.com/whisparr/supported#trakt).",
		Attributes: map[string]schema.Attribute{
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
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationTrakt name.",
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
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access Token.",
				Required:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Access Token.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Required:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "expires.",
				Computed:            true,
			},
		},
	}
}

func (r *NotificationTraktResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationTraktResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationTrakt
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationTraktResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTraktResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationTrakt current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationTraktResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTraktResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationTrakt
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationTraktResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTraktResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationTrakt

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationTrakt current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationTraktResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationTraktResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationTraktResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationTraktResourceName+": "+req.ID)
}

func (n *NotificationTrakt) write(ctx context.Context, notification *whisparr.NotificationResource) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification)
	n.fromNotification(genericNotification)
}

func (n *NotificationTrakt) read(ctx context.Context) *whisparr.NotificationResource {
	return n.toNotification().read(ctx)
}

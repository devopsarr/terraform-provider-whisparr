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
	notificationTelegramResourceName   = "notification_telegram"
	notificationTelegramImplementation = "Telegram"
	notificationTelegramConfigContract = "TelegramSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationTelegramResource{}
	_ resource.ResourceWithImportState = &NotificationTelegramResource{}
)

func NewNotificationTelegramResource() resource.Resource {
	return &NotificationTelegramResource{}
}

// NotificationTelegramResource defines the notification implementation.
type NotificationTelegramResource struct {
	client *whisparr.APIClient
}

// NotificationTelegram describes the notification data model.
type NotificationTelegram struct {
	Tags                        types.Set    `tfsdk:"tags"`
	ChatID                      types.String `tfsdk:"chat_id"`
	Name                        types.String `tfsdk:"name"`
	BotToken                    types.String `tfsdk:"bot_token"`
	ID                          types.Int64  `tfsdk:"id"`
	SendSilently                types.Bool   `tfsdk:"send_silently"`
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

func (n NotificationTelegram) toNotification() *Notification {
	return &Notification{
		Tags:                        n.Tags,
		ChatID:                      n.ChatID,
		BotToken:                    n.BotToken,
		SendSilently:                n.SendSilently,
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
		ConfigContract:              types.StringValue(notificationTelegramConfigContract),
		Implementation:              types.StringValue(notificationTelegramImplementation),
	}
}

func (n *NotificationTelegram) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.ChatID = notification.ChatID
	n.BotToken = notification.BotToken
	n.SendSilently = notification.SendSilently
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

func (r *NotificationTelegramResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationTelegramResourceName
}

func (r *NotificationTelegramResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Telegram resource.\nFor more information refer to [Notification](https://wiki.servarr.com/whisparr/settings#connect) and [Telegram](https://wiki.servarr.com/whisparr/supported#telegram).",
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
				MarkdownDescription: "NotificationTelegram name.",
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
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Send silently flag.",
				Optional:            true,
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Required:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationTelegramResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *NotificationTelegramResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationTelegram
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationTelegram current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationTelegram
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationTelegram current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationTelegramResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationTelegramResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationTelegramResourceName+": "+req.ID)
}

func (n *NotificationTelegram) write(ctx context.Context, notification *whisparr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationTelegram) read(ctx context.Context, diags *diag.Diagnostics) *whisparr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}

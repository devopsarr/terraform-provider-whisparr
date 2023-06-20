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

const remotePathMappingResourceName = "remote_path_mapping"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &RemotePathMappingResource{}
	_ resource.ResourceWithImportState = &RemotePathMappingResource{}
)

func NewRemotePathMappingResource() resource.Resource {
	return &RemotePathMappingResource{}
}

// RemotePathMappingResource defines the remote path mapping implementation.
type RemotePathMappingResource struct {
	client *whisparr.APIClient
}

// RemotePathMapping describes the remote path mapping data model.
type RemotePathMapping struct {
	Host       types.String `tfsdk:"host"`
	RemotePath types.String `tfsdk:"remote_path"`
	LocalPath  types.String `tfsdk:"local_path"`
	ID         types.Int64  `tfsdk:"id"`
}

func (r *RemotePathMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Remote Path Mapping resource.\nFor more information refer to [Remote Path Mapping](https://wiki.servarr.com/whisparr/settings#remote-path-mappings) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Remote Path Mapping ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Download Client host.",
				Required:            true,
			},
			"remote_path": schema.StringAttribute{
				MarkdownDescription: "Download Client remote path.",
				Required:            true,
			},
			"local_path": schema.StringAttribute{
				MarkdownDescription: "Local path.",
				Required:            true,
			},
		},
	}
}

func (r *RemotePathMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + remotePathMappingResourceName
}

func (r *RemotePathMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *RemotePathMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var mapping *RemotePathMapping

	resp.Diagnostics.Append(req.Plan.Get(ctx, &mapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new RemotePathMapping
	request := mapping.read()

	response, _, err := r.client.RemotePathMappingApi.CreateRemotePathMapping(ctx).RemotePathMappingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	mapping.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &mapping)...)
}

func (r *RemotePathMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var mapping *RemotePathMapping

	resp.Diagnostics.Append(req.State.Get(ctx, &mapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get remotePathMapping current value
	response, _, err := r.client.RemotePathMappingApi.GetRemotePathMappingById(ctx, int32(mapping.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	mapping.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &mapping)...)
}

func (r *RemotePathMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var mapping *RemotePathMapping

	resp.Diagnostics.Append(req.Plan.Get(ctx, &mapping)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update RemotePathMapping
	request := mapping.read()

	response, _, err := r.client.RemotePathMappingApi.UpdateRemotePathMapping(ctx, strconv.Itoa(int(request.GetId()))).RemotePathMappingResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+remotePathMappingResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	mapping.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &mapping)...)
}

func (r *RemotePathMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete remotePathMapping current value
	_, err := r.client.RemotePathMappingApi.DeleteRemotePathMapping(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, remotePathMappingResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+remotePathMappingResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *RemotePathMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+remotePathMappingResourceName+": "+req.ID)
}

func (r *RemotePathMapping) write(remotePathMapping *whisparr.RemotePathMappingResource) {
	r.ID = types.Int64Value(int64(remotePathMapping.GetId()))
	r.Host = types.StringValue(remotePathMapping.GetHost())
	r.RemotePath = types.StringValue(remotePathMapping.GetRemotePath())
	r.LocalPath = types.StringValue(remotePathMapping.GetLocalPath())
}

func (r *RemotePathMapping) read() *whisparr.RemotePathMappingResource {
	mapping := whisparr.NewRemotePathMappingResource()
	mapping.SetHost(r.Host.ValueString())
	mapping.SetLocalPath(r.LocalPath.ValueString())
	mapping.SetRemotePath(r.RemotePath.ValueString())
	mapping.SetId(int32(r.ID.ValueInt64()))

	return mapping
}

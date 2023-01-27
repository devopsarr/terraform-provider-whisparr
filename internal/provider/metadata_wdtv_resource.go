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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	metadataWdtvResourceName   = "metadata_wdtv"
	metadataWdtvImplementation = "WdtvMetadata"
	metadataWdtvConfigContract = "WdtvMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataWdtvResource{}
	_ resource.ResourceWithImportState = &MetadataWdtvResource{}
)

func NewMetadataWdtvResource() resource.Resource {
	return &MetadataWdtvResource{}
}

// MetadataWdtvResource defines the Wdtv metadata implementation.
type MetadataWdtvResource struct {
	client *whisparr.APIClient
}

// MetadataWdtv describes the Wdtv metadata data model.
type MetadataWdtv struct {
	Tags          types.Set    `tfsdk:"tags"`
	Name          types.String `tfsdk:"name"`
	ID            types.Int64  `tfsdk:"id"`
	Enable        types.Bool   `tfsdk:"enable"`
	MovieMetadata types.Bool   `tfsdk:"movie_metadata"`
	MovieImages   types.Bool   `tfsdk:"movie_images"`
}

func (m MetadataWdtv) toMetadata() *Metadata {
	return &Metadata{
		Tags:          m.Tags,
		Name:          m.Name,
		ID:            m.ID,
		Enable:        m.Enable,
		MovieMetadata: m.MovieMetadata,
		MovieImages:   m.MovieImages,
	}
}

func (m *MetadataWdtv) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.Enable = metadata.Enable
	m.MovieMetadata = metadata.MovieMetadata
	m.MovieImages = metadata.MovieImages
}

func (r *MetadataWdtvResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataWdtvResourceName
}

func (r *MetadataWdtvResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Wdtv resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/whisparr/settings#metadata) and [WDTV](https://wiki.servarr.com/whisparr/supported#wdtvmetadata).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Metadata name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Metadata ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"movie_metadata": schema.BoolAttribute{
				MarkdownDescription: "Movie metadata flag.",
				Required:            true,
			},
			"movie_images": schema.BoolAttribute{
				MarkdownDescription: "Movie images flag.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataWdtvResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataWdtvResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataWdtv
	request := metadata.read(ctx)

	response, _, err := r.client.MetadataApi.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataWdtvResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataWdtvResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataWdtvResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataWdtv current value
	response, _, err := r.client.MetadataApi.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataWdtvResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataWdtvResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataWdtvResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataWdtv
	request := metadata.read(ctx)

	response, _, err := r.client.MetadataApi.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+metadataWdtvResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataWdtvResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataWdtvResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var metadata *MetadataWdtv

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataWdtv current value
	_, err := r.client.MetadataApi.DeleteMetadata(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataWdtvResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataWdtvResourceName+": "+strconv.Itoa(int(metadata.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataWdtvResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataWdtvResourceName+": "+req.ID)
}

func (m *MetadataWdtv) write(ctx context.Context, metadata *whisparr.MetadataResource) {
	genericMetadata := Metadata{
		Name:   types.StringValue(metadata.GetName()),
		ID:     types.Int64Value(int64(metadata.GetId())),
		Enable: types.BoolValue(metadata.GetEnable()),
	}
	genericMetadata.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, metadata.Tags)
	genericMetadata.writeFields(metadata.Fields)
	m.fromMetadata(&genericMetadata)
}

func (m *MetadataWdtv) read(ctx context.Context) *whisparr.MetadataResource {
	tags := make([]*int32, len(m.Tags.Elements()))
	tfsdk.ValueAs(ctx, m.Tags, &tags)

	metadata := whisparr.NewMetadataResource()
	metadata.SetEnable(m.Enable.ValueBool())
	metadata.SetId(int32(m.ID.ValueInt64()))
	metadata.SetConfigContract(metadataWdtvConfigContract)
	metadata.SetImplementation(metadataWdtvImplementation)
	metadata.SetName(m.Name.ValueString())
	metadata.SetTags(tags)
	metadata.SetFields(m.toMetadata().readFields())

	return metadata
}

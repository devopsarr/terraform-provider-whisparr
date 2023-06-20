package provider

import (
	"context"
	"fmt"
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
	metadataKodiResourceName   = "metadata_kodi"
	metadataKodiImplementation = "XbmcMetadata"
	metadataKodiConfigContract = "XbmcMetadataSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MetadataKodiResource{}
	_ resource.ResourceWithImportState = &MetadataKodiResource{}
)

func NewMetadataKodiResource() resource.Resource {
	return &MetadataKodiResource{}
}

// MetadataKodiResource defines the Kodi metadata implementation.
type MetadataKodiResource struct {
	client *whisparr.APIClient
}

// MetadataKodi describes the Kodi metadata data model.
type MetadataKodi struct {
	Tags                  types.Set    `tfsdk:"tags"`
	Name                  types.String `tfsdk:"name"`
	ID                    types.Int64  `tfsdk:"id"`
	MovieMetadataLanguage types.Int64  `tfsdk:"movie_metadata_language"`
	Enable                types.Bool   `tfsdk:"enable"`
	MovieMetadata         types.Bool   `tfsdk:"movie_metadata"`
	MovieMetadataURL      types.Bool   `tfsdk:"movie_metadata_url"`
	MovieImages           types.Bool   `tfsdk:"movie_images"`
	UseMovieNfo           types.Bool   `tfsdk:"use_movie_nfo"`
}

func (m MetadataKodi) toMetadata() *Metadata {
	return &Metadata{
		Tags:                  m.Tags,
		Name:                  m.Name,
		ID:                    m.ID,
		MovieMetadataLanguage: m.MovieMetadataLanguage,
		Enable:                m.Enable,
		MovieMetadata:         m.MovieMetadata,
		MovieMetadataURL:      m.MovieMetadataURL,
		MovieImages:           m.MovieImages,
		UseMovieNfo:           m.UseMovieNfo,
		Implementation:        types.StringValue(metadataKodiImplementation),
		ConfigContract:        types.StringValue(metadataKodiConfigContract),
	}
}

func (m *MetadataKodi) fromMetadata(metadata *Metadata) {
	m.ID = metadata.ID
	m.Name = metadata.Name
	m.Tags = metadata.Tags
	m.MovieMetadataLanguage = metadata.MovieMetadataLanguage
	m.Enable = metadata.Enable
	m.MovieMetadata = metadata.MovieMetadata
	m.MovieImages = metadata.MovieImages
	m.MovieMetadataURL = metadata.MovieMetadataURL
	m.UseMovieNfo = metadata.UseMovieNfo
}

func (r *MetadataKodiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataKodiResourceName
}

func (r *MetadataKodiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Metadata -->Metadata Kodi resource.\nFor more information refer to [Metadata](https://wiki.servarr.com/whisparr/settings#metadata) and [KODI](https://wiki.servarr.com/whisparr/supported#xbmcmetadata).",
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
			"use_movie_nfo": schema.BoolAttribute{
				MarkdownDescription: "Use movie nfo flag.",
				Required:            true,
			},
			"movie_images": schema.BoolAttribute{
				MarkdownDescription: "Movie images flag.",
				Required:            true,
			},
			"movie_metadata": schema.BoolAttribute{
				MarkdownDescription: "Movie metadata flag.",
				Required:            true,
			},
			"movie_metadata_url": schema.BoolAttribute{
				MarkdownDescription: "Movie metadata URL flag.",
				Required:            true,
			},
			"movie_metadata_language": schema.Int64Attribute{
				MarkdownDescription: "Movie metadata language.",
				Required:            true,
			},
		},
	}
}

func (r *MetadataKodiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MetadataKodiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new MetadataKodi
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataApi.CreateMetadata(ctx).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.State.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get MetadataKodi current value
	response, _, err := r.client.MetadataApi.GetMetadataById(ctx, int32(metadata.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var metadata *MetadataKodi

	resp.Diagnostics.Append(req.Plan.Get(ctx, &metadata)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update MetadataKodi
	request := metadata.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.MetadataApi.UpdateMetadata(ctx, strconv.Itoa(int(request.GetId()))).MetadataResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, fmt.Sprintf("Unable to update "+metadataKodiResourceName+", got error: %s", err))

		return
	}

	tflog.Trace(ctx, "updated "+metadataKodiResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	metadata.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &metadata)...)
}

func (r *MetadataKodiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete MetadataKodi current value
	_, err := r.client.MetadataApi.DeleteMetadata(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, metadataKodiResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+metadataKodiResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *MetadataKodiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+metadataKodiResourceName+": "+req.ID)
}

func (m *MetadataKodi) write(ctx context.Context, metadata *whisparr.MetadataResource, diags *diag.Diagnostics) {
	genericMetadata := m.toMetadata()
	genericMetadata.write(ctx, metadata, diags)
	m.fromMetadata(genericMetadata)
}

func (m *MetadataKodi) read(ctx context.Context, diags *diag.Diagnostics) *whisparr.MetadataResource {
	return m.toMetadata().read(ctx, diags)
}

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityDefinitionResourceName = "quality_definition"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QualityDefinitionResource{}
	_ resource.ResourceWithImportState = &QualityDefinitionResource{}
)

func NewQualityDefinitionResource() resource.Resource {
	return &QualityDefinitionResource{}
}

// QualityDefinitionResource defines the quality definition implementation.
type QualityDefinitionResource struct {
	client *whisparr.APIClient
}

// QualityDefinition describes the quality definition data model.
type QualityDefinition struct {
	Title         types.String  `tfsdk:"title"`
	QualityName   types.String  `tfsdk:"quality_name"`
	Source        types.String  `tfsdk:"source"`
	MinSize       types.Float64 `tfsdk:"min_size"`
	MaxSize       types.Float64 `tfsdk:"max_size"`
	PreferredSize types.Float64 `tfsdk:"preferred_size"`
	ID            types.Int64   `tfsdk:"id"`
	QualityID     types.Int64   `tfsdk:"quality_id"`
	Resolution    types.Int64   `tfsdk:"resolution"`
}

func (r *QualityDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityDefinitionResourceName
}

func (r *QualityDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Quality Definition resource.\nFor more information refer to [Quality Definition](https://wiki.servarr.com/whisparr/settings#quality-1) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Definition ID.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Quality Definition Title.",
				Required:            true,
			},
			"min_size": schema.Float64Attribute{
				MarkdownDescription: "Minimum size MB/min.",
				Optional:            true,
				Computed:            true,
			},
			"max_size": schema.Float64Attribute{
				MarkdownDescription: "Maximum size MB/min.",
				Optional:            true,
				Computed:            true,
			},
			"preferred_size": schema.Float64Attribute{
				MarkdownDescription: "Preferred size MB/min.",
				Optional:            true,
				Computed:            true,
			},
			"quality_id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Quality Resolution.",
				Computed:            true,
			},
			"quality_name": schema.StringAttribute{
				MarkdownDescription: "Quality Name.",
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Quality source.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *QualityDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *QualityDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var definition *QualityDefinition

	resp.Diagnostics.Append(req.Plan.Get(ctx, &definition)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := definition.read()

	// Read to get the quality ID
	read, _, err := r.client.QualityDefinitionApi.GetQualityDefinitionById(ctx, request.GetId()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, qualityDefinitionResourceName, err))

		return
	}

	request.Quality.SetId(read.Quality.GetId())
	request.Quality.SetSource(read.Quality.GetSource())

	// Create new QualityDefinition
	response, _, err := r.client.QualityDefinitionApi.UpdateQualityDefinition(ctx, strconv.Itoa(int(request.GetId()))).QualityDefinitionResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, qualityDefinitionResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityDefinitionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	definition.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *QualityDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var definition *QualityDefinition

	resp.Diagnostics.Append(req.State.Get(ctx, &definition)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualitydefinition current value
	response, _, err := r.client.QualityDefinitionApi.GetQualityDefinitionById(ctx, int32(definition.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityDefinitionResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityDefinitionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	definition.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *QualityDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var definition *QualityDefinition

	resp.Diagnostics.Append(req.Plan.Get(ctx, &definition)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := definition.read()

	// Update QualityDefinition
	response, _, err := r.client.QualityDefinitionApi.UpdateQualityDefinition(ctx, strconv.Itoa(int(request.GetId()))).QualityDefinitionResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, qualityDefinitionResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityDefinitionResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	definition.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &definition)...)
}

func (r *QualityDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// QualityDefinitions cannot be really deleted just removing configuration
	tflog.Trace(ctx, "decoupled "+qualityDefinitionResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *QualityDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+qualityDefinitionResourceName+": "+req.ID)
}

func (p *QualityDefinition) write(definition *whisparr.QualityDefinitionResource) {
	p.ID = types.Int64Value(int64(definition.GetId()))
	p.MinSize = types.Float64Value(definition.GetMinSize())
	p.MaxSize = types.Float64Value(definition.GetMaxSize())
	p.PreferredSize = types.Float64Value(definition.GetPreferredSize())
	p.Title = types.StringValue(definition.GetTitle())
	p.QualityName = types.StringValue(definition.Quality.GetName())
	p.QualityID = types.Int64Value(int64(definition.Quality.GetId()))
	p.Source = types.StringValue(string(definition.Quality.GetSource()))
	p.Resolution = types.Int64Value(int64(definition.Quality.GetResolution()))
}

func (p *QualityDefinition) read() *whisparr.QualityDefinitionResource {
	quality := whisparr.NewQuality()
	quality.SetId(int32(p.QualityID.ValueInt64()))
	quality.SetName(p.QualityName.ValueString())
	quality.SetResolution(int32(p.Resolution.ValueInt64()))
	quality.SetSource(whisparr.Source(p.Source.ValueString()))

	definition := whisparr.NewQualityDefinitionResource()
	definition.SetId(int32(p.ID.ValueInt64()))
	definition.SetMaxSize(p.MaxSize.ValueFloat64())
	definition.SetMinSize(p.MinSize.ValueFloat64())
	definition.SetPreferredSize(p.PreferredSize.ValueFloat64())
	definition.SetTitle(p.Title.ValueString())
	definition.SetQuality(*quality)

	return definition
}

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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityProfileResourceName = "quality_profile"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QualityProfileResource{}
	_ resource.ResourceWithImportState = &QualityProfileResource{}
)

func NewQualityProfileResource() resource.Resource {
	return &QualityProfileResource{}
}

// QualityProfileResource defines the quality profile implementation.
type QualityProfileResource struct {
	client *whisparr.APIClient
}

// QualityProfile describes the quality profile data model.
type QualityProfile struct {
	QualityGroups     types.Set    `tfsdk:"quality_groups"`
	FormatItems       types.Set    `tfsdk:"format_items"`
	Name              types.String `tfsdk:"name"`
	Language          types.Object `tfsdk:"language"`
	ID                types.Int64  `tfsdk:"id"`
	Cutoff            types.Int64  `tfsdk:"cutoff"`
	MinFormatScore    types.Int64  `tfsdk:"min_format_score"`
	CutoffFormatScore types.Int64  `tfsdk:"cutoff_format_score"`
	UpgradeAllowed    types.Bool   `tfsdk:"upgrade_allowed"`
}

// QualityGroup is part of QualityProfile.
type QualityGroup struct {
	Qualities types.Set    `tfsdk:"qualities"`
	Name      types.String `tfsdk:"name"`
	ID        types.Int64  `tfsdk:"id"`
}

// Quality is part of QualityGroup.
type Quality struct {
	Name       types.String `tfsdk:"name"`
	Source     types.String `tfsdk:"source"`
	ID         types.Int64  `tfsdk:"id"`
	Resolution types.Int64  `tfsdk:"resolution"`
}

// FormatItem is part of QualityProfile.
type FormatItem struct {
	Name   types.String `tfsdk:"name"`
	Format types.Int64  `tfsdk:"format"`
	Score  types.Int64  `tfsdk:"score"`
}

// Language is part of QualityProfile.
type Language struct {
	Name types.String `tfsdk:"name"`
	ID   types.Int64  `tfsdk:"id"`
}

func (r *QualityProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfileResourceName
}

func (r *QualityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Profiles -->Quality Profile resource.\nFor more information refer to [Quality Profile](https://wiki.servarr.com/whisparr/settings#quality-profiles) documentation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality Profile ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality Profile Name.",
				Required:            true,
			},
			"upgrade_allowed": schema.BoolAttribute{
				MarkdownDescription: "Upgrade allowed flag.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff": schema.Int64Attribute{
				MarkdownDescription: "Quality ID to which cutoff.",
				Optional:            true,
				Computed:            true,
			},
			"cutoff_format_score": schema.Int64Attribute{
				MarkdownDescription: "Cutoff format score.",
				Optional:            true,
				Computed:            true,
			},
			"min_format_score": schema.Int64Attribute{
				MarkdownDescription: "Min format score.",
				Optional:            true,
				Computed:            true,
			},
			"language": schema.SingleNestedAttribute{
				MarkdownDescription: "Language.",
				Required:            true,
				Attributes:          r.getLanguageSchema().Attributes,
			},
			"quality_groups": schema.SetNestedAttribute{
				MarkdownDescription: "Quality groups.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualityGroupSchema().Attributes,
				},
			},
			"format_items": schema.SetNestedAttribute{
				MarkdownDescription: "Format items.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getFormatItemsSchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualityGroupSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality group ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality group name.",
				Optional:            true,
				Computed:            true,
			},
			"qualities": schema.SetNestedAttribute{
				MarkdownDescription: "Qualities in group.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: r.getQualitySchema().Attributes,
				},
			},
		},
	}
}

func (r QualityProfileResource) getQualitySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Quality ID.",
				Optional:            true,
				Computed:            true,
			},
			"resolution": schema.Int64Attribute{
				MarkdownDescription: "Resolution.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Quality name.",
				Optional:            true,
				Computed:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Source.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r QualityProfileResource) getFormatItemsSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"format": schema.Int64Attribute{
				MarkdownDescription: "Format.",
				Optional:            true,
				Computed:            true,
			},
			"score": schema.Int64Attribute{
				MarkdownDescription: "Score.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r QualityProfileResource) getLanguageSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *QualityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *QualityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	request := profile.read(ctx)

	// Create new QualityProfile
	response, _, err := r.client.QualityProfileApi.CreateQualityProfile(ctx).QualityProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get qualityprofile current value
	response, _, err := r.client.QualityProfileApi.GetQualityProfileById(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var profile *QualityProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	request := profile.read(ctx)

	// Update QualityProfile
	response, _, err := r.client.QualityProfileApi.UpdateQualityProfile(ctx, strconv.Itoa(int(request.GetId()))).QualityProfileResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+qualityProfileResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	profile.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &profile)...)
}

func (r *QualityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var profile *QualityProfile

	resp.Diagnostics.Append(req.State.Get(ctx, &profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete qualityprofile current value
	_, err := r.client.QualityProfileApi.DeleteQualityProfile(ctx, int32(profile.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, qualityProfileResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+qualityProfileResourceName+": "+strconv.Itoa(int(profile.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *QualityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+qualityProfileResourceName+": "+req.ID)
}

func (p *QualityProfile) write(ctx context.Context, profile *whisparr.QualityProfileResource) {
	p.UpgradeAllowed = types.BoolValue(profile.GetUpgradeAllowed())
	p.ID = types.Int64Value(int64(profile.GetId()))
	p.Name = types.StringValue(profile.GetName())
	p.Cutoff = types.Int64Value(int64(profile.GetCutoff()))
	p.CutoffFormatScore = types.Int64Value(int64(profile.GetCutoffFormatScore()))
	p.MinFormatScore = types.Int64Value(int64(profile.GetMinFormatScore()))
	p.QualityGroups = types.SetValueMust(QualityProfileResource{}.getQualityGroupSchema().Type(), nil)
	p.FormatItems = types.SetValueMust(QualityProfileResource{}.getFormatItemsSchema().Type(), nil)

	qualityGroups := make([]QualityGroup, len(profile.Items))
	for n, g := range profile.GetItems() {
		qualityGroups[n].write(ctx, g)
	}

	formatItems := make([]FormatItem, len(profile.FormatItems))
	for n, f := range profile.FormatItems {
		formatItems[n].write(f)
	}

	language := Language{}
	language.write(profile.Language)
	tfsdk.ValueFrom(ctx, language, QualityProfileResource{}.getLanguageSchema().Type(), &p.Language)

	tfsdk.ValueFrom(ctx, qualityGroups, p.QualityGroups.Type(ctx), &p.QualityGroups)
}

func (q *QualityGroup) write(ctx context.Context, group *whisparr.QualityProfileQualityItemResource) {
	var (
		name      string
		id        int64
		qualities []Quality
	)

	if len(group.Items) == 0 {
		name = group.Quality.GetName()
		id = int64(group.Quality.GetId())
		qualities = []Quality{{
			ID:         types.Int64Value(int64(group.Quality.GetId())),
			Name:       types.StringValue(group.Quality.GetName()),
			Source:     types.StringValue(string(group.Quality.GetSource())),
			Resolution: types.Int64Value(int64(group.Quality.GetResolution())),
		}}
	} else {
		name = group.GetName()
		id = int64(group.GetId())
		qualities = make([]Quality, len(group.GetItems()))
		for m, q := range group.GetItems() {
			qualities[m].write(q)
		}
	}

	q.Name = types.StringValue(name)
	q.ID = types.Int64Value(id)
	q.Qualities = types.SetValueMust(QualityProfileResource{}.getQualitySchema().Type(), nil)

	tfsdk.ValueFrom(ctx, qualities, q.Qualities.Type(ctx), &q.Qualities)
}

func (q *Quality) write(quality *whisparr.QualityProfileQualityItemResource) {
	q.ID = types.Int64Value(int64(quality.Quality.GetId()))
	q.Name = types.StringValue(quality.Quality.GetName())
	q.Source = types.StringValue(string(quality.Quality.GetSource()))
	q.Resolution = types.Int64Value(int64(quality.Quality.GetResolution()))
}

func (f *FormatItem) write(format *whisparr.ProfileFormatItemResource) {
	f.Name = types.StringValue(format.GetName())
	f.Format = types.Int64Value(int64(format.GetFormat()))
	f.Score = types.Int64Value(int64(format.GetScore()))
}

func (l *Language) write(language *whisparr.Language) {
	l.Name = types.StringValue(language.GetName())
	l.ID = types.Int64Value(int64(language.GetId()))
}

func (p *QualityProfile) read(ctx context.Context) *whisparr.QualityProfileResource {
	groups := make([]QualityGroup, len(p.QualityGroups.Elements()))
	tfsdk.ValueAs(ctx, p.QualityGroups, &groups)
	qualities := make([]*whisparr.QualityProfileQualityItemResource, len(groups))

	for n, g := range groups {
		q := make([]Quality, len(g.Qualities.Elements()))
		tfsdk.ValueAs(ctx, g.Qualities, &q)

		if len(q) == 0 {
			quality := whisparr.NewQuality()
			quality.SetId(int32(g.ID.ValueInt64()))
			quality.SetName(g.Name.ValueString())

			item := whisparr.NewQualityProfileQualityItemResource()
			item.SetAllowed(true)
			item.SetQuality(*quality)

			qualities[n] = item

			continue
		}

		items := make([]*whisparr.QualityProfileQualityItemResource, len(q))
		for m, q := range q {
			items[m] = q.read()
		}

		quality := whisparr.NewQualityProfileQualityItemResource()
		quality.SetId(int32(g.ID.ValueInt64()))
		quality.SetName(g.Name.ValueString())
		quality.SetAllowed(true)
		quality.SetItems(items)
		qualities[n] = quality
	}

	formats := make([]FormatItem, len(p.FormatItems.Elements()))
	tfsdk.ValueAs(ctx, p.FormatItems, &formats)
	formatItems := make([]*whisparr.ProfileFormatItemResource, len(formats))

	for n, f := range formats {
		formatItems[n] = f.read()
	}

	language := Language{}
	tfsdk.ValueAs(ctx, p.Language, &language)

	profile := whisparr.NewQualityProfileResource()
	profile.SetUpgradeAllowed(p.UpgradeAllowed.ValueBool())
	profile.SetId(int32(p.ID.ValueInt64()))
	profile.SetCutoff(int32(p.Cutoff.ValueInt64()))
	profile.SetMinFormatScore(int32(p.MinFormatScore.ValueInt64()))
	profile.SetCutoffFormatScore(int32(p.CutoffFormatScore.ValueInt64()))
	profile.SetName(p.Name.ValueString())
	profile.SetLanguage(*language.read())
	profile.SetItems(qualities)
	profile.SetFormatItems(formatItems)

	return profile
}

func (q *Quality) read() *whisparr.QualityProfileQualityItemResource {
	quality := whisparr.NewQuality()
	quality.SetName(q.Name.ValueString())
	quality.SetId(int32(q.ID.ValueInt64()))
	quality.SetSource(whisparr.Source(q.Source.ValueString()))
	quality.SetResolution(int32(q.Resolution.ValueInt64()))

	item := whisparr.NewQualityProfileQualityItemResource()
	item.SetAllowed(true)
	item.SetQuality(*quality)

	return item
}

func (f *FormatItem) read() *whisparr.ProfileFormatItemResource {
	formatItem := whisparr.NewProfileFormatItemResource()
	formatItem.SetFormat(int32(f.Format.ValueInt64()))
	formatItem.SetName(f.Name.ValueString())
	formatItem.SetScore(int32(f.Score.ValueInt64()))

	return formatItem
}

func (l *Language) read() *whisparr.Language {
	language := whisparr.NewLanguage()
	language.SetId(int32(l.ID.ValueInt64()))
	language.SetName(l.Name.ValueString())

	return language
}

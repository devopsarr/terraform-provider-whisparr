package provider

import (
	"context"
	"fmt"
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

const movieResourceName = "movie"

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &MovieResource{}
	_ resource.ResourceWithImportState = &MovieResource{}
)

func NewMovieResource() resource.Resource {
	return &MovieResource{}
}

// MovieResource defines the movie implementation.
type MovieResource struct {
	client *whisparr.APIClient
}

// Movie describes the movie data model.
type Movie struct {
	Genres              types.Set    `tfsdk:"genres"`
	Tags                types.Set    `tfsdk:"tags"`
	OriginalLanguage    types.Object `tfsdk:"original_language"`
	Title               types.String `tfsdk:"title"`
	Path                types.String `tfsdk:"path"`
	MinimumAvailability types.String `tfsdk:"minimum_availability"`
	OriginalTitle       types.String `tfsdk:"original_title"`
	Status              types.String `tfsdk:"status"`
	IMDBID              types.String `tfsdk:"imdb_id"`
	YouTubeTrailerID    types.String `tfsdk:"youtube_trailer_id"`
	Overview            types.String `tfsdk:"overview"`
	Website             types.String `tfsdk:"website"`
	ID                  types.Int64  `tfsdk:"id"`
	QualityProfileID    types.Int64  `tfsdk:"quality_profile_id"`
	TMDBID              types.Int64  `tfsdk:"tmdb_id"`
	Year                types.Int64  `tfsdk:"year"`
	IsAvailable         types.Bool   `tfsdk:"is_available"`
	Monitored           types.Bool   `tfsdk:"monitored"`

	// TODO: future Implementation
	// SortTitle      types.String  `tfsdk:"sortTitle"`
	// SizeOnDisk     types.Int64   `tfsdk:"sizeOnDisk"`
	// RemotePoster   types.String  `tfsdk:"remotePoster"`
	// HasFile        types.Bool    `tfsdk:"hasFile"`
	// Studio         types.String  `tfsdk:"studio"`
	// RootFolderPath types.String  `tfsdk:"root_folder_path"`
	// FolderName     types.String  `tfsdk:"folderName"`
	// Runtime        types.Int64   `tfsdk:"runtime"`
	// CleanTitle     types.String  `tfsdk:"cleanTitle"`
	// TitleSlug      types.String  `tfsdk:"titleSlug"`
	// Folder         types.String  `tfsdk:"folder"`
	// Certification  types.String  `tfsdk:"certification"`
	// Added          types.String  `tfsdk:"added"`
	// Popularity     types.Float64 `tfsdk:"popularity"`
	// Images         types.Set     `tfsdk:"images"`
	// Ratings        types.Object  `tfsdk:"ratings"`
	// MovieFile      types.Object  `tfsdk:"movieFile"`
	// Collection     types.Object  `tfsdk:"collection"`
}

func (r *MovieResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + movieResourceName
}

func (r *MovieResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Movies -->Movie resource.\nFor more information refer to [Movies](https://wiki.servarr.com/whisparr/library#movies) documentation.",
		Attributes: map[string]schema.Attribute{
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Required:            true,
			},
			"is_available": schema.BoolAttribute{
				MarkdownDescription: "Availability flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Required:            true,
			},
			"tmdb_id": schema.Int64Attribute{
				MarkdownDescription: "TMDB ID.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Movie ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"year": schema.Int64Attribute{
				MarkdownDescription: "Year.",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Movie title.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Full movie path.",
				Required:            true,
			},
			"minimum_availability": schema.StringAttribute{
				MarkdownDescription: "Minimum availability.\nAllowed values: 'tba', 'announced', 'inCinemas', 'released', 'deleted'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("tba", "announced", "inCinemas", "released", "deleted"),
				},
			},
			"original_title": schema.StringAttribute{
				MarkdownDescription: "Movie original title.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Movie status.",
				Computed:            true,
			},
			"overview": schema.StringAttribute{
				MarkdownDescription: "Overview.",
				Computed:            true,
			},
			"website": schema.StringAttribute{
				MarkdownDescription: "Website.",
				Computed:            true,
			},
			"imdb_id": schema.StringAttribute{
				MarkdownDescription: "IMDB ID.",
				Computed:            true,
			},
			"youtube_trailer_id": schema.StringAttribute{
				MarkdownDescription: "Youtube trailer ID.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"genres": schema.SetAttribute{
				MarkdownDescription: "List genres.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"original_language": schema.SingleNestedAttribute{
				MarkdownDescription: "Original language.",
				Computed:            true,
				Attributes:          QualityProfileResource{}.getQualityLanguageSchema().Attributes,
			},
		},
	}
}

func (r *MovieResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *MovieResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var movie *Movie

	resp.Diagnostics.Append(req.Plan.Get(ctx, &movie)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new Movie
	request := movie.read(ctx)
	// TODO: can parametrize AddMovieOptions
	options := whisparr.NewAddMovieOptions()
	options.SetSearchForMovie(true)

	response, _, err := r.client.MovieApi.CreateMovie(ctx).MovieResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, movieResourceName, err))

		return
	}

	tflog.Trace(ctx, "created movie: "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	movie.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &movie)...)
}

func (r *MovieResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var movie *Movie

	resp.Diagnostics.Append(req.State.Get(ctx, &movie)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get movie current value
	response, _, err := r.client.MovieApi.GetMovieById(ctx, int32(movie.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, movieResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+movieResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	movie.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &movie)...)
}

func (r *MovieResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var movie *Movie

	resp.Diagnostics.Append(req.Plan.Get(ctx, &movie)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update Movie
	request := movie.read(ctx)

	response, _, err := r.client.MovieApi.UpdateMovie(ctx, fmt.Sprint(request.GetId())).MovieResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, movieResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+movieResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	movie.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &movie)...)
}

func (r *MovieResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var movie *Movie

	resp.Diagnostics.Append(req.State.Get(ctx, &movie)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete movie current value
	_, err := r.client.MovieApi.DeleteMovie(ctx, int32(movie.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, movieResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+movieResourceName+": "+strconv.Itoa(int(movie.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *MovieResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+movieResourceName+": "+req.ID)
}

func (m *Movie) write(ctx context.Context, movie *whisparr.MovieResource) {
	m.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, movie.GetTags())
	m.Monitored = types.BoolValue(movie.GetMonitored())
	m.ID = types.Int64Value(int64(movie.GetId()))
	m.Title = types.StringValue(movie.GetTitle())
	m.Path = types.StringValue(movie.GetPath())
	m.QualityProfileID = types.Int64Value(int64(movie.GetQualityProfileId()))
	m.TMDBID = types.Int64Value(int64(movie.GetTmdbId()))
	m.MinimumAvailability = types.StringValue(string(movie.GetMinimumAvailability()))
	// Read only values
	m.IsAvailable = types.BoolValue(movie.GetIsAvailable())
	m.OriginalTitle = types.StringValue(movie.GetOriginalTitle())
	m.Status = types.StringValue(string(movie.GetStatus()))
	m.Year = types.Int64Value(int64(movie.GetYear()))
	m.IMDBID = types.StringValue(movie.GetImdbId())
	m.YouTubeTrailerID = types.StringValue(movie.GetYouTubeTrailerId())
	m.Overview = types.StringValue(movie.GetOverview())
	m.Website = types.StringValue(movie.GetWebsite())
	language := QualityLanguage{}
	language.write(movie.OriginalLanguage)
	tfsdk.ValueFrom(ctx, language, QualityProfileResource{}.getQualityLanguageSchema().Type(), &m.OriginalLanguage)
	m.Genres = types.SetValueMust(types.StringType, nil)
	tfsdk.ValueFrom(ctx, movie.Genres, m.Genres.Type(ctx), &m.Genres)
}

func (m *Movie) read(ctx context.Context) *whisparr.MovieResource {
	tags := make([]*int32, len(m.Tags.Elements()))
	tfsdk.ValueAs(ctx, m.Tags, &tags)

	movie := whisparr.NewMovieResource()
	movie.SetMonitored(m.Monitored.ValueBool())
	movie.SetTitle(m.Title.ValueString())
	movie.SetPath(m.Path.ValueString())
	movie.SetQualityProfileId(int32(m.QualityProfileID.ValueInt64()))
	movie.SetTmdbId(int32(m.TMDBID.ValueInt64()))
	movie.SetId(int32(m.ID.ValueInt64()))
	movie.SetTags(tags)

	if !m.MinimumAvailability.IsNull() && !m.MinimumAvailability.IsUnknown() {
		movie.SetMinimumAvailability(whisparr.MovieStatusType(m.MinimumAvailability.ValueString()))
	}

	return movie
}

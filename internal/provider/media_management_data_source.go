package provider

import (
	"context"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const mediaManagementDataSourceName = "media_management"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MediaManagementDataSource{}

func NewMediaManagementDataSource() datasource.DataSource {
	return &MediaManagementDataSource{}
}

// MediaManagementDataSource defines the media management implementation.
type MediaManagementDataSource struct {
	client *whisparr.APIClient
}

func (d *MediaManagementDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + mediaManagementDataSourceName
}

func (d *MediaManagementDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->[Media Management](../resources/media_management).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"auto_rename_folders": schema.BoolAttribute{
				MarkdownDescription: "Auto rename folders.",
				Computed:            true,
			},
			"auto_unmonitor_previously_downloaded_movies": schema.BoolAttribute{
				MarkdownDescription: "Auto unmonitor previously downloaded movies.",
				Computed:            true,
			},
			"copy_using_hardlinks": schema.BoolAttribute{
				MarkdownDescription: "Use hardlinks instead of copy.",
				Computed:            true,
			},
			"create_empty_movie_folders": schema.BoolAttribute{
				MarkdownDescription: "Create empty movies directories.",
				Computed:            true,
			},
			"delete_empty_folders": schema.BoolAttribute{
				MarkdownDescription: "Delete empty movies directories.",
				Computed:            true,
			},
			"enable_media_info": schema.BoolAttribute{
				MarkdownDescription: "Scan files details.",
				Computed:            true,
			},
			"import_extra_files": schema.BoolAttribute{
				MarkdownDescription: "Import extra files. If enabled it will leverage 'extra_file_extensions'.",
				Computed:            true,
			},
			"paths_default_static": schema.BoolAttribute{
				MarkdownDescription: "Path default static.",
				Computed:            true,
			},
			"set_permissions_linux": schema.BoolAttribute{
				MarkdownDescription: "Set permission for imported files.",
				Computed:            true,
			},
			"skip_free_space_check_when_importing": schema.BoolAttribute{
				MarkdownDescription: "Skip free space check before importing.",
				Computed:            true,
			},
			"minimum_free_space_when_importing": schema.Int64Attribute{
				MarkdownDescription: "Minimum free space in MB to allow import.",
				Computed:            true,
			},
			"recycle_bin_cleanup_days": schema.Int64Attribute{
				MarkdownDescription: "Recyle bin days of retention.",
				Computed:            true,
			},
			"chmod_folder": schema.StringAttribute{
				MarkdownDescription: "Permission in linux format.",
				Computed:            true,
			},
			"chown_group": schema.StringAttribute{
				MarkdownDescription: "Group used for permission.",
				Computed:            true,
			},
			"download_propers_and_repacks": schema.StringAttribute{
				MarkdownDescription: "Download proper and repack policy. valid inputs are: 'preferAndUpgrade', 'doNotUpgrade', and 'doNotPrefer'.",
				Computed:            true,
			},
			"extra_file_extensions": schema.StringAttribute{
				MarkdownDescription: "Comma separated list of extra files to import (.nfo will be imported as .nfo-orig).",
				Computed:            true,
			},
			"file_date": schema.StringAttribute{
				MarkdownDescription: "Define the file date modification. valid inputs are: 'none', 'localAirDate, and 'utcAirDate'.",
				Computed:            true,
			},
			"recycle_bin": schema.StringAttribute{
				MarkdownDescription: "Recycle bin absolute path.",
				Computed:            true,
			},
			"rescan_after_refresh": schema.StringAttribute{
				MarkdownDescription: "Rescan after refresh policy. valid inputs are: 'always', 'afterManual' and 'never'.",
				Computed:            true,
			},
		},
	}
}

func (d *MediaManagementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MediaManagementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.MediaManagementConfigApi.GetMediaManagementConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, mediaManagementDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+mediaManagementDataSourceName)

	state := MediaManagement{}
	state.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

package provider

import (
	"context"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const downloadClientConfigDataSourceName = "download_client_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DownloadClientConfigDataSource{}

func NewDownloadClientConfigDataSource() datasource.DataSource {
	return &DownloadClientConfigDataSource{}
}

// DownloadClientConfigDataSource defines the download client config implementation.
type DownloadClientConfigDataSource struct {
	client *whisparr.APIClient
}

func (d *DownloadClientConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientConfigDataSourceName
}

func (d *DownloadClientConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Download Clients -->[Download Client Config](../resources/download_client_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client Config ID.",
				Computed:            true,
			},
			"enable_completed_download_handling": schema.BoolAttribute{
				MarkdownDescription: "Enable Completed Download Handling flag.",
				Computed:            true,
			},
			"auto_redownload_failed": schema.BoolAttribute{
				MarkdownDescription: "Auto Redownload Failed flag.",
				Computed:            true,
			},
			"check_for_finished_download_interval": schema.Int64Attribute{
				MarkdownDescription: "Check for finished download interval.",
				Computed:            true,
			},
			"download_client_working_folders": schema.StringAttribute{
				MarkdownDescription: "Download Client Working Folders.",
				Computed:            true,
			},
		},
	}
}

func (d *DownloadClientConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *DownloadClientConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.DownloadClientConfigApi.GetDownloadClientConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientConfigDataSourceName)

	config := DownloadClientConfig{}
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-whisparr/tools"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerConfigDataSourceName = "indexer_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerConfigDataSource{}

func NewIndexerConfigDataSource() datasource.DataSource {
	return &IndexerConfigDataSource{}
}

// IndexerConfigDataSource defines the indexer config implementation.
type IndexerConfigDataSource struct {
	client *whisparr.APIClient
}

func (d *IndexerConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerConfigDataSourceName
}

func (d *IndexerConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->[Indexer Config](../resources/indexer_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"maximum_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum size.",
				Computed:            true,
			},
			"minimum_age": schema.Int64Attribute{
				MarkdownDescription: "Minimum age.",
				Computed:            true,
			},
			"retention": schema.Int64Attribute{
				MarkdownDescription: "Retention.",
				Computed:            true,
			},
			"rss_sync_interval": schema.Int64Attribute{
				MarkdownDescription: "RSS sync interval.",
				Computed:            true,
			},
			"availability_delay": schema.Int64Attribute{
				MarkdownDescription: "Availability delay.",
				Computed:            true,
			},
			"whitelisted_hardcoded_subs": schema.StringAttribute{
				MarkdownDescription: "Whitelisted hardconded subs.",
				Computed:            true,
			},
			"prefer_indexer_flags": schema.BoolAttribute{
				MarkdownDescription: "Prefer indexer flags.",
				Computed:            true,
			},
			"allow_hardcoded_subs": schema.BoolAttribute{
				MarkdownDescription: "Allow hardcoded subs.",
				Computed:            true,
			},
		},
	}
}

func (d *IndexerConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*whisparr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *whisparr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *IndexerConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.IndexerConfigApi.GetIndexerConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerConfigDataSourceName)

	status := IndexerConfig{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}

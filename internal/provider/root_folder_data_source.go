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

const rootFolderDataSourceName = "root_folder"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootFolderDataSource{}

func NewRootFolderDataSource() datasource.DataSource {
	return &RootFolderDataSource{}
}

// RootFolderDataSource defines the root folders implementation.
type RootFolderDataSource struct {
	client *whisparr.APIClient
}

func (d *RootFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + rootFolderDataSourceName
}

func (d *RootFolderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->Single [Root Folder](../resources/root_folder).",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Root Folder absolute path.",
				Required:            true,
			},
			"accessible": schema.BoolAttribute{
				MarkdownDescription: "Access flag.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Root Folder ID.",
				Computed:            true,
			},
			"unmapped_folders": schema.SetNestedAttribute{
				MarkdownDescription: "List of folders with no associated series.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							MarkdownDescription: "Path of unmapped folder.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of unmapped folder.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RootFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RootFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var folder *RootFolder

	resp.Diagnostics.Append(req.Config.Get(ctx, &folder)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get rootfolders current value
	response, _, err := d.client.RootFolderApi.ListRootFolder(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", rootFolderDataSourceName, err))

		return
	}

	// Map response body to resource schema attribute
	rootFolder, err := findRootFolder(folder.Path.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(tools.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", rootFolderDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+rootFolderDataSourceName)
	folder.write(ctx, rootFolder)
	resp.Diagnostics.Append(resp.State.Set(ctx, &folder)...)
}

func findRootFolder(path string, folders []*whisparr.RootFolderResource) (*whisparr.RootFolderResource, error) {
	for _, f := range folders {
		if f.GetPath() == path {
			return f, nil
		}
	}

	return nil, tools.ErrDataNotFoundError(rootFolderDataSourceName, "path", path)
}

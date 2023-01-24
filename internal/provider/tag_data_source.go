package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tagDataSourceName = "tag"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TagDataSource{}

func NewTagDataSource() datasource.DataSource {
	return &TagDataSource{}
}

// TagDataSource defines the tag implementation.
type TagDataSource struct {
	client *whisparr.APIClient
}

func (d *TagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagDataSourceName
}

func (d *TagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Tags -->Single [Tag](../resources/tag).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Tag ID.",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Tag label.",
				Required:            true,
			},
		},
	}
}

func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tag *Tag

	resp.Diagnostics.Append(req.Config.Get(ctx, &tag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, _, err := d.client.TagApi.ListTag(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagDataSourceName, err))

		return
	}

	value, err := findTag(tag.Label.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", tagDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+tagDataSourceName)
	tag.write(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func findTag(label string, tags []*whisparr.TagResource) (*whisparr.TagResource, error) {
	for _, t := range tags {
		if t.GetLabel() == label {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(tagDataSourceName, "label", label)
}

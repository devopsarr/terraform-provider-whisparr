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

const customFormatDataSourceName = "custom_format"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatDataSource{}

func NewCustomFormatDataSource() datasource.DataSource {
	return &CustomFormatDataSource{}
}

// CustomFormatDataSource defines the custom_format implementation.
type CustomFormatDataSource struct {
	client *whisparr.APIClient
}

func (d *CustomFormatDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatDataSourceName
}

func (d *CustomFormatDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->Single [Download Client](../resources/custom_format).",
		Attributes: map[string]schema.Attribute{
			"include_custom_format_when_renaming": schema.BoolAttribute{
				MarkdownDescription: "Include custom format when renaming flag.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Custom Format name.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom Format ID.",
				Computed:            true,
			},
			"specifications": schema.SetNestedAttribute{
				MarkdownDescription: "Specifications.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"negate": schema.BoolAttribute{
							MarkdownDescription: "Negate flag.",
							Computed:            true,
						},
						"required": schema.BoolAttribute{
							MarkdownDescription: "Computed flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Specification name.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Implementation.",
							Computed:            true,
						},
						// Field values
						"value": schema.StringAttribute{
							MarkdownDescription: "Value.",
							Computed:            true,
						},
						"min": schema.Int64Attribute{
							MarkdownDescription: "Min.",
							Computed:            true,
						},
						"max": schema.Int64Attribute{
							MarkdownDescription: "Max.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *CustomFormatDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *CustomFormatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormat

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get customFormat current value
	response, _, err := d.client.CustomFormatApi.ListCustomFormat(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, customFormatDataSourceName, err))

		return
	}

	customFormat, err := findCustomFormat(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", customFormatDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatDataSourceName)
	data.write(ctx, customFormat)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findCustomFormat(name string, customFormats []*whisparr.CustomFormatResource) (*whisparr.CustomFormatResource, error) {
	for _, i := range customFormats {
		if i.GetName() == name {
			return i, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(customFormatDataSourceName, "name", name)
}

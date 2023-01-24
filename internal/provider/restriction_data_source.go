package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const restrictionDataSourceName = "restriction"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RestrictionDataSource{}

func NewRestrictionDataSource() datasource.DataSource {
	return &RestrictionDataSource{}
}

// RestrictionDataSource defines the remote path restriction implementation.
type RestrictionDataSource struct {
	client *whisparr.APIClient
}

func (d *RestrictionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + restrictionDataSourceName
}

func (d *RestrictionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->Single [Restriction](../resources/restriction).",
		Attributes: map[string]schema.Attribute{
			"required": schema.StringAttribute{
				MarkdownDescription: "Required.",
				Computed:            true,
			},
			"ignored": schema.StringAttribute{
				MarkdownDescription: "Ignored.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Restriction ID.",
				Required:            true,
			},
		},
	}
}

func (d *RestrictionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *RestrictionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var restriction *Restriction

	resp.Diagnostics.Append(req.Config.Get(ctx, &restriction)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get remote path restriction current value
	response, _, err := d.client.RestrictionApi.ListRestriction(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, restrictionDataSourceName, err))

		return
	}

	// Map response body to resource schema attribute
	value, err := findRestriction(restriction.ID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", restrictionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+restrictionDataSourceName)

	restriction.write(ctx, value)
	resp.Diagnostics.Append(resp.State.Set(ctx, &restriction)...)
}

func findRestriction(id int64, restrictions []*whisparr.RestrictionResource) (*whisparr.RestrictionResource, error) {
	for _, m := range restrictions {
		if int64(m.GetId()) == id {
			return m, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(restrictionDataSourceName, "id", strconv.Itoa(int(id)))
}
